package main

import (
	"fmt"
	"math"
	"regexp"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func generateFingerTipMinMaxBins(hands, ctrls []*regexp.Regexp, directory *string, wBins int) ([]float64, []float64, []int) {
	fmt.Println("Getting min/max/bin values for hand")
	handFilename := "hand.sofastates.csv"
	handFiles := ListAllFilesRecursivelyByFilename(*directory, handFilename)

	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(handFiles, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}

	bar := pb.StartNew(iterations)

	handMin := make([]float64, 3, 3) // x, y, z -> we bin per dimension
	handMax := make([]float64, 3, 3) // x, y, z -> we bin per dimension

	first := true

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			behaviours := Select(handFiles, *hand)
			behaviours = Select(behaviours, *ctrl)

			for _, s := range behaviours {
				data := ReadCSVToFloat(s)
				rows := len(data)
				cols := len(data[0])
				for row := 0; row < rows; row++ {
					for col := 0; col < cols; col++ {
						if first {
							handMin[col%3] = data[0][col]
							handMax[col%3] = data[0][col]
						}
						handMax[col%3] = math.Max(data[row][col], handMax[col%3])
						handMin[col%3] = math.Min(data[row][col], handMin[col%3])
					}
					first = false
				}
				bar.Increment()
			}
		}
	}
	bar.Finish()

	minFingerTip := make([]float64, 12, 12)
	maxFingerTip := make([]float64, 12, 12)
	binsFingerTip := make([]int, 12, 12)

	for i := 0; i < 11; i++ {
		minFingerTip[i] = handMin[i%3]
		maxFingerTip[i] = handMax[i%3]
		binsFingerTip[i] = wBins
	}

	return minFingerTip, maxFingerTip, binsFingerTip
}

func generateControllerMinMaxBins(hands, ctrls []*regexp.Regexp, directory *string, aBins int) ([]float64, []float64, []int) {
	fmt.Println("Getting min/max/bin values for controller")
	ctrlFilename := "control.states.csv"
	ctrlFiles := ListAllFilesRecursivelyByFilename(*directory, ctrlFilename)

	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(ctrlFiles, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}

	bar := pb.StartNew(iterations)

	ctrlMin := make([]float64, 6, 6) // ctrl pressure states
	ctrlMax := make([]float64, 6, 6)

	first := true

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			controls := Select(ctrlFiles, *hand)
			controls = Select(controls, *ctrl)

			for _, s := range controls {
				data := ReadControlFile(s)
				rows := len(data)
				for row := 0; row < rows; row++ {
					if first {
						for c := 0; c < 6; c++ {
							ctrlMin[c] = data[row][c+8]
							ctrlMax[c] = data[row][c+8]
						}
					}
					for c := 0; c < 6; c++ {
						ctrlMin[c] = math.Min(data[row][c+8], ctrlMin[c])
						ctrlMax[c] = math.Max(data[row][c+8], ctrlMax[c])
					}
				}
				first = false
				bar.Increment()
			}
		}
	}
	bar.Finish()

	ctrlBins := make([]int, 6, 6)
	for i := 0; i < 6; i++ {
		ctrlBins[i] = aBins
	}

	return ctrlMin, ctrlMin, ctrlBins
}
