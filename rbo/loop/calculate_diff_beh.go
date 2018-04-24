package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func CalculateDifferenceBehaviour(input, output string, hands, ctrls []*regexp.Regexp, directory *string) {
	fmt.Println("Calculating difference behaviour")
	files := ListAllFilesRecursivelyByFilename(*directory, input)

	selectedFiles := SelectFiles(files, hands, ctrls)
	bar := pb.StartNew(len(selectedFiles))

	for _, s := range selectedFiles {
		prescritiveFilename := convertFilenameToPrescriptive(s)
		prescritiveData := ReadCSVToData(prescritiveFilename)
		outfile := strings.Replace(s, input, output, 1)
		if _, err := os.Stat(outfile); os.IsNotExist(err) {
			data := ReadCSVToData(s) // returns 2d-array of pose
			diff := calculateDifferencePositionOnly(data, prescritiveData)
			WritePositions(outfile, diff)
		}
		bar.Increment()
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
			// diff := PoseSub(g, p)
			var diff Pose
			diff.Position.X = g.Position.X - p.Position.X
			diff.Position.Y = g.Position.Y - p.Position.Y
			diff.Position.Z = g.Position.Z - p.Position.Z
			r.Trajectories[trajectoryIndex].Frame[frameIndex] = diff
		}
	}
	return r
}

func convertFilenameToPrescriptive(input string) (output string) {
	directory := regexp.MustCompile("/(rbohand[-a-zA-Z0-9]+)/")
	position := regexp.MustCompile("(-*[0-9]+.[0-9]+_-*[0-9]+.[0-9]+_-*[0-9]+.[0-9])")
	object := regexp.MustCompile("(object[a-zA-Z]+)")
	output = directory.ReplaceAllString(input, "/$1-prescriptive/")
	output = position.ReplaceAllString(output, "0.0_0.0_0.0")
	output = object.ReplaceAllString(output, "objectcylinder")
	return
}
