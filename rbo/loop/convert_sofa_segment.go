package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStatesSegment(input, output string, hands, ctrls []*regexp.Regexp, directory string) {
	fmt.Println("Converting sofa state files:", input)
	files := ListAllFilesRecursivelyByFilename(directory, input)

	selectedFiles := SelectFiles(files, hands, ctrls)
	bar := pb.StartNew(len(selectedFiles))

	for _, s := range selectedFiles {
		outfile := strings.Replace(s, "raw", "analysis", 1)
		outfile = strings.Replace(outfile, input, output, 1)
		if _, err := os.Stat(outfile); os.IsNotExist(err) {
			data := ReadSofaSates(s) // returns 6d-array of pose
			data = transformSegment(data)
			CreateDir(outfile)
			WritePositions(outfile, data)
		}
		bar.Increment()
	}

	bar.Finish()
}

func transformSegment(data Data) Data {
	// data is the original data from hand.sofastates.txt, including full wrist frame data (trajectory 0)
	// index finder  indices = {1,2,3,4,5}
	// middle finger indices = {6,7,8,9,10}
	// ring finger   indices = {11,12,13,14,15}
	// pinky finger  indices = {16,17,18,19,20}
	// palm          indices = {21,22,23,24,25}
	// thumb         indices = {26,27,28,29,30}
	// we want to have two points per segment ->
	// index finder  coordinate frame 5 in local coordinates of frame 3, then
	// index finder  coordinate frame 3 in local coordinates of frame 1
	// middle finder coordinate frame 10 in local coordinates of frame 8, then
	// middle finder coordinate frame 8 in local coordinates of frame 6
	// ring finder   coordinate frame 15 in local coordinates of frame 13, then
	// ring finder   coordinate frame 13 in local coordinates of frame 11
	// pinky finder  coordinate frame 20 in local coordinates of frame 18, then
	// pinky finder  coordinate frame 18 in local coordinates of frame 16
	// palm          coordinate frame 25 in local coordinate of frame 23, then
	// palm          coordinate frame 23 in local coordinate of frame 21
	// thumb         coordinate frame 30 in local coordinate of frame 28, then
	// thumb         coordinate frame 28 in local coordinate of frame 26

	r := Data{Trajectories: make([]Trajectory, 12, 12), NrOfDataPoints: data.NrOfDataPoints, NrOfTrajectories: 12}

	// each finger has two segments so that correlations can be calculated
	indices := [][]int{
		{3, 5},   // index finger  - second half
		{1, 3},   // index finger  - first half
		{8, 10},  // middle finger - second half
		{6, 8},   // middle finger - first half
		{13, 15}, // ring finger	 - second half
		{11, 13}, // ring finger	 - first half
		{18, 20}, // pinky finger  - second half
		{16, 18}, // pinky finger  - first half
		{23, 25}, // palm          - second half
		{21, 23}, // palm		 	     - first half
		{28, 30}, // thumb				 - second half
		{26, 28}} // thumb				 - first half

	for i, v := range indices {
		root := data.Trajectories[v[0]]
		tip := data.Trajectories[v[1]]
		r.Trajectories[i].Frame = make([]Pose, data.NrOfDataPoints, data.NrOfDataPoints)
		for j := 0; j < data.NrOfDataPoints; j++ {
			r.Trajectories[i].Frame[j] = PoseSub(tip.Frame[j], root.Frame[j])
		}
	}

	return r
}
