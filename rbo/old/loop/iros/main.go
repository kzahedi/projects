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
	CreateResultsContainer(hands, ctrls, directory, &results)
	n := len(results)

	////////////////////////////////////////////////////////////
	// convert SOFA files to CSV
	// including preprocessing (conversation to wrist frame)
	////////////////////////////////////////////////////////////

	// ConvertSofaStates("hand.sofastates.txt", hands, ctrls, directory, true)
	// ConvertSofaStates("obstacle.sofastates.txt", hands, ctrls, directory, false)

	ConvertSofaStatesFull("hand.sofastates.txt", hands, ctrls, directory, true)
	ConvertSofaStatesFull("obstacle.sofastates.txt", hands, ctrls, directory, false)

	////////////////////////////////////////////////////////////
	// calculate difference behaviour (grasp - prescriptive)
	////////////////////////////////////////////////////////////

	// CalculateDifferenceBehaviour(hands, ctrls, rbohand2p, directory)

	////////////////////////////////////////////////////////////
	// calculate co-variance matrices
	////////////////////////////////////////////////////////////

	CalculateCovarianceMatrices(hands, ctrls, directory, 75)

	////////////////////////////////////////////////////////////
	// determine if successful or not
	////////////////////////////////////////////////////////////

	CalculateSuccess(hands, ctrls, directory, 50.0, &results)

	////////////////////////////////////////////////////////////
	// Calculating MC_W
	////////////////////////////////////////////////////////////

	CalculateMCW(hands, ctrls, directory, 100, 30, &results)

	////////////////////////////////////////////////////////////
	// Calculating Grasp Distance
	////////////////////////////////////////////////////////////

	CalculateGraspDistance(hands, ctrls, directory, 10, 500, &results)

	////////////////////////////////////////////////////////////
	// Convert object position to integer values
	////////////////////////////////////////////////////////////

	ExtractObjectPosition(&results)

	////////////////////////////////////////////////////////////
	// Convert object type to integer values
	////////////////////////////////////////////////////////////

	ExtractObjectType(&results)

	////////////////////////////////////////////////////////////
	// Calculate t-SNE
	////////////////////////////////////////////////////////////

	CalculateTSNE(rbohand2, controller0, directory, 10000, false, &results)

	// PrintResults(results)
	WriteResults("/Users/zahedi/Desktop/iros_results.csv", &results)

	CalculateTSNE(hands, ctrls, directory, 10000, false, &results)

	// PrintResults(results)
	WriteResults("/Users/zahedi/Desktop/iros_results.csv", &results)

}
