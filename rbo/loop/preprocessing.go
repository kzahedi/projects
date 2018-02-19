package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/gonum/matrix/mat64"
	"github.com/gonum/stat"
	"github.com/kzahedi/goent/dh"
	"github.com/kzahedi/goent/discrete"
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

func CalculateMCW(hands, ctrls []*regexp.Regexp, directory *string, wBins, aBins int) {
	fmt.Println("Calculating MC_W (discrete)")
	handFilename := "hand.sofastates.csv"
	handFiles := ListAllFilesRecursivelyByFilename(*directory, handFilename)

	ctrlFilename := "control.states.csv"
	ctrlFiles := ListAllFilesRecursivelyByFilename(*directory, ctrlFilename)

	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(handFiles, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}

	fmt.Println("  Getting min/max values for x,y,z")
	bar := pb.StartNew(iterations)

	handMin := make([]float64, 3, 3) // x, y, z -> we bin per dimension
	handMax := make([]float64, 3, 3) // x, y, z -> we bin per dimension

	first := true

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			behaviours := Select(handFiles, *hand)
			behaviours = Select(behaviours, *ctrl)

			for _, s := range behaviours {
				data := ReadCSVToFloat(s)
				rows := len(data)
				cols := len(data[0])
				for row := 0; row < rows; row++ {
					for col := 0; col < cols; col++ {
						if first {
							handMin[col%3] = data[0][col]
							handMax[col%3] = data[0][col]
						}
						handMax[col%3] = math.Max(data[row][col], handMax[col%3])
						handMin[col%3] = math.Min(data[row][col], handMin[col%3])
					}
					first = false
				}
				bar.Increment()
			}
		}
	}

	fmt.Println("  Getting min/max values for ctrl")
	bar = pb.StartNew(iterations)

	ctrlMin := make([]float64, 6, 6) // ctrl pressure states
	ctrlMax := make([]float64, 6, 6)

	first = true

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			controls := Select(ctrlFiles, *hand)
			controls = Select(controls, *ctrl)

			for _, s := range controls {
				data := ReadControlFile(s)
				rows := len(data)
				for row := 0; row < rows; row++ {
					if first {
						for c := 0; c < 6; c++ {
							ctrlMin[c] = data[row][c+8]
							ctrlMax[c] = data[row][c+8]
						}
					}
					for c := 0; c < 6; c++ {
						ctrlMin[c] = math.Min(data[row][c+8], ctrlMin[c])
						ctrlMax[c] = math.Max(data[row][c+8], ctrlMax[c])
					}
				}
				first = false
			}
			bar.Increment()
		}
	}

	fmt.Println("  Calculating MC_W on fingertips")
	bar = pb.StartNew(iterations)

	minFingerTip := make([]float64, 12, 12)
	maxFingerTip := make([]float64, 12, 12)
	binsFingerTip := make([]int, 12, 12)

	for i := 0; i < 11; i++ {
		minFingerTip[i] = handMin[i%3]
		maxFingerTip[i] = handMax[i%3]
		binsFingerTip[i] = wBins
	}

	binsCtrl := make([]int, 6, 6)
	for i := 0; i < 6; i++ {
		binsCtrl[i] = aBins
	}

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			behaviours := Select(handFiles, *hand)
			behaviours = Select(behaviours, *ctrl)

			for _, s := range behaviours {
				ftd := ReadCSVToFloat(s)
				fingerTipData := extractFingerTipData(ftd)
				discretisedFingerTipData := dh.Discrestise(fingerTipData, binsFingerTip, minFingerTip, maxFingerTip)
				univariateFingerTipData := dh.MakeUnivariateRelabelled(discretisedFingerTipData, binsFingerTip)

				c := strings.Replace(s, "analysis", "raw", -1)
				c = strings.Replace(c, handFilename, ctrlFilename, -1)
				ctd := ReadCSVToFloat(c)
				ctrlData := extractControllerData(ctd)
				discretisedCtrlData := dh.Discrestise(ctrlData, binsCtrl, ctrlMin, ctrlMax)
				univariateCtrlData := dh.MakeUnivariateRelabelled(discretisedCtrlData, binsCtrl)

				w2w1a1 := mergeDataForMCW(univariateFingerTipData, univariateCtrlData)
				pw2w1a1 := discrete.Emperical3D(w2w1a1)
				mc_w := discrete.MorphologicalComputationW(pw2w1a1)
				fmt.Println(mc_w)
			}
			bar.Increment()
		}
	}

	// output := strings.Replace(s, filename, "mc_w.csv", 1)
	// WriteCSVFloat(output, data)
	bar.Finish()
}

// takes world states W and action states A and returns the sequence W',W,A
func mergeDataForMCW(w, a []int) [][]int {
	w2w1a1 := make([][]int, len(w)-1, len(w)-1)
	for i := 0; i < len(w)-1; i++ {
		w2w1a1[i] = make([]int, 3, 3)
		w2w1a1[i][0] = w[i+1]
		w2w1a1[i][1] = w[i]
		w2w1a1[i][2] = a[i]
	}
	return w2w1a1
}

func extractFingerTipData(data [][]float64) [][]float64 {
	r := make([][]float64, len(data), len(data))

	for row := 0; row < len(data); row++ {
		r[row] = make([]float64, 4*3, 4*3)

		r[row][0] = data[row][4*3+0] // index finger x
		r[row][1] = data[row][4*3+1] // index finger y
		r[row][2] = data[row][4*3+2] // index finger z

		r[row][3] = data[row][9*3+0] // index finger x
		r[row][4] = data[row][9*3+1] // index finger y
		r[row][5] = data[row][9*3+2] // index finger z

		r[row][6] = data[row][14*3+0] // index finger x
		r[row][7] = data[row][14*3+1] // index finger y
		r[row][8] = data[row][14*3+2] // index finger z

		r[row][9] = data[row][19*3+0]  // index finger x
		r[row][10] = data[row][19*3+1] // index finger y
		r[row][11] = data[row][19*3+2] // index finger z
	}

	return r
}

func extractControllerData(data [][]float64) [][]float64 {
	r := make([][]float64, len(data), len(data))

	for row := 0; row < len(data); row++ {
		r[row] = data[row][8:]
	}

	return r
}
