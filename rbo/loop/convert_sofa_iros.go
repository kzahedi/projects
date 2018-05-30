package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStatesIROS(input, output string, hands, ctrls []*regexp.Regexp, directory string) {
	fmt.Println("Converting sofa state files:", input)
	files := ListAllFilesRecursivelyByFilename(directory, input)

	selectedFiles := SelectFiles(files, hands, ctrls)
	bar := pb.StartNew(len(selectedFiles))

	for _, s := range selectedFiles {
		outfile := strings.Replace(s, "raw", "analysis", 1)
		outfile = strings.Replace(outfile, input, output, 1)
		if _, err := os.Stat(outfile); os.IsNotExist(err) {
			data := ReadSofaSates(s) // returns 2d-array of pose
			data = transformIntoWristFramePositionOnly(data)
			CreateDir(outfile)
			WritePositions(outfile, data)
		}
		bar.Increment()
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
