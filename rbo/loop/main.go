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
	outputDirectory := flag.String("o", "/Users/zahedi/Desktop/", "Output Directory")
	iros := flag.Bool("iros", false, "Calculate IROS Results")
	// segment := flag.Bool("segment", false, "Calculate with segment tips transformed into local coordinate systems of the segment roots")
	// frameByFrame := flag.Bool("fbf", false, "Calculate with coordinate system transformed into the local predecessor coordinate system")
	percentage := flag.Float64("p", 0.15, "Cut-off percentage for intelligent and stupid")
	maxGraspDistance := flag.Float64("mgd", 250.0, "Cut-off for grasp distance")
	minLiftHeight := flag.Float64("mlh", 50.0, "Min lifting height for successful grasps.")
	stabilFactor := flag.Float64("s", 2.0, "How much bigger the mean values must be compared to the standard deviation.")
	trajectoryLength := flag.Int("t", 75, "The number of data points for covariance calculations.")
	tsneIterations := flag.Int("tsne", 10000, "Number of iterations for t-SNE")
	k := flag.Int("k", 10, "k-nearest neighbour after clustering")
	flag.Parse()

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

	////////////////////////////////////////////////////////////
	// IROS Results
	////////////////////////////////////////////////////////////

	// Creating container for all results

	if *iros == true {
		outputDir := fmt.Sprintf("%s/%s", *outputDirectory, "iros")
		fmt.Println("Checking ", outputDir)
		if _, err := os.Stat(outputDir); err != nil {
			os.MkdirAll(outputDir, 0755)
		}

		fmt.Println(">>> Calculating IROS Results")
		irosResults := make(Results)
		CreateResultsContainer(grasps, ctrls, directory, &irosResults)

		irosHandSofaStatesFile := "iros.hand.sofastates.csv"
		irosDiffHandSofaStatesFile := "iros.diffed.hand.sofastates.csv"
		irosCovarianceFile := "iros.covariance.csv"
		irosClusterInfoFile := "iros.cluster.info.txt"
		irosResultsFile := "iros.results.csv"
		irosIntelligentFile := "iros.intelligent.csv"
		irosStupidFile := "iros.stupid.csv"
		irosIntelligentHumanFile := "iros.intelligent.human.csv"
		irosStupidHumanFile := "iros.stupid.human.csv"
		irosAnalysisFile := "iros.analysis.txt"
		irosTSNEFile := "iros.tsne.results.txt"

		if utils.FileExists(fmt.Sprintf("%s/%s", outputDir, irosTSNEFile)) == true {
			irosResults = ReadResults(fmt.Sprintf("%s/%s", outputDir, irosTSNEFile))
		} else {
			// these four liens are needed for grasp success calculations
			ConvertSofaStatesPreprocessing(obstacleSofaStates, grasps, ctrls, directory, false)
			ConvertSofaStatesPreprocessing(obstacleSofaStates, prescritives, ctrls, directory, false)
			ConvertSofaStatesPreprocessing(handSofaStates, grasps, ctrls, directory, false)
			ConvertSofaStatesPreprocessing(handSofaStates, prescritives, ctrls, directory, false)

			// convert SOFA files to CSV
			// including preprocessing

			ConvertSofaStatesIROS(handSofaStates, irosHandSofaStatesFile, grasps, ctrls, directory)
			ConvertSofaStatesIROS(handSofaStates, irosHandSofaStatesFile, prescritives, ctrls, directory)

			// calculate difference behaviour (grasp - prescriptive)

			CalculateDifferenceBehaviour(irosHandSofaStatesFile, irosDiffHandSofaStatesFile, grasps, ctrls, directory)

			// calculate co-variance matrices

			CalculateCovarianceMatrices(irosDiffHandSofaStatesFile, irosCovarianceFile, grasps, ctrls, directory, *trajectoryLength, MODE_FULL)

			// determine if successful or not

			irosResults = CalculateSuccess(obstacleSofaStatesCsv, grasps, ctrls, directory, *minLiftHeight, irosResults) // checked

			// Calculating MC_W

			irosResults = CalculateMCW(grasps, ctrls, directory, 100, 30, irosResults) // checked

			// Calculating Grasp Distance

			irosResults = CalculateGraspDistance(grasps, ctrls, directory, 10, 500, irosResults) // checked

			// Convert object position to integer values

			irosResults = ExtractObjectPosition(irosResults) // checked

			// Convert object type to integer values

			irosResults = ExtractObjectType(irosResults) // checked

			// Calculate t-SNE

			irosResults = CalculateTSNE(irosCovarianceFile, grasps, ctrls, directory, *tsneIterations, false, irosResults, irosTSNEFile, outputDir)
		}

		// Calculate Clusters

		fmt.Println("hier 0")
		irosResults = CalculateInterestingClusters(irosResults, *maxGraspDistance, *percentage, *k, irosClusterInfoFile, outputDir) // checked
		fmt.Println("hier 1")

		WriteResults(irosResultsFile, irosResults, outputDir)
		fmt.Println("hier 2")

		intelligent := AnalyseIntelligent(irosResults, *directory, irosCovarianceFile, irosIntelligentFile, outputDir)
		fmt.Println("hier 3")
		stupid := AnalyseStupid(irosResults, *directory, irosCovarianceFile, irosStupidFile, outputDir)
		fmt.Println("hier 4")

		ConvertIROSMatrixResults(outputDir, irosIntelligentFile, irosIntelligentHumanFile)
		fmt.Println("hier 5")
		ConvertIROSMatrixResults(outputDir, irosStupidFile, irosStupidHumanFile)
		fmt.Println("hier 6")

		irosAnalysis := AnalyseData(intelligent, stupid, *stabilFactor)
		fmt.Println("hier 7")

		ConvertIROSAnalysisResults(outputDir, irosAnalysisFile, irosAnalysis)
		fmt.Println("hier 8")

	}
}

//
// 	////////////////////////////////////////////////////////////
// 	// Segment Results
// 	////////////////////////////////////////////////////////////
//
// 	// Creating container for all results
//
// 	if *segment == true {
// 		fmt.Println(">>> Calculating Segment Results")
// 		segmentResults := make(Results)
// 		CreateResultsContainer(grasps, ctrls, directory, &segmentResults)
//
// 		segmentHandSofaStates := "segment.hand.sofastates.csv"
// 		segmentDiffHandSofaStates := "segment.diffed.hand.sofastates.csv"
// 		segmentCovariance := "segment.covariance.csv"
//
// 		// convert SOFA files to CSV
// 		// including preprocessing
//
// 		ConvertSofaStatesSegment(handSofaStates, segmentHandSofaStates, grasps, ctrls, directory)       // checked
// 		ConvertSofaStatesSegment(handSofaStates, segmentHandSofaStates, prescritives, ctrls, directory) // checked
//
// 		// calculate difference behaviour (grasp - prescriptive)
//
// 		CalculateDifferenceBehaviour(segmentHandSofaStates, segmentDiffHandSofaStates, grasps, ctrls, directory) // checked
//
// 		// calculate co-variance matrices
//
// 		CalculateCovarianceMatrices(segmentDiffHandSofaStates, segmentCovariance, grasps, ctrls, directory, *trajectoryLength, MODE_SEGMENT) // checked
//
// 		// determine if successful or not
//
// 		segmentResults = CalculateSuccess(obstacleSofaStatesCsv, grasps, ctrls, directory, *minLiftHeight, segmentResults) // checked
//
// 		// Calculating MC_W
//
// 		segmentResults = CalculateMCW(grasps, ctrls, directory, 100, 30, segmentResults) // checked
//
// 		// Calculating Grasp Distance
//
// 		segmentResults = CalculateGraspDistance(grasps, ctrls, directory, 10, 500, segmentResults) // checked
//
// 		// Convert object position to integer values
//
// 		segmentResults = ExtractObjectPosition(segmentResults) // checked
//
// 		// Convert object type to integer values
//
// 		segmentResults = ExtractObjectType(segmentResults) // checked
//
// 		// Calculate t-SNE
//
// 		// segmentResults = CalculateTSNE("segment.covariance.csv", []*regexp.Regexp{rbohand2}, ctrls, directory, *tsneIterations, false, segmentResults) // checked
// 		segmentResults = CalculateTSNE("segment.covariance.csv", grasps, ctrls, directory, *tsneIterations, false, segmentResults) // checked
//
// 		// Calculate Clusters
//
// 		segmentResults = CalculateInterestingClusters(segmentResults, *maxGraspDistance, *percentage, *k, "/Users/zahedi/Desktop/segment.cluster.info.txt") // checked
//
// 		WriteResults("/Users/zahedi/Desktop/segment.results.csv", segmentResults)
//
// 		// WriteResults("/Users/zahedi/Desktop/segment.results.csv", &segmentResults)
//
// 		AnalyseIntelligent(segmentResults, *directory, segmentCovariance, "/Users/zahedi/Desktop/segment.intelligent.csv") // checked
// 		AnalyseStupid(segmentResults, *directory, segmentCovariance, "/Users/zahedi/Desktop/segment.stupid.csv")           // checked
//
// 		ConvertSegmentMatrixResults("/Users/zahedi/Desktop/segment.intelligent.csv")
// 		ConvertSegmentMatrixResults("/Users/zahedi/Desktop/segment.stupid.csv")
//
// 	}
//
// 	////////////////////////////////////////////////////////////
// 	// Frame by Frame Results
// 	////////////////////////////////////////////////////////////
//
// 	if *frameByFrame == true {
// 		fmt.Println(">>> Calculating Frame By Frame Results")
// 		frameByFrameResults := make(Results)
// 		CreateResultsContainer(grasps, ctrls, directory, &frameByFrameResults)
//
// 		frameByFrameHandSofaStates := "frame.by.frame.hand.sofastates.csv"
// 		frameByFrameDiffHandSofaStates := "frame.by.frame.diffed.hand.sofastates.csv"
// 		frameByFrameCovariance := "frame.by.frame.covariance.csv"
//
// 		// convert SOFA files to CSV
// 		// including preprocessing
//
// 		ConvertSofaStatesFrameByFrame(handSofaStates, frameByFrameHandSofaStates, grasps, ctrls, directory)
// 		ConvertSofaStatesFrameByFrame(handSofaStates, frameByFrameHandSofaStates, prescritives, ctrls, directory)
//
// 		// calculate difference behaviour (grasp - prescriptive)
//
// 		CalculateDifferenceBehaviour(frameByFrameHandSofaStates, frameByFrameDiffHandSofaStates, grasps, ctrls, directory)
//
// 		// calculate co-variance matrices
//
// 		CalculateCovarianceMatrices(frameByFrameDiffHandSofaStates, frameByFrameCovariance, grasps, ctrls, directory, *trajectoryLength, MODE_FRAME_BY_FRAME)
//
// 		// determine if successful or not
//
// 		frameByFrameResults = CalculateSuccess(obstacleSofaStatesCsv, grasps, ctrls, directory, *minLiftHeight, frameByFrameResults) // checked
//
// 		// Calculating MC_W
//
// 		frameByFrameResults = CalculateMCW(grasps, ctrls, directory, 100, 30, frameByFrameResults) // checked
//
// 		// Calculating Grasp Distance
//
// 		frameByFrameResults = CalculateGraspDistance(grasps, ctrls, directory, 10, 500, frameByFrameResults) // checked
//
// 		// Convert object position to integer values
//
// 		frameByFrameResults = ExtractObjectPosition(frameByFrameResults) // checked
//
// 		// Convert object type to integer values
//
// 		frameByFrameResults = ExtractObjectType(frameByFrameResults) // checked
//
// 		// Calculate t-SNE
//
// 		// frameByFrameResults = CalculateTSNE("frame.by.frame.covariance.csv", []*regexp.Regexp{rbohand2}, ctrls, directory, *tsneIterations, false, frameByFrameResults)
// 		frameByFrameResults = CalculateTSNE("frame.by.frame.covariance.csv", grasps, ctrls, directory, *tsneIterations, false, frameByFrameResults)
//
// 		// Calculate Clusters
//
// 		frameByFrameResults = CalculateInterestingClusters(frameByFrameResults, *maxGraspDistance, *percentage, *k, "/Users/zahedi/Desktop/frameByFrame.cluster.info.txt") // checked
//
// 		WriteResults("/Users/zahedi/Desktop/frame.by.frame.results.csv", frameByFrameResults)
//
// 		AnalyseIntelligent(frameByFrameResults, *directory, frameByFrameCovariance, "/Users/zahedi/Desktop/frameByFrame.intelligent.csv")
// 		AnalyseStupid(frameByFrameResults, *directory, frameByFrameCovariance, "/Users/zahedi/Desktop/frameByFrame.stupid.csv")
//
// 		ConvertFrameByFrameMatrixResults("/Users/zahedi/Desktop/frameByFrame.intelligent.csv")
// 		ConvertFrameByFrameMatrixResults("/Users/zahedi/Desktop/frameByFrame.stupid.csv")
//
// 	}
//
