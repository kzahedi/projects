package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStatesFrameByFrame(input, output string, hands, ctrls []*regexp.Regexp, directory *string) {
	fmt.Println("Converting sofa state files per segment:", input)
	files := ListAllFilesRecursivelyByFilename(*directory, input)

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
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{5, 6},
		{6, 7},
		{7, 8},
		{8, 9},
		{10, 11},
		{11, 12},
		{12, 13},
		{13, 14},
		{15, 16},
		{16, 17},
		{17, 18},
		{18, 19},
		{20, 21},
		{21, 22},
		{22, 23},
		{23, 24},
		{24, 25},
		{25, 26},
		{26, 27},
		{27, 28},
		{28, 29}}

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
