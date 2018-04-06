package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func CalculateGraspDistance(hands, ctrls []*regexp.Regexp, directory *string, lastNSteps, cutOff int, results Results) Results {
	objectFilename := "obstacle.sofastates.csv"
	handFilename := "hand.sofastates.csv"
	handFiles := ListAllFilesRecursivelyByFilename(*directory, handFilename)
	fcutOff := float64(cutOff)

	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(handFiles, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}

	fmt.Println("Calculating Grasp Distances")
	bar := pb.StartNew(iterations)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			behaviours := Select(handFiles, *hand)
			behaviours = Select(behaviours, *ctrl)

			for _, s := range behaviours {
				handData := ReadCSVToFloat(s)
				handData = extractMiddleFingerRoot(handData)
				objF := strings.Replace(s, handFilename, objectFilename, -1)
				objectData := ReadCSVToFloat(objF)

				gd := calculateGD(handData, objectData, lastNSteps)

				key := GetKey(s)
				v := results[key]
				v.GraspDistance = math.Min(gd, fcutOff)
				results[key] = v
				bar.Increment()
			}
		}
	}
	bar.Finish()
	return results
}

func extractMiddleFingerRoot(data [][]float64) [][]float64 {
	r := make([][]float64, len(data), len(data))

	for i := 0; i < len(data); i++ {
		r[i] = make([]float64, 3, 3)
		// hand.sofastates.csv is the original data, therefore 6 and not 5
		r[i][0] = data[i][6*3+0]
		r[i][1] = data[i][6*3+1]
		r[i][2] = data[i][6*3+2]
	}

	return r
}

func calculateGD(a, b [][]float64, n int) float64 {
	l := len(a)
	r := 0.0

	for i := l - n - 1; i < l; i++ {
		r += dist(a[i], b[i])
	}

	return r / float64(n)
}

func dist(a, b []float64) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	dz := a[2] - b[2]
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
