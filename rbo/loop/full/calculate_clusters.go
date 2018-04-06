package main

func CalculateInteretingClusters(filename string, graspDistanceCutoff, percentage float64) Results {

	data := ReadResults(filename)

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
	}

	mcwDiff := maxMCW * percentage
	gdDiff := maxGraspDistance * percentage

	stupidMCW := minMCW + mcwDiff
	intelligentMCW := maxMCW - mcwDiff

	stupidGraspDistance := maxGraspDistance - gdDiff
	intelligentGraspDistance := minGraspDistance + gdDiff

	for key, value := range data {
		mcw := value.MC_W
		gd := value.GraspDistance

		value.Stupid = false
		value.Intelligent = false

		if mcw < stupidMCW && gd > stupidGraspDistance {
			value.Stupid = true
			value.Intelligent = false
		} else if mcw > intelligentMCW && gd < intelligentGraspDistance {
			value.Stupid = false
			value.Intelligent = true
		}
		data[key] = value
	}

	return data
}
