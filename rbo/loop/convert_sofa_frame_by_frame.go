package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStatesFrameByFrame(input, output string, hands, ctrls []*regexp.Regexp, directory string) {
	fmt.Println("Converting sofa state files per segment:", input)
	files := ListAllFilesRecursivelyByFilename(directory, input)

	selectedFiles := SelectFiles(files, hands, ctrls)
	bar := pb.StartNew(len(selectedFiles))

	for _, s := range selectedFiles {
		outfile := strings.Replace(s, "raw", "analysis", 1)
		outfile = strings.Replace(outfile, input, output, 1)
		if _, err := os.Stat(outfile); os.IsNotExist(err) {
			data := ReadSofaSates(s) // returns 2d-array of pose
			data = transformFrameByFrame(data)
			CreateDir(outfile)
			WritePositions(outfile, data)
		}
		bar.Increment()
	}

	bar.Finish()
}

func transformFrameByFrame(data Data) Data {
	n := data.NrOfTrajectories - 1
	m := data.NrOfDataPoints
	r := Data{Trajectories: make([]Trajectory, n-1, n-1), NrOfDataPoints: m, NrOfTrajectories: n - 1}

	indices := [][]int{
		{29, 30}, // thumb
		{28, 29},
		{27, 28},
		{26, 27},
		{25, 26},
		{24, 25}, // palm
		{23, 24},
		{22, 23},
		{21, 22},
		{0, 21},
		{19, 20}, // pinky finger
		{18, 19},
		{17, 18},
		{16, 17},
		{0, 16},
		{14, 15}, // ring finger
		{13, 14},
		{12, 13},
		{11, 12},
		{0, 11},
		{9, 10}, // middle finger
		{8, 9},
		{7, 8},
		{6, 7},
		{0, 6},
		{5, 6}, // index finger
		{4, 5},
		{3, 4},
		{1, 2},
		{0, 1}}

	for i, v := range indices {
		root := data.Trajectories[v[0]]
		tip := data.Trajectories[v[1]]
		r.Trajectories[i].Frame = make([]Pose, data.NrOfDataPoints, data.NrOfDataPoints)
		for j := 0; j < m; j++ {
			r.Trajectories[i].Frame[j] = PoseSub(tip.Frame[j], root.Frame[j])
		}
	}

	return r
}
