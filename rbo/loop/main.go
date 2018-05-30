package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kzahedi/utils"
)

type fnConvertSofaState func(string, string, []*regexp.Regexp, []*regexp.Regexp, string)
type fnConvertMatrixResults func(string, string, string)
type fnConvertAnalysisResults func(dir, output string, analysis []Analysis)

func main() {
	// directory := flag.String("d", "", "Directory")
	directory := flag.String("d", "/Users/zahedi/projects/TU.Berlin/experiments/run2017011101/", "Directory")
	outputDirectory := flag.String("o", "/Users/zahedi/Desktop/", "Output Directory")
	iros := flag.Bool("iros", false, "Calculate IROS Results")
	segment := flag.Bool("segment", false, "Calculate with segment tips transformed into local coordinate systems of the segment roots")
	frameByFrame := flag.Bool("fbf", false, "Calculate with coordinate system transformed into the local predecessor coordinate system")
	percentage := flag.Float64("p", 0.15, "Cut-off percentage for intelligent and stupid")
	maxGraspDistance := flag.Float64("mgd", 250.0, "Cut-off for grasp distance")
	minLiftHeight := flag.Float64("mlh", 50.0, "Min lifting height for successful grasps.")
	stabilFactor := flag.Float64("s", 2.0, "How much bigger the mean values must be compared to the standard deviation.")
	trajectoryLength := flag.Int("t", 75, "The number of data points for covariance calculations.")
	tsneIterations := flag.Int("tsne", 10000, "Number of iterations for t-SNE")
	k := flag.Int("k", 10, "k-nearest neighbour after clustering")
	flag.Parse()

	////////////////////////////////////////////////////////////
	// IROS Results
	////////////////////////////////////////////////////////////

	// Creating container for all results

	// default
	prefix := ""

	if *iros == true {
		prefix = "iros"
		outputDir := fmt.Sprintf("%s/%s", *outputDirectory, prefix)
		convertSofaStates := ConvertSofaStatesIROS
		convertMatrixResults := ConvertIROSMatrixResults
		convertAnalysisResults := ConvertIROSAnalysisResults
		doIt(*directory, outputDir, prefix, convertSofaStates, convertMatrixResults, convertAnalysisResults, *percentage, *maxGraspDistance, *minLiftHeight, *stabilFactor, *trajectoryLength, *tsneIterations, *k)
	}
	if *segment == true {
		prefix = "segment"
		outputDir := fmt.Sprintf("%s/%s", *outputDirectory, prefix)
		convertSofaStates := ConvertSofaStatesSegment
		convertMatrixResults := ConvertSegmentMatrixResults
		convertAnalysisResults := ConvertIROSAnalysisResults
		// convertAnalysisResults := ConvertSegmentAnalysisResults
		doIt(*directory, outputDir, prefix, convertSofaStates, convertMatrixResults, convertAnalysisResults, *percentage, *maxGraspDistance, *minLiftHeight, *stabilFactor, *trajectoryLength, *tsneIterations, *k)
	}
	if *frameByFrame == true {
		prefix = "frame.by.frame"
		outputDir := fmt.Sprintf("%s/%s", *outputDirectory, prefix)
		convertSofaStates := ConvertSofaStatesFrameByFrame
		convertMatrixResults := ConvertFrameByFrameMatrixResults
		convertAnalysisResults := ConvertIROSAnalysisResults
		// convertAnalysisResults := ConvertFrameByFrameAnalysisResults
		doIt(*directory, outputDir, prefix, convertSofaStates, convertMatrixResults, convertAnalysisResults, *percentage, *maxGraspDistance, *minLiftHeight, *stabilFactor, *trajectoryLength, *tsneIterations, *k)
	}

	if prefix == "" {
		fmt.Println("Please choose a method")
		os.Exit(0)
	}

}

func doIt(directory, outputDir, prefix string,
	convertSofaStates fnConvertSofaState,
	convertMatrixResults fnConvertMatrixResults,
	convertAnalysisResults fnConvertAnalysisResults,
	percentage, maxGraspDistance, minLiftHeight,
	stabilFactor float64, trajectoryLength, tsneIterations, k int) {

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

	fmt.Println("Checking ", outputDir)
	if _, err := os.Stat(outputDir); err != nil {
		os.MkdirAll(outputDir, 0755)
	}

	fmt.Println(fmt.Sprintf(">>> Calculating %s Results", strings.ToUpper(prefix)))
	results := make(Results)
	CreateResultsContainer(grasps, ctrls, directory, &results)

	handSofaStatesFile := fmt.Sprintf("%s.hand.sofastates.csv", prefix)
	diffHandSofaStatesFile := fmt.Sprintf("%s.diffed.hand.sofastates.csv", prefix)
	covarianceFile := fmt.Sprintf("%s.covariance.csv", prefix)
	clusterInfoFile := fmt.Sprintf("%s.cluster.info.txt", prefix)
	resultsFile := fmt.Sprintf("%s.results.csv", prefix)
	intelligentFile := fmt.Sprintf("%s.intelligent.csv", prefix)
	stupidFile := fmt.Sprintf("%s.stupid.csv", prefix)
	intelligentHumanFile := fmt.Sprintf("%s.intelligent.human.csv", prefix)
	stupidHumanFile := fmt.Sprintf("%s.stupid.human.csv", prefix)
	analysisFile := fmt.Sprintf("%s.analysis.txt", prefix)
	tsneFile := fmt.Sprintf("%s.tsne.results.txt", prefix)

	if utils.FileExists(fmt.Sprintf("%s/%s", outputDir, tsneFile)) == true {
		results = ReadResults(fmt.Sprintf("%s/%s", outputDir, tsneFile))
	} else {
		// these four liens are needed for grasp success calculations
		ConvertSofaStatesPreprocessing(obstacleSofaStates, grasps, ctrls, directory, false)
		ConvertSofaStatesPreprocessing(obstacleSofaStates, prescritives, ctrls, directory, false)
		ConvertSofaStatesPreprocessing(handSofaStates, grasps, ctrls, directory, false)
		ConvertSofaStatesPreprocessing(handSofaStates, prescritives, ctrls, directory, false)

		// convert SOFA files to CSV
		// including preprocessing

		convertSofaStates(handSofaStates, handSofaStatesFile, grasps, ctrls, directory)
		convertSofaStates(handSofaStates, handSofaStatesFile, prescritives, ctrls, directory)

		// calculate difference behaviour (grasp - prescriptive)

		CalculateDifferenceBehaviour(handSofaStatesFile, diffHandSofaStatesFile, grasps, ctrls, directory)

		// calculate co-variance matrices

		CalculateCovarianceMatrices(diffHandSofaStatesFile, covarianceFile, grasps, ctrls, directory, trajectoryLength, MODE_FULL)

		// determine if successful or not

		results = CalculateSuccess(obstacleSofaStatesCsv, grasps, ctrls, directory, minLiftHeight, results) // checked

		// Calculating MC_W

		results = CalculateMCW(grasps, ctrls, directory, 100, 30, results) // checked

		// Calculating Grasp Distance

		results = CalculateGraspDistance(grasps, ctrls, directory, 10, 500, results) // checked

		// Convert object position to integer values

		results = ExtractObjectPosition(results) // checked

		// Convert object type to integer values

		results = ExtractObjectType(results) // checked

		// Calculate t-SNE

		results = CalculateTSNE(covarianceFile, grasps, ctrls, directory, tsneIterations, false, results, tsneFile, outputDir)
	}

	// Calculate Clusters

	results = CalculateInterestingClusters(results, maxGraspDistance, percentage, k, clusterInfoFile, outputDir) // checked

	WriteResults(resultsFile, results, outputDir)

	intelligent := AnalyseIntelligent(results, directory, covarianceFile, intelligentFile, outputDir)
	stupid := AnalyseStupid(results, directory, covarianceFile, stupidFile, outputDir)

	convertMatrixResults(outputDir, intelligentFile, intelligentHumanFile)
	convertMatrixResults(outputDir, stupidFile, stupidHumanFile)

	analysis := AnalyseData(intelligent, stupid, stabilFactor)

	convertAnalysisResults(outputDir, analysisFile, analysis)
}
