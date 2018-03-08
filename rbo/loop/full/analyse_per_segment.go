package main

import (
	"fmt"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStatesSegment(filename string, hands, ctrls []*regexp.Regexp, directory *string, convertToWristFrame bool) {
	fmt.Println("Converting sofa state files per segment:", filename)
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
				if convertToWristFrame {
					data = transformSegment(data)
				}
				outfile := strings.Replace(s, "raw", "analysis", 1)
				outfile = strings.Replace(outfile, filename, "segment.sofastates.csv", 1)
				CreateDir(outfile)
				WritePositions(outfile, data)
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
	r := Data{Trajectories: make([]Trajectory, 6, 6), NrOfDataPoints: data.NrOfDataPoints, NrOfTrajectories: 6}

	indices := [][]int{{24, 29}, {20, 24}, {15, 19}, {10, 14}, {5, 9}, {0, 4}}

	for i, v := range indices {
		root := data.Trajectories[v[0]]
		tip := data.Trajectories[v[1]]
		r.Trajectories[5-i].Frame = make([]Pose, data.NrOfDataPoints, data.NrOfDataPoints)
		for j := 0; j < data.NrOfDataPoints; j++ {
			r.Trajectories[5-i].Frame[j] = PoseSub(tip.Frame[j], root.Frame[j])
		}
	}

	return r
}
