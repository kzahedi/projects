package main

import (
	"fmt"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStatesIROS(input, output string, hands, ctrls []*regexp.Regexp, directory *string, convertToWritsFrame bool) {
	fmt.Println("Converting sofa state files:", input)
	files := ListAllFilesRecursivelyByFilename(*directory, input)
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
				if convertToWritsFrame {
					data = transformIntoWristFramePositionOnly(data)
				}
				outfile := strings.Replace(s, "raw", "analysis", 1)
				outfile = strings.Replace(outfile, input, output, 1)
				CreateDir(outfile)
				WritePositions(outfile, data)
				bar.Increment()
			}
		}
	}
	bar.Finish()
}

// tested
func calculateDifferencePositionOnlyIROS(grasp, prescriptive Data) Data {
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

func CalculateDifferenceBehaviourIROS(input, output string, hands, ctrls []*regexp.Regexp, prescriptive *regexp.Regexp, directory *string) {
	fmt.Println("Calculating difference behaviour")
	files := ListAllFilesRecursivelyByFilename(*directory, input)

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
				diff := calculateDifferencePositionOnlyIROS(data, prescritiveData)
				output := strings.Replace(s, input, output, 1)
				WritePositions(output, diff)
				bar.Increment()
			}
		}
	}
	bar.Finish()
}

// transformIntoWristFramePositionOnly transforms all coordinate frames
func transformIntoWristFramePositionOnly(data Data) Data {
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
