package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/kzahedi/goent/dh"
	"github.com/kzahedi/goent/discrete"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func CalculateMCW(hands, ctrls []*regexp.Regexp, directory *string, wBins, aBins int, results Results) Results {
	fmt.Println("Calculating MC_W")
	ctrlFilename := "control.states.csv"
	handFilename := "hand.sofastates.csv"
	handFiles := ListAllFilesRecursivelyByFilename(*directory, handFilename)

	handMin, handMax, handBins := generateFingerTipMinMaxBins(hands, ctrls, directory, wBins)
	ctrlMin, ctrlMax, ctrlBins := generateControllerMinMaxBins(hands, ctrls, directory, wBins)

	selectedFiles := SelectFiles(handFiles, hands, ctrls)
	bar := pb.StartNew(len(selectedFiles))

	for _, s := range selectedFiles {
		ftd := ReadCSVToFloat(s)
		fingerTipData := extractFingerTipData(ftd)
		discretisedFingerTipData := dh.Discretise(fingerTipData, handBins, handMin, handMax)
		univariateFingerTipData := dh.MakeUnivariateRelabelled(discretisedFingerTipData, handBins)

		c := strings.Replace(s, "analysis", "raw", -1)
		c = strings.Replace(c, handFilename, ctrlFilename, -1)
		ctd := ReadControlFile(c)
		ctrlData := extractControllerData(ctd)
		discretisedCtrlData := dh.Discretise(ctrlData, ctrlBins, ctrlMin, ctrlMax)
		univariateCtrlData := dh.MakeUnivariateRelabelled(discretisedCtrlData, ctrlBins)

		w2w1a1 := mergeDataForMCW(univariateFingerTipData, univariateCtrlData)
		pw2w1a1 := discrete.Emperical3D(w2w1a1)
		mc_w := discrete.MorphologicalComputationW(pw2w1a1)

		key := GetKey(s)

		v := results[key]
		v.MC_W = mc_w
		results[key] = v
		bar.Increment()
	}

	bar.Finish()

	return results
}

// takes world states W and action states A and returns the sequence W',W,A
func mergeDataForMCW(w, a []int) [][]int {
	w2w1a1 := make([][]int, len(w)-1, len(w)-1)
	for i := 0; i < len(w)-1; i++ {
		w2w1a1[i] = make([]int, 3, 3)
		w2w1a1[i][0] = w[i+1]
		w2w1a1[i][1] = w[i]
		w2w1a1[i][2] = a[i]
	}
	return w2w1a1
}

func extractFingerTipData(data [][]float64) [][]float64 {
	r := make([][]float64, len(data), len(data))

	for row := 0; row < len(data); row++ {
		r[row] = make([]float64, 4*3, 4*3)

		r[row][0] = data[row][5*3+0] // index finger x
		r[row][1] = data[row][5*3+1] // index finger y
		r[row][2] = data[row][5*3+2] // index finger z

		r[row][3] = data[row][10*3+0] // index finger x
		r[row][4] = data[row][10*3+1] // index finger y
		r[row][5] = data[row][10*3+2] // index finger z

		r[row][6] = data[row][15*3+0] // index finger x
		r[row][7] = data[row][15*3+1] // index finger y
		r[row][8] = data[row][15*3+2] // index finger z

		r[row][9] = data[row][20*3+0]  // index finger x
		r[row][10] = data[row][20*3+1] // index finger y
		r[row][11] = data[row][20*3+2] // index finger z
	}

	return r
}

func extractControllerData(data [][]float64) [][]float64 {
	r := make([][]float64, len(data), len(data))

	for row := 0; row < len(data); row++ {
		r[row] = data[row][8:]
	}

	return r
}
