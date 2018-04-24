package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/kzahedi/utils"
)

func main() {
	// directory := flag.String("d", "", "Directory")
	directory := flag.String("d", "/Users/zahedi/projects/TU.Berlin/experiments/run2017011101/", "Directory")
	iros := flag.Bool("iros", false, "Calculate IROS Results")
	segment := flag.Bool("segment", false, "Calculate with segment tips transformed into local coordinate systems of the segment roots")
	frameByFrame := flag.Bool("fbf", false, "Calculate with coordinate system transformed into the local predecessor coordinate system")
	percentage := flag.Float64("p", 0.15, "Cut-off percentage for intelligent and stupid")
	maxGraspDistance := flag.Float64("mgd", 250.0, "Cut-off for grasp distance")
	minLiftHeight := flag.Float64("mlh", 50.0, "Min lifting height for successful grasps.")
	trajectoryLength := flag.Int("t", 75, "The number of data points for covariance calculations.")
	tsneIterations := flag.Int("tsne", 10000, "Number of iterations for t-SNE")
	k := flag.Int("k", 10, "k-nearest neighbour after clustering")
	test := flag.String("test", "", "Test")
	flag.Parse()

	if *test != "" {
		header := []string{"a", "b", "c", "d"}
		header2 := []string{"# a", "b", "c", "d"}
		data := [][]string{{"1", "2", "3", "4"}, {"11", "12", "13", "14"}, {"21", "22", "23", "24"}, {"31", "32", "33", "34"}, {"41", "42", "43", "44"}}
		utils.WriteCsv("/Users/zahedi/Desktop/test.1.csv", data, header)
		utils.WriteCsv("/Users/zahedi/Desktop/test.2.csv", data, nil)
		utils.WriteCsv("/Users/zahedi/Desktop/test.3.csv", data, header2)

		d1, h1 := utils.ReadCsv("/Users/zahedi/Desktop/test.1.csv")
		d2, h2 := utils.ReadCsv("/Users/zahedi/Desktop/test.2.csv")
		d3, h3 := utils.ReadCsv("/Users/zahedi/Desktop/test.3.csv")

		fmt.Println(d1)
		fmt.Println(h1)

		fmt.Println(d2)
		fmt.Println(h2)

		fmt.Println(d3)
		fmt.Println(h3)

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

	grasps := []*regexp.Regexp{rbohand2, rbohandkz1, rbohandkz2}
	prescritives := []*regexp.Regexp{rbohand2p, rbohandkz1p, rbohandkz2p}
	ctrls := []*regexp.Regexp{controller0, controller1, controller2}

	////////////////////////////////////////////////////////////
	// Preprocessing
	////////////////////////////////////////////////////////////

	handSofaStates := "hand.sofastates.txt"
	obstacleSofaStates := "obstacle.sofastates.txt"
	obstacleSofaStatesCsv := "obstacle.sofastates.csv"

	// these four liens are needed for grasp success calculations
	ConvertSofaStatesPreprocessing(obstacleSofaStates, grasps, ctrls, directory, false)
	ConvertSofaStatesPreprocessing(obstacleSofaStates, prescritives, ctrls, directory, false)
	ConvertSofaStatesPreprocessing(handSofaStates, grasps, ctrls, directory, false)
	ConvertSofaStatesPreprocessing(handSofaStates, prescritives, ctrls, directory, false)

	////////////////////////////////////////////////////////////
	// IROS Results
	////////////////////////////////////////////////////////////

	// Creating container for all results

	if *iros == true {
		fmt.Println(">>> Calculating IROS Results")
		irosResults := make(Results)
		CreateResultsContainer(grasps, ctrls, directory, &irosResults)

		irosHandSofaStates := "iros.hand.sofastates.csv"
		irosDiffHandSofaStates := "iros.diffed.hand.sofastates.csv"
		irosCovariance := "iros.covariance.csv"

		// convert SOFA files to CSV
		// including preprocessing

		ConvertSofaStatesIROS(handSofaStates, irosHandSofaStates, grasps, ctrls, directory)
		ConvertSofaStatesIROS(handSofaStates, irosHandSofaStates, prescritives, ctrls, directory)

		// calculate difference behaviour (grasp - prescriptive)

		CalculateDifferenceBehaviour(irosHandSofaStates, irosDiffHandSofaStates, grasps, ctrls, directory)

		// calculate co-variance matrices

		CalculateCovarianceMatrices(irosDiffHandSofaStates, irosCovariance, grasps, ctrls, directory, *trajectoryLength, MODE_FULL)

		// determine if successful or not

		irosResults = CalculateSuccess(obstacleSofaStatesCsv, grasps, ctrls, directory, *minLiftHeight, irosResults) // checked

		// Calculating MC_W

		irosResults = CalculateMCW(grasps, ctrls, directory, 100, 30, irosResults)

		// Calculating Grasp Distance

		irosResults = CalculateGraspDistance(grasps, ctrls, directory, 10, 500, irosResults)

		// Convert object position to integer values

		irosResults = ExtractObjectPosition(irosResults)

		// Convert object type to integer values

		irosResults = ExtractObjectType(irosResults)

		// Calculate t-SNE

		irosResults = CalculateTSNE("iros.covariance.csv", rbohand2, controller0, directory, *tsneIterations, false, irosResults)

		// Calculate Clusters

		irosResults = CalculateInteretingClusters(irosResults, *maxGraspDistance, *percentage, *k)

		WriteResults("/Users/zahedi/Desktop/iros.results.csv", irosResults)

		AnalyseIntelligent(irosResults, *directory, irosCovariance, "/Users/zahedi/Desktop/iros.intelligent.csv")
		AnalyseStupid(irosResults, *directory, irosCovariance, "/Users/zahedi/Desktop/iros.stupid.csv")

		ConvertIROSMatrixResults("/Users/zahedi/Desktop/iros.intelligent.csv")
		ConvertIROSMatrixResults("/Users/zahedi/Desktop/iros.stupid.csv")

	}

	////////////////////////////////////////////////////////////
	// Segment Results
	////////////////////////////////////////////////////////////

	// Creating container for all results

	if *segment == true {
		fmt.Println(">>> Calculating Segment Results")
		segmentResults := make(Results)
		CreateResultsContainer(grasps, ctrls, directory, &segmentResults)

		segmentHandSofaStates := "segment.hand.sofastates.csv"
		segmentDiffHandSofaStates := "segment.diffed.hand.sofastates.csv"
		segmentCovariance := "segment.covariance.csv"

		// convert SOFA files to CSV
		// including preprocessing

		ConvertSofaStatesSegment(handSofaStates, segmentHandSofaStates, grasps, ctrls, directory)       // checked
		ConvertSofaStatesSegment(handSofaStates, segmentHandSofaStates, prescritives, ctrls, directory) // checked

		// calculate difference behaviour (grasp - prescriptive)

		CalculateDifferenceBehaviour(segmentHandSofaStates, segmentDiffHandSofaStates, grasps, ctrls, directory) // checked

		// calculate co-variance matrices

		CalculateCovarianceMatrices(segmentDiffHandSofaStates, segmentCovariance, grasps, ctrls, directory, *trajectoryLength, MODE_SEGMENT) // checked

		// determine if successful or not

		segmentResults = CalculateSuccess(obstacleSofaStatesCsv, grasps, ctrls, directory, *minLiftHeight, segmentResults) // checked

		// Calculating MC_W

		segmentResults = CalculateMCW(grasps, ctrls, directory, 100, 30, segmentResults)

		// Calculating Grasp Distance

		segmentResults = CalculateGraspDistance(grasps, ctrls, directory, 10, 500, segmentResults)

		// Convert object position to integer values

		segmentResults = ExtractObjectPosition(segmentResults)

		// Convert object type to integer values

		segmentResults = ExtractObjectType(segmentResults)

		// Calculate t-SNE

		segmentResults = CalculateTSNE("segment.covariance.csv", rbohand2, controller0, directory, *tsneIterations, false, segmentResults)

		// Calculate Clusters

		segmentResults = CalculateInteretingClusters(segmentResults, *maxGraspDistance, *percentage, *k)

		WriteResults("/Users/zahedi/Desktop/segment.results.csv", segmentResults)

		// WriteResults("/Users/zahedi/Desktop/segment.results.csv", &segmentResults)

		AnalyseIntelligent(segmentResults, *directory, segmentCovariance, "/Users/zahedi/Desktop/segment.intelligent.csv")
		AnalyseStupid(segmentResults, *directory, segmentCovariance, "/Users/zahedi/Desktop/segment.stupid.csv")

		ConvertSegmentMatrixResults("/Users/zahedi/Desktop/segment.intelligent.csv")
		ConvertSegmentMatrixResults("/Users/zahedi/Desktop/segment.stupid.csv")

	}

	////////////////////////////////////////////////////////////
	// Frame by Frame Results
	////////////////////////////////////////////////////////////

	if *frameByFrame == true {
		fmt.Println(">>> Calculating Frame By Frame Results")
		frameByFrameResults := make(Results)
		CreateResultsContainer(grasps, ctrls, directory, &frameByFrameResults)

		frameByFrameHandSofaStates := "frame.by.frame.hand.sofastates.csv"
		frameByFrameDiffHandSofaStates := "frame.by.frame.diffed.hand.sofastates.csv"
		frameByFrameCovariance := "frame.by.frame.covariance.csv"

		// convert SOFA files to CSV
		// including preprocessing

		ConvertSofaStatesFrameByFrame(handSofaStates, frameByFrameHandSofaStates, grasps, ctrls, directory)
		ConvertSofaStatesFrameByFrame(handSofaStates, frameByFrameHandSofaStates, prescritives, ctrls, directory)

		// calculate difference behaviour (grasp - prescriptive)

		CalculateDifferenceBehaviour(frameByFrameHandSofaStates, frameByFrameDiffHandSofaStates, grasps, ctrls, directory)

		// calculate co-variance matrices

		CalculateCovarianceMatrices(frameByFrameDiffHandSofaStates, frameByFrameCovariance, grasps, ctrls, directory, *trajectoryLength, MODE_FRAME_BY_FRAME)

		// determine if successful or not

		frameByFrameResults = CalculateSuccess(obstacleSofaStatesCsv, grasps, ctrls, directory, *minLiftHeight, frameByFrameResults) // checked

		// Calculating MC_W

		frameByFrameResults = CalculateMCW(grasps, ctrls, directory, 100, 30, frameByFrameResults)

		// Calculating Grasp Distance

		frameByFrameResults = CalculateGraspDistance(grasps, ctrls, directory, 10, 500, frameByFrameResults)

		// Convert object position to integer values

		frameByFrameResults = ExtractObjectPosition(frameByFrameResults)

		// Convert object type to integer values

		frameByFrameResults = ExtractObjectType(frameByFrameResults)

		// Calculate t-SNE

		frameByFrameResults = CalculateTSNE("frame.by.frame.covariance.csv", rbohand2, controller0, directory, *tsneIterations, false, frameByFrameResults)

		// Calculate Clusters

		frameByFrameResults = CalculateInteretingClusters(frameByFrameResults, *maxGraspDistance, *percentage, *k)

		WriteResults("/Users/zahedi/Desktop/frame.by.frame.results.csv", frameByFrameResults)

		AnalyseIntelligent(frameByFrameResults, *directory, frameByFrameCovariance, "/Users/zahedi/Desktop/frameByFrame.intelligent.csv")
		AnalyseStupid(frameByFrameResults, *directory, frameByFrameCovariance, "/Users/zahedi/Desktop/frameByFrame.stupid.csv")
	}
}
