package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/gonum/matrix/mat64"
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

func ConvertSofaStates(filename string, hand, controller *regexp.Regexp, directory *string, convertToWritsFrame bool) {
	files := ListAllFilesRecursivelyByFilename(*directory, filename)
	rbohand2Files := Select(files, *hand)
	rbohand2Files = Select(rbohand2Files, *controller)

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

func CalculateDifferenceBehaviour(grasp, prescriptive, controller *regexp.Regexp, directory *string) {
	filename := "hand.sofastates.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	rbohand2Grasps := Select(files, *grasp)
	rbohand2Grasps = Select(rbohand2Grasps, *controller)

	// we only take the first prescriptive, because they are all identical (should be)
	rbohand2Prescriptives := Select(files, *prescriptive)
	rbohand2Prescriptives = Select(rbohand2Prescriptives, *controller)
	prescritiveData := ReadCSVToData(rbohand2Prescriptives[0])

	for _, s := range rbohand2Grasps {
		data := ReadCSVToData(s) // returns 2d-array of pose
		diff := calculateDifferencePositionOnly(data, prescritiveData)
		output := strings.Replace(s, filename, fmt.Sprintf("difference.%s", filename), 1)
		WritePositions(output, diff)
	}
}
