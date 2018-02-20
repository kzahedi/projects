package main

import (
	"fmt"
	"math"
	"path"
	"regexp"
	"strings"

	"github.com/kzahedi/goent/dh"
	"github.com/kzahedi/goent/discrete"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func generateFingerTipMinMaxBins(hands, ctrls []*regexp.Regexp, directory *string, wBins int) ([]float64, []float64, []int) {
	fmt.Println("  Getting min/max/bin values for hand")
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

	minFingerTip := make([]float64, 12, 12)
	maxFingerTip := make([]float64, 12, 12)
	binsFingerTip := make([]int, 12, 12)

	for i := 0; i < 11; i++ {
		minFingerTip[i] = handMin[i%3]
		maxFingerTip[i] = handMax[i%3]
		binsFingerTip[i] = wBins
	}

	bar.Finish()

	return minFingerTip, maxFingerTip, binsFingerTip
}

func generateControllerMinMaxBins(hands, ctrls []*regexp.Regexp, directory *string, aBins int) ([]float64, []float64, []int) {
	fmt.Println("  Getting min/max/bin values for controller")
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

	ctrlBins := make([]int, 6, 6)
	for i := 0; i < 6; i++ {
		ctrlBins[i] = aBins
	}

	bar.Finish()

	return ctrlMin, ctrlMin, ctrlBins
}

func CalculateMCW(hands, ctrls []*regexp.Regexp, directory *string, wBins, aBins int, results *Results) {
	fmt.Println("Calculating MC_W (discrete)")
	ctrlFilename := "control.states.csv"
	handFilename := "hand.sofastates.csv"
	handFiles := ListAllFilesRecursivelyByFilename(*directory, handFilename)

	handMin, handMax, handBins := generateFingerTipMinMaxBins(hands, ctrls, directory, wBins)
	ctrlMin, ctrlMax, ctrlBins := generateControllerMinMaxBins(hands, ctrls, directory, wBins)

	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(handFiles, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}

	fmt.Println("  Calculating MC_W on fingertips")
	bar := pb.StartNew(iterations)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			behaviours := Select(handFiles, *hand)
			behaviours = Select(behaviours, *ctrl)

			for _, s := range behaviours {
				ftd := ReadCSVToFloat(s)
				fingerTipData := extractFingerTipData(ftd)
				discretisedFingerTipData := dh.Discrestise(fingerTipData, handBins, handMin, handMax)
				univariateFingerTipData := dh.MakeUnivariateRelabelled(discretisedFingerTipData, handBins)

				c := strings.Replace(s, "analysis", "raw", -1)
				c = strings.Replace(c, handFilename, ctrlFilename, -1)
				ctd := ReadControlFile(c)
				ctrlData := extractControllerData(ctd)
				discretisedCtrlData := dh.Discrestise(ctrlData, ctrlBins, ctrlMin, ctrlMax)
				univariateCtrlData := dh.MakeUnivariateRelabelled(discretisedCtrlData, ctrlBins)

				w2w1a1 := mergeDataForMCW(univariateFingerTipData, univariateCtrlData)
				pw2w1a1 := discrete.Emperical3D(w2w1a1)
				mc_w := discrete.MorphologicalComputationW(pw2w1a1)

				key := strings.Replace(path.Dir(s), "/analysis", "", -1)

				if v, found := (*results)[key]; found == false {
					r := Result{MC_W: mc_w, GraspDistance: 0.0, Point: []float64{0.0, 0.0}}
					(*results)[key] = r
				} else {
					v.MC_W = mc_w
					(*results)[key] = v
				}
				bar.Increment()
			}
		}
	}

	// output := strings.Replace(s, filename, "mc_w.csv", 1)
	// WriteCSVFloat(output, data)
	bar.Finish()
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

		r[row][0] = data[row][4*3+0] // index finger x
		r[row][1] = data[row][4*3+1] // index finger y
		r[row][2] = data[row][4*3+2] // index finger z

		r[row][3] = data[row][9*3+0] // index finger x
		r[row][4] = data[row][9*3+1] // index finger y
		r[row][5] = data[row][9*3+2] // index finger z

		r[row][6] = data[row][14*3+0] // index finger x
		r[row][7] = data[row][14*3+1] // index finger y
		r[row][8] = data[row][14*3+2] // index finger z

		r[row][9] = data[row][19*3+0]  // index finger x
		r[row][10] = data[row][19*3+1] // index finger y
		r[row][11] = data[row][19*3+2] // index finger z
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

func CalculateGraspDistance(hands, ctrls []*regexp.Regexp, directory *string, lastNSteps, cutOff int, results *Results) {
	fmt.Println("Calculating MC_W (discrete)")
	objectFilename := "obstacle.sofastates.csv"
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

	fmt.Println("  Calculating Grasp Distances")
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

				key := strings.Replace(path.Dir(s), "/analysis", "", -1)

				if v, found := (*results)[key]; found == false {
					r := Result{MC_W: 0.0, GraspDistance: gd, Point: []float64{0.0, 0.0}}
					(*results)[key] = r
				} else {
					v.GraspDistance = math.Min(gd, cutOff)
					(*results)[key] = v
				}
				bar.Increment()
			}
		}
	}

	// output := strings.Replace(s, filename, "mc_w.csv", 1)
	// WriteCSVFloat(output, data)
	bar.Finish()
}

func extractMiddleFingerRoot(data [][]float64) [][]float64 {
	r := make([][]float64, len(data), len(data))

	for i := 0; i < len(data); i++ {
		r[i] = make([]float64, 3, 3)
		r[i][0] = data[i][5*3+0]
		r[i][1] = data[i][5*3+1]
		r[i][2] = data[i][5*3+2]
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
