package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/gonum/matrix/mat64"
	"github.com/gonum/stat"
	pb "gopkg.in/cheggaaa/pb.v1"
)

// CalculateGlobalVelocities This function calculates the velocities and angular
// velocities in the world coordinate system, which means that
// velocity(t)         = position(t) - position(t-1) and
// angular velocity(t) = orientation(t) - orientation(t-1)
// velocity(0)         = 0
// angular velocity(0) = 0
func CalculateGlobalVelocities(data Data) Data {
	for trajectoryIndex := 0; trajectoryIndex < data.NrOfTrajectories; trajectoryIndex++ {
		data.Trajectories[trajectoryIndex].GlobalVelocity = make([]Pose, data.NrOfDataPoints, data.NrOfDataPoints)
		for frameIndex := 1; frameIndex < data.NrOfDataPoints; frameIndex++ {
			diff :=
				PoseSub(
					data.Trajectories[trajectoryIndex].Frame[frameIndex],
					data.Trajectories[trajectoryIndex].Frame[frameIndex-1])
			data.Trajectories[trajectoryIndex].GlobalVelocity[frameIndex] = diff
		}
	}
	return data
}

func rollPitchYawRotationMatrixInverse(z, y, x float64) *mat64.Dense {
	r := mat64.NewDense(3, 3, nil)
	cx := math.Cos(x)
	cy := math.Cos(y)
	cz := math.Cos(z)

	sx := math.Sin(x)
	sy := math.Sin(y)
	sz := math.Sin(z)

	r.Set(0, 0, cy*cz)
	r.Set(0, 1, -cy*sz)
	r.Set(0, 2, sy)

	r.Set(1, 0, cz*sx*sy+cx*sz)
	r.Set(1, 1, cx*cz-sx*sy*sz)
	r.Set(1, 2, -cy*sx)

	r.Set(2, 0, -cx*cz*sy+sx*sz)
	r.Set(2, 1, cz*sx+cx*sy*sz)
	r.Set(2, 2, cx*cy)

	return r
}

func transformRotationToLocalCoordinateFrame(pos P3D, orientation P3D) P3D {
	rotationMatrix := rollPitchYawRotationMatrixInverse(orientation.Z, orientation.Y, orientation.X)
	rotationMatrix.Inverse(rotationMatrix)

	return P3D{X: 1.0, Y: 2.0, Z: 3.0}
}

func printFrames(t Trajectory) {
	for _, f := range t.Frame {
		fmt.Println(
			f.Position.X, " ", f.Position.Y, " ", f.Position.Z,
			f.Orientation.X, " ", f.Orientation.Y, " ", f.Orientation.Z)
	}
}

// TransformIntoWristFrame transforms all coordinate frames
// (position and orientation) into the coordinate frame
// located in the wrist
// TESTED
func TransformIntoWristFrame(data Data) Data {
	// wrist frame is the 'first' trajectory in the data set
	r := Data{Trajectories: make([]Trajectory, data.NrOfTrajectories-1, data.NrOfTrajectories-1),
		NrOfDataPoints: data.NrOfDataPoints, NrOfTrajectories: data.NrOfTrajectories - 1}

	wrist := data.Trajectories[0]

	// copying data without wrist frame
	for trajectoryIndex := 1; trajectoryIndex < data.NrOfTrajectories; trajectoryIndex++ {
		for frameIndex := 0; frameIndex < data.NrOfDataPoints; frameIndex++ {
			r.Trajectories[trajectoryIndex-1].Frame =
				append(r.Trajectories[trajectoryIndex-1].Frame,
					PoseCopy(data.Trajectories[trajectoryIndex].Frame[frameIndex]))
		}
	}

	// translate all frames with respect to wrist frame: ONLY POSITION
	for trajectoryIndex := 0; trajectoryIndex < r.NrOfTrajectories; trajectoryIndex++ {
		for frameIndex := 0; frameIndex < r.NrOfDataPoints; frameIndex++ {
			origPosition := r.Trajectories[trajectoryIndex].Frame[frameIndex].Position
			wristPosition := wrist.Frame[frameIndex].Position
			diff := P3DSub(origPosition, wristPosition)
			r.Trajectories[trajectoryIndex].Frame[frameIndex].Position = diff
		}
	}

	return r
}

func ConvertSofaStates(filename string, hands, ctrls []*regexp.Regexp, directory *string, convertToWritsFrame bool) {
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
				data = ConvertAngles(data)
				if convertToWritsFrame {
					data = TransformIntoWristFrame(data)
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
						di := data[:][i]
						dj := data[:][j]
						di = di[start:stop]
						dj = dj[start:stop]
						r[i][j] = fmt.Sprintf("%f", stat.Covariance(di, dj, nil))
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
