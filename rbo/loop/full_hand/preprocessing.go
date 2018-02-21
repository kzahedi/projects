package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gonum/stat"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func printFrames(t Trajectory) {
	for _, f := range t.Frame {
		fmt.Println(
			f.Position.X, " ", f.Position.Y, " ", f.Position.Z,
			f.Quaternion.X, " ", f.Quaternion.Y, " ", f.Quaternion.Z, " ", f.Quaternion.W)
	}
}

// transformIntoWristFrame transforms all coordinate frames
// (position and orientation) into the coordinate frame
// located in the wrist
func transformIntoWristFrame(data Data) Data {
	// wrist frame is the 'first' trajectory in the data set
	r := Data{Trajectories: make([]Trajectory, data.NrOfTrajectories-1, data.NrOfTrajectories-1),
		NrOfDataPoints: data.NrOfDataPoints, NrOfTrajectories: data.NrOfTrajectories - 1}

	wrist := data.Trajectories[0]

	for trajectoryIndex := 1; trajectoryIndex < data.NrOfTrajectories; trajectoryIndex++ {
		for frameIndex := 0; frameIndex < data.NrOfDataPoints; frameIndex++ {
			r.Trajectories[trajectoryIndex-1].Frame =
				append(r.Trajectories[trajectoryIndex-1].Frame,
					PoseCopy(data.Trajectories[trajectoryIndex].Frame[frameIndex]))
		}
	}

	// translate all frames with respect to wrist frame
	for trajectoryIndex := 0; trajectoryIndex < r.NrOfTrajectories; trajectoryIndex++ {
		for frameIndex := 0; frameIndex < r.NrOfDataPoints; frameIndex++ {
			origPose := r.Trajectories[trajectoryIndex].Frame[frameIndex]
			wristPose := wrist.Frame[frameIndex]
			newPose := PoseSub(origPose, wristPose)
			r.Trajectories[trajectoryIndex].Frame[frameIndex] = newPose
		}
	}

	return r
}

func ConvertSofaStates(filename string, hands, ctrls []*regexp.Regexp, directory *string, convertToWristFrame bool) {
	fmt.Println("Converting sofa state files:", filename)
	files := ListAllFilesRecursivelyByFilename(*directory, filename)
	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(files, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}

	bar := pb.StartNew(iterations)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(files, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			for _, s := range rbohand2Files {
				data := ReadSofaSates(s) // returns 2d-array of pose
				if convertToWristFrame {
					data = transformIntoWristFrame(data)
					// e.g. finger tip to finger root, palm tip to palm root, thumb tip to thumb root
					data = extractTipToRootData(data)
				}
				outfile := strings.Replace(s, "raw", "analysis", 1)
				outfile = strings.Replace(outfile, "txt", "csv", 1)
				CreateDir(outfile)
				WritePositions(outfile, data)
				bar.Increment()
			}
		}
	}
	bar.Finish()
}

func extractTipToRootData(data Data) Data {
	// 6 = index finder tip to root
	// 5 = middle finger tip to root
	// 4 = ring finger tip to root
	// 3 = pinky finger tip to root
	// 2 = palm finger tip to root
	// 1 = thumb finger tip to root
	r := Data{Trajectories: make([]Trajectory, 6, 6), NrOfDataPoints: data.NrOfDataPoints, NrOfTrajectories: 6}

	indices = [][]int{{24, 29}, {20, 24}, {15, 19}, {10, 14}, {5, 9}, {0, 4}}

	for i, v := range indices {
		root = data.Trajectories[v[0]]
		tip = data.Trajectories[v[1]]
		r.Trajectories[i].Frame = make([]Data, data.NrOfDataPoints, data.NrOfDataPoints)
		for j := 0; j < data.NrOfDataPoints; j++ {
			r.Trajectories[i].Frame[j] = PoseSub(tip, root)
		}
	}

}

// tested
func calculateDifferencePositionOnly(grasp, prescriptive Data) Data {
	r := Data{Trajectories: make([]Trajectory, grasp.NrOfTrajectories, grasp.NrOfTrajectories), NrOfTrajectories: grasp.NrOfTrajectories, NrOfDataPoints: grasp.NrOfDataPoints}

	for trajectoryIndex := 0; trajectoryIndex < grasp.NrOfTrajectories; trajectoryIndex++ {
		r.Trajectories[trajectoryIndex].Frame = make([]Pose, r.NrOfDataPoints, r.NrOfDataPoints)
		for frameIndex := 0; frameIndex < grasp.NrOfDataPoints; frameIndex++ {
			g := grasp.Trajectories[trajectoryIndex].Frame[frameIndex]
			p := prescriptive.Trajectories[trajectoryIndex].Frame[frameIndex]
			diff := PoseSub(g, p)
			r.Trajectories[trajectoryIndex].Frame[frameIndex] = diff
		}
	}
	return r
}

func CalculateDifferenceBehaviour(hands, ctrls []*regexp.Regexp, prescriptive *regexp.Regexp, directory *string) {
	fmt.Println("Calculating difference behaviour")
	filename := "hand.sofastates.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(files, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}

	bar := pb.StartNew(iterations)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Grasps := Select(files, *hand)
			rbohand2Grasps = Select(rbohand2Grasps, *ctrl)

			rbohand2Prescriptives := Select(files, *prescriptive)
			rbohand2Prescriptives = Select(rbohand2Prescriptives, *ctrl)
			prescritiveData := ReadCSVToData(rbohand2Prescriptives[0])

			for _, s := range rbohand2Grasps {
				data := ReadCSVToData(s) // returns 2d-array of pose
				diff := calculateDifferencePositionOnly(data, prescritiveData)
				output := strings.Replace(s, filename, fmt.Sprintf("difference.%s", filename), 1)
				WritePositions(output, diff)
				bar.Increment()
			}
		}
	}
	bar.Finish()
}

func CalculateCovarianceMatrices(hands, ctrls []*regexp.Regexp, directory *string, max int) {
	fmt.Println("Calculating covariance matrices")
	filename := "difference.hand.sofastates.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(files, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}
	bar := pb.StartNew(iterations)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			differenceBehaviours := Select(files, *hand)
			differenceBehaviours = Select(differenceBehaviours, *ctrl)

			for _, s := range differenceBehaviours {

				data := ReadCSVToFloat(s)
				cols := len(data[0])

				start := 10
				stop := 10 + max
				r := make([][]string, cols, cols)
				for i := 0; i < cols; i++ {
					r[i] = make([]string, cols, cols)
					for j := 0; j < cols; j++ {
						di := getColumn(data, i)
						dj := getColumn(data, j)
						di = di[start:stop]
						dj = dj[start:stop]
						r[i][j] = fmt.Sprintf("%f", stat.Covariance(di, dj, nil))
						// r[i][j] = fmt.Sprintf("%f", stat.Correlation(di, dj, nil))
					}
				}

				output := strings.Replace(s, filename, "covariance.csv", 1)
				WriteCSV(output, r)
				bar.Increment()
			}
		}
	}
	bar.Finish()
}

func CreateResultsContainer(hands, ctrls []*regexp.Regexp, directory *string, results *Results) {
	filename := "hand.sofastates.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			hfiles := Select(files, *hand)
			hfiles = Select(hfiles, *ctrl)
			for _, s := range hfiles {
				key := GetKey(s)
				r := Result{MC_W: 0.0, GraspDistance: 0.0, Point: []float64{0.0, 0.0}, ObjectType: -1, ObjectPosition: -1, ClusteredByTSE: false}
				(*results)[key] = r
			}
		}
	}

}

func GetKey(s string) string {
	re := regexp.MustCompile("rbo[a-zA-Z0-9-]+/[a-zA-Z0-9_.-]+")
	return re.FindAllString(s, -1)[0]
}

func GetObjectName(s string) string {
	re := regexp.MustCompile("object[a-zA-Z0-9-]+")
	return re.FindAllString(s, -1)[0]
}

func getColumn(data [][]float64, col int) []float64 {
	r := make([]float64, len(data), len(data))
	for row := 0; row < len(data); row++ {
		r[row] = data[row][col]
	}
	return r
}
