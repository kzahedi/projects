package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStatesSegment(input, output string, hands, ctrls []*regexp.Regexp, directory *string) {
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
				outfile := strings.Replace(s, "raw", "analysis", 1)
				outfile = strings.Replace(outfile, input, output, 1)
				if _, err := os.Stat(outfile); os.IsNotExist(err) {
					data := ReadSofaSates(s) // returns 2d-array of pose
					data = transformSegment(data)
					CreateDir(outfile)
					WritePositions(outfile, data)
				}
				bar.Increment()
			}
		}
	}
	bar.Finish()
}

func transformSegment(data Data) Data {
	// 1 = index finder tip to root
	// 2 = middle finger tip to root
	// 3 = ring finger tip to root
	// 4 = pinky finger tip to root
	// 5 = palm finger tip to root
	// 6 = thumb finger tip to root
	r := Data{Trajectories: make([]Trajectory, 12, 12), NrOfDataPoints: data.NrOfDataPoints, NrOfTrajectories: 12}

	// each finger has two segments so that correlations can be calculated
	indices := [][]int{
		{0, 2},   // index finger - first half
		{2, 4},   // index finger - second half
		{5, 7},   // middle finger - first half
		{7, 9},   // middle finger - second half
		{10, 12}, // ring finger - first half
		{12, 14}, // ring finger - second half
		{15, 17}, // pinky finger - first half
		{17, 19}, // pinky finger - second half
		{20, 22}, // palm - first half
		{22, 24}, // palm - second half
		{24, 26}, // thumb - first half
		{26, 29}} // thumb - second half

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
