package main

import (
	"fmt"
	"math"
	"sort"
)

func CalculateInteretingClusters(data Results, graspDistanceCutoff, percentage float64, k int) Results {

	maxMCW := -1.0
	minMCW := -1.0

	minGraspDistance := 0.0
	maxGraspDistance := 0.0

	for _, value := range data {

		if value.ClusteredByTSE {
			mcw := value.MC_W
			gd := value.GraspDistance

			if gd > graspDistanceCutoff {
				gd = graspDistanceCutoff
			}

			if maxMCW < 0.0 {
				maxMCW = mcw
				minMCW = mcw
				minGraspDistance = gd
				maxGraspDistance = gd
			}

			if maxMCW < mcw {
				maxMCW = mcw
			}

			if minMCW > mcw {
				minMCW = mcw
			}

			if minGraspDistance > gd {
				minGraspDistance = gd
			}
			if maxGraspDistance < gd {
				maxGraspDistance = gd
			}
		}
	}

	mcwDiff := maxMCW * percentage
	gdDiff := maxGraspDistance * percentage

	stupidMCW := minMCW + mcwDiff
	intelligentMCW := maxMCW - mcwDiff

	stupidGraspDistance := maxGraspDistance - gdDiff
	intelligentGraspDistance := minGraspDistance + gdDiff

	nrOfIntelligent := 0
	nrOfStupid := 0
	nrOfNone := 0
	for key, value := range data {
		if value.ClusteredByTSE {
			mcw := value.MC_W
			gd := value.GraspDistance

			value.Stupid = false
			value.Intelligent = false

			if mcw < stupidMCW && gd > stupidGraspDistance {
				value.Stupid = true
				value.Intelligent = false
				nrOfStupid++
			} else if mcw > intelligentMCW && gd < intelligentGraspDistance {
				value.Stupid = false
				value.Intelligent = true
				nrOfIntelligent++
			} else {
				nrOfNone++
			}
			data[key] = value
		}
	}
	fmt.Println(fmt.Sprintf("Intelligent, MC > %.5f and Grasp Distance < %.5f", intelligentMCW, intelligentGraspDistance))
	fmt.Println(fmt.Sprintf("Stupid,      MC < %.5f and Grasp Distance > %.5f", stupidMCW, stupidGraspDistance))

	fmt.Println(fmt.Sprintf("Intelligent %d", nrOfIntelligent))
	fmt.Println(fmt.Sprintf("Stupid      %d", nrOfStupid))
	fmt.Println(fmt.Sprintf("None        %d", nrOfNone))
	fmt.Println(fmt.Sprintf("Sum         %d", nrOfIntelligent+nrOfStupid+nrOfNone))

	intelligent := SelectIntelligent(data)
	stupid := SelectStupid(data)

	intelligent = getMeaningfulCluster(intelligent, k)
	stupid = getMeaningfulCluster(stupid, k)

	for key, value := range data {
		for ikey, ivalue := range intelligent {
			if key == ikey {
				value.SelectedForAnalysis = ivalue.SelectedForAnalysis
				data[key] = value
			}
		}
		for skey, svalue := range stupid {
			if key == skey {
				value.SelectedForAnalysis = svalue.SelectedForAnalysis
				data[key] = value
			}
		}
	}

	return data
}

func getMeaningfulCluster(data Results, k int) Results {
	for key1, value1 := range data {
		distances := make([]float64, len(data), len(data))

		index := 0
		for _, value2 := range data {
			dx := value1.PosX - value2.PosX
			dy := value1.PosY - value2.PosY
			dist := math.Sqrt(dx*dx + dy*dy)
			distances[index] = dist
			index++
		}
		sort.Slice(distances, func(a, b int) bool {
			return distances[a] < distances[b]
		})
		value1.Distance = distances[k]
		value1.SelectedForAnalysis = true // will be set false below
		data[key1] = value1
	}

	for key1, value1 := range data {
		for _, value2 := range data {
			if value1.Distance > value2.Distance {
				value1.SelectedForAnalysis = false
				data[key1] = value1
			}
		}
	}

	// var centerKey string
	var center Result

	for _, value := range data {
		if value.SelectedForAnalysis == true {
			// centerKey = key
			center = value
		}
	}

	for key, value := range data {
		dx := center.PosX - value.PosX
		dy := center.PosY - value.PosY
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < center.Distance {
			value.SelectedForAnalysis = true
			data[key] = value
		}
	}

	return data
}

func SelectIntelligent(data Results) Results {
	r := make(Results)

	for key, value := range data {
		if value.Intelligent {
			r[key] = value
		}
	}
	return r
}

func SelectStupid(data Results) Results {
	r := make(Results)

	for key, value := range data {
		if value.Stupid {
			r[key] = value
		}
	}
	return r
}
