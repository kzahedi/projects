package main

import "fmt"

func CalculateInteretingClusters(data Results, graspDistanceCutoff, percentage float64) Results {

	maxMCW := -1.0
	minMCW := -1.0

	minGraspDistance := 0.0
	maxGraspDistance := 0.0

	for _, value := range data {

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

	fmt.Println(fmt.Sprintf("Intelligent %d", nrOfIntelligent))
	fmt.Println(fmt.Sprintf("Stupid      %d", nrOfStupid))
	fmt.Println(fmt.Sprintf("None        %d", nrOfNone))
	fmt.Println(fmt.Sprintf("Sum         %d", nrOfIntelligent+nrOfStupid+nrOfNone))

	return data
}
