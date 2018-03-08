package main

import (
	"fmt"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStatesFrameToFrame(filename string, hands, ctrls []*regexp.Regexp, directory *string, convertToWristFrame bool) {
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
					data = transformFrameByFrame(data)
				}
				outfile := strings.Replace(s, "raw", "analysis", 1)
				outfile = strings.Replace(outfile, filename, "frame.by.frame.sofastates.csv", 1)
				CreateDir(outfile)
				WritePositions(outfile, data)
				bar.Increment()
			}
		}
	}
	bar.Finish()
}

func transformFrameByFrame(data Data) Data {
	n := data.NrOfTrajectories - 1
	r := Data{Trajectories: make([]Trajectory, n, n), NrOfDataPoints: data.NrOfDataPoints, NrOfTrajectories: 6}

	indices := [][]int{
		{29, 30}, // thumb and palm
		{28, 29},
		{27, 28},
		{26, 27},
		{25, 30},
		{24, 25},
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
		{5, 6}, // thumb
		{4, 5},
		{3, 4},
		{2, 3},
		{0, 1}}

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

func CalculateDifferenceBehaviourFrameToFrame(hands, ctrls []*regexp.Regexp, prescriptive *regexp.Regexp, directory *string) {
	fmt.Println("Calculating difference behaviour for segments")
	filename := "frame.to.frame.sofastates.csv"
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
				diff := calculateDifferencePositionOnlyFrameByFrame(data, prescritiveData)
				output := strings.Replace(s, filename, fmt.Sprintf("difference.%s", filename), 1)
				WritePositions(output, diff)
				bar.Increment()
			}
		}
	}
	bar.Finish()
}

func calculateDifferencePositionOnlyFrameByFrame(grasp, prescriptive Data) Data {
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
