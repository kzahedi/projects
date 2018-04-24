package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/kzahedi/goent/dh"
	"github.com/kzahedi/goent/discrete"
	"github.com/sacado/tsne4go"
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

func CalculateSuccess(hands, ctrls []*regexp.Regexp, directory *string, height float64, results *Results) {
	objectFilename := "obstacle.sofastates.csv"
	objectFiles := ListAllFilesRecursivelyByFilename(*directory, objectFilename)

	osizes := make(map[string]float64)
	osizes["objectcylinder"] = 20.0
	osizes["objectcylinderB"] = 40.0
	osizes["objectbox"] = 35.0
	osizes["objectboxB"] = 20.0
	osizes["objectsphere"] = 35.0
	osizes["objectsphereB"] = 20.0
	osizes["objectegg"] = 35.0
	osizes["objecteggB"] = 20.0

	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(objectFiles, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}

	fmt.Println("Calculating Success")
	bar := pb.StartNew(iterations)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			objects := Select(objectFiles, *hand)
			objects = Select(objects, *ctrl)

			for _, s := range objects {
				data := ReadCSVToFloat(s)
				maxHeight := data[20][1]
				for i := 20; i < len(data); i++ {
					maxHeight = math.Max(maxHeight, data[i][1])
				}

				key := GetKey(s)
				objectName := GetObjectName(s)

				v := (*results)[key]
				v.Successful = ((maxHeight - osizes[objectName]) > height)
				(*results)[key] = v

				bar.Increment()
			}
		}
	}
	bar.Finish()
}

func CalculateMCW(hands, ctrls []*regexp.Regexp, directory *string, wBins, aBins int, results *Results) {
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

	fmt.Println("Calculating MC_W on fingertips")
	bar := pb.StartNew(iterations)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			behaviours := Select(handFiles, *hand)
			behaviours = Select(behaviours, *ctrl)

			for _, s := range behaviours {
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

				v := (*results)[key]
				v.MC_W = mc_w
				(*results)[key] = v
				bar.Increment()
			}
		}
	}
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
				v := (*results)[key]
				v.GraspDistance = math.Min(gd, fcutOff)
				(*results)[key] = v
				bar.Increment()
			}
		}
	}
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

func CalculateTSNE(hand, controller []*regexp.Regexp, directory *string, iterations int, successfulOnly bool, results *Results) {
	fmt.Println("Calculating TSNE")
	filename := "covariance.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	covariances := Select(files, *hand)
	covariances = Select(covariances, *controller)

	var selected []string
	if successfulOnly == false {
		selected = covariances
	} else {
		fmt.Println("number:", len(covariances))
		for _, v := range covariances {
			key := GetKey(v)
			elem := (*results)[key]
			if elem.Successful {
				selected = append(selected, v)
			}
		}
		fmt.Println("number:", len(selected))
	}

	var data tsne4go.VectorDistancer
	data = make([][]float64, len(selected), len(selected))
	for i, f := range selected {
		data[i] = ReadCSVToArray(f)
	}

	tsne := tsne4go.New(data, nil)

	WriteCSVFloat("/Users/zahedi/Desktop/data.csv", data)

	bar := pb.StartNew(iterations)

	for i := 0; i < iterations; i++ {
		tsne.Step()
		bar.Increment()
	}
	bar.Finish()

	for i := 0; i < len(selected); i++ {
		key := GetKey(selected[i])
		v := (*results)[key]
		v.Point[0] = tsne.Solution[i][0]
		v.Point[1] = tsne.Solution[i][1]
		v.ClusteredByTSE = true
		(*results)[key] = v
	}
}

func ExtractObjectPosition(results *Results) {
	bar := pb.StartNew(len(*results))

	r := make(map[string]int)

	index := 0

	for key, _ := range *results {
		s := extractPositionString(key)
		if _, ok := r[s]; ok == false {
			r[s] = index
			index++
		}
		bar.Increment()
	}
	bar.Finish()

	for key, value := range *results {
		s := extractPositionString(key)
		value.ObjectPosition = r[s]
		(*results)[key] = value
	}

}

func extractPositionString(in string) string {
	re := regexp.MustCompile("-?[0-9]{1,2}.[0-9]{0,2}_-?[0-9]{1,2}.[0-9]{0,2}_-?[0-9]{1,2}.[0-9]{0,2}")
	return re.FindAllString(in, -1)[0]
}

func extractObjectString(in string) string {
	re := regexp.MustCompile("object[a-zA-Z]+")
	return re.FindAllString(in, -1)[0]
}

func ExtractObjectType(results *Results) {
	bar := pb.StartNew(len(*results))

	r := make(map[string]int)

	index := 0

	for key, _ := range *results {
		s := extractObjectString(key)
		if _, ok := r[s]; ok == false {
			r[s] = index
			index++
		}
		bar.Increment()
	}
	bar.Finish()

	for key, value := range *results {
		s := extractObjectString(key)
		value.ObjectType = r[s]
		(*results)[key] = value
	}
}
