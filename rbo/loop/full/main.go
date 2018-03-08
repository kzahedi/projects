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
	iros := flag.Bool("iros", false, "Calculate IROS Results")
	segment := flag.Bool("segment", false, "Calculate with segment tips transformed into local coordinate systems of the segment roots")
	full := flag.Bool("full", false, "Calculate with coordinate system transformed into the local predecessor coordinate system")
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
	// Preprocessing
	////////////////////////////////////////////////////////////

	// handSofaStates := "hand.sofastates.txt"
	// obstacleSofaStates := "obstacle.sofastates.txt"
	// ConvertSofaStates(obstacleSofaStates, hands, ctrls, directory, false)
	// ConvertSofaStates(handSofaStates, hands, ctrls, directory, false)

	////////////////////////////////////////////////////////////
	// IROS Results
	////////////////////////////////////////////////////////////

	// Creating container for all results

	if *iros == true {
		fmt.Println("Calculating IROS results")
		results := make(Results)
		CreateResultsContainer(hands, ctrls, directory, &results)

		// irosHandSofaStates := "iros.hand.sofastates.csv"
		// irosDiffHandSofaStates := "iros.diffed.hand.sofastates.csv"
		// irosCovariance := "iros.covariance.csv"

		// convert SOFA files to CSV
		// including preprocessing

		// ConvertSofaStatesIROS(handSofaStates, irosHandSofaStates, hands, ctrls, directory, true)

		// calculate difference behaviour (grasp - prescriptive)

		// CalculateDifferenceBehaviourIROS(irosHandSofaStates, irosDiffHandSofaStates, hands, ctrls, rbohand2p, directory)

		// calculate co-variance matrices

		// CalculateCovarianceMatrices(irosDiffHandSofaStates, irosCovariance, hands, ctrls, directory, 75)

		// determine if successful or not

		// CalculateSuccess(obstacleSofaStates, hands, ctrls, directory, 50.0, &results)

		// Calculating MC_W

		CalculateMCW(hands, ctrls, directory, 100, 30, &results)

		// Calculating Grasp Distance

		CalculateGraspDistance(hands, ctrls, directory, 10, 500, &results)

		// Convert object position to integer values

		ExtractObjectPosition(&results)

		// Convert object type to integer values

		ExtractObjectType(&results)

		// Calculate t-SNE

		CalculateTSNE("iros.covariance.csv", rbohand2, controller0, directory, 10000, false, &results)

		WriteResults("/Users/zahedi/Desktop/iros.results.csv", &results)
	}

	////////////////////////////////////////////////////////////
	// Segment Results
	////////////////////////////////////////////////////////////

	// Creating container for all results

	if *segment == true {
		fmt.Println("Calculating SEGMENT results")
		// segmentResults := make(Results)
		// CreateResultsContainer(hands, ctrls, directory, &segmentResults)

		// convert SOFA files to CSV
		// including preprocessing

		// ConvertSofaStatesSegment("hand.sofastates.txt", hands, ctrls, directory, true)

		// calculate difference behaviour (grasp - prescriptive)

		// CalculateDifferenceBehaviourSegment(hands, ctrls, rbohand2p, directory)

		// calculate co-variance matrices

		// CalculateCovarianceMatrices("difference.segment.hand.sofastates.csv", hands, ctrls, directory, 75)

		// determine if successful or not

		// CalculateSuccess(hands, ctrls, directory, 50.0, &segmentResults)

		// Calculating MC_W

		// CalculateMCW(hands, ctrls, directory, 100, 30, &segmentResults)

		// Calculating Grasp Distance

		// CalculateGraspDistance(hands, ctrls, directory, 10, 500, &segmentResults)

		// Convert object position to integer values

		// ExtractObjectPosition(&segmentResults)

		// Convert object type to integer values

		// ExtractObjectType(&segmentResults)

		// Calculate t-SNE

		// CalculateTSNESegments(rbohand2, controller0, directory, 10000, false, &segmentResults)

		// WriteResults("/Users/zahedi/Desktop/segment.results.csv", &segmentResults)
	}

	////////////////////////////////////////////////////////////
	// Frame by Frame Results
	////////////////////////////////////////////////////////////

	// Creating container for all results

	if *full == true {
		fmt.Println("Calculating Frame-by-Frame results")
		// frameByFrameResults := make(Results)
		// CreateResultsContainer(hands, ctrls, directory, &frameByFrameResults)

		// convert SOFA files to CSV
		// including preprocessing

		// ConvertSofaStatesFrameByFrame("hand.sofastates.txt", hands, ctrls, directory, true)

		// calculate difference behaviour (grasp - prescriptive)

		// CalculateDifferenceBehaviourFrameByFrame(hands, ctrls, rbohand2p, directory)

		// calculate co-variance matrices

		// CalculateCovarianceMatrices("difference.frame.by.frame.hand.sofastates.csv", hands, ctrls, directory, 75)

		// determine if successful or not

		// CalculateSuccess(hands, ctrls, directory, 50.0, &frameByFrameResults)

		// Calculating MC_W

		// CalculateMCW(hands, ctrls, directory, 100, 30, &frameByFrameResults)

		// Calculating Grasp Distance

		// CalculateGraspDistance(hands, ctrls, directory, 10, 500, &frameByFrameResults)

		// Convert object position to integer values

		// ExtractObjectPosition(&frameByFrameResults)

		// Convert object type to integer values

		// ExtractObjectType(&frameByFrameResults)

		// Calculate t-SNE

		// CalculateTSNEFrameByFrame(rbohand2, controller0, directory, 10000, false, &frameByFrameResults)

		// WriteResults("/Users/zahedi/Desktop/frame.by.frame.results.csv", &frameByFrameResults)
	}

}
