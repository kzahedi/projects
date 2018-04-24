package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	// directory := flag.String("d", "", "Directory")
	directory := flag.String("d", "/Users/zahedi/projects/TU.Berlin/experiments/run2017011101/", "Directory")
	flag.Parse()

	if *directory == "" {
		fmt.Println("Please provide a directory to analyse.")
		os.Exit(0)
	}

	////////////////////////////////////////////////////////////
	// define regexp patterns
	////////////////////////////////////////////////////////////
	rbohand2 := regexp.MustCompile(".*/rbohand2/.*")
	rbohand2p := regexp.MustCompile(".*/rbohand2-prescriptive/.*")
	rbohandkz1 := regexp.MustCompile(".*/rbohandkz1/.*")
	rbohandkz1p := regexp.MustCompile(".*/rbohandkz1-prescriptive/.*")
	rbohandkz2 := regexp.MustCompile(".*/rbohandkz2/.*")
	rbohandkz2p := regexp.MustCompile(".*/rbohandkz2-prescriptive/.*")

	controller0 := regexp.MustCompile(".*-controller0-.*")
	controller1 := regexp.MustCompile(".*-controller1-.*")
	controller2 := regexp.MustCompile(".*-controller2-.*")

	hands := []*regexp.Regexp{rbohand2, rbohand2p, rbohandkz1, rbohandkz1p, rbohandkz2, rbohandkz2p}
	ctrls := []*regexp.Regexp{controller0, controller1, controller2}

	////////////////////////////////////////////////////////////
	// Creating container for all results
	////////////////////////////////////////////////////////////

	results := make(Results)
	segmentResults := make(Results)
	CreateResultsContainer(hands, ctrls, directory, &results)
	CreateResultsContainer(hands, ctrls, directory, &segmentResults)

	////////////////////////////////////////////////////////////
	// convert SOFA files to CSV
	// including preprocessing (conversation to wrist frame)
	////////////////////////////////////////////////////////////

	ConvertSofaStates("hand.sofastates.txt", hands, ctrls, directory, true)
	ConvertSofaStates("obstacle.sofastates.txt", hands, ctrls, directory, false)
	ConvertSofaStatesSegment("hand.sofastates.txt", hands, ctrls, directory, true)

	////////////////////////////////////////////////////////////
	// calculate difference behaviour (grasp - prescriptive)
	////////////////////////////////////////////////////////////

	CalculateDifferenceBehaviour(hands, ctrls, rbohand2p, directory)
	CalculateDifferenceBehaviourSegment(hands, ctrls, rbohand2p, directory)

	////////////////////////////////////////////////////////////
	// calculate co-variance matrices
	////////////////////////////////////////////////////////////

	CalculateCovarianceMatrices(hands, ctrls, directory, 75)
	CalculateCovarianceMatricesSegments(hands, ctrls, directory, 75)

	////////////////////////////////////////////////////////////
	// determine if successful or not
	////////////////////////////////////////////////////////////

	CalculateSuccess(hands, ctrls, directory, 50.0, &results)
	CalculateSuccess(hands, ctrls, directory, 50.0, &segmentResults)

	////////////////////////////////////////////////////////////
	// Calculating MC_W
	////////////////////////////////////////////////////////////

	CalculateMCW(hands, ctrls, directory, 100, 30, &results)
	CalculateMCW(hands, ctrls, directory, 100, 30, &segmentResults)

	////////////////////////////////////////////////////////////
	// Calculating Grasp Distance
	////////////////////////////////////////////////////////////

	CalculateGraspDistance(hands, ctrls, directory, 10, 500, &results)
	CalculateGraspDistance(hands, ctrls, directory, 10, 500, &segmentResults)

	////////////////////////////////////////////////////////////
	// Convert object position to integer values
	////////////////////////////////////////////////////////////

	ExtractObjectPosition(&results)
	ExtractObjectPosition(&segmentResults)

	////////////////////////////////////////////////////////////
	// Convert object type to integer values
	////////////////////////////////////////////////////////////

	ExtractObjectType(&results)
	ExtractObjectType(&segmentResults)

	////////////////////////////////////////////////////////////
	// Calculate t-SNE
	////////////////////////////////////////////////////////////

	CalculateTSNE(rbohand2, controller0, directory, 10000, false, &results)
	CalculateTSNESegments(rbohand2, controller0, directory, 10000, false, &segmentResults)

	WriteResults("/Users/zahedi/Desktop/results.csv", &results)
	WriteResults("/Users/zahedi/Desktop/segment.full.results.csv", &segmentResults)

	PrintResults(results)
}
