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
	// rbohand2p := regexp.MustCompile(".*/rbohand2-prescriptive/.*")

	controller0 := regexp.MustCompile(".*-controller0-.*")

	// rbohandkz1 := regexp.MustCompile(".*/rbohandkz1/.*")
	// rbohandkz1p := regexp.MustCompile(".*/rbohandkz1-prescriptive/.*")
	// rbohandkz2 := regexp.MustCompile(".*/rbohandkz2/.*")
	// rbohandkz2p := regexp.MustCompile(".*/rbohandkz2-prescriptive/.*")

	////////////////////////////////////////////////////////////
	// convert SOFA files to CSV
	// including preprocessing (conversation to wrist frame)
	////////////////////////////////////////////////////////////

	// ConvertSofaStates("hand.sofastates.txt", rbohand2, controller0, directory, true)
	// ConvertSofaStates("obstacle.sofastates.txt", rbohand2, controller0, directory, false)

	// ConvertSofaStates("hand.sofastates.txt", rbohand2p, controller0, directory, true)
	// ConvertSofaStates("obstacle.sofastates.txt", rbohand2p, controller0, directory, false)

	////////////////////////////////////////////////////////////
	// calculate difference behaviour (grasp - prescriptive)
	////////////////////////////////////////////////////////////

	// CalculateDifferenceBehaviour(rbohand2, rbohand2p, controller0, directory)

	////////////////////////////////////////////////////////////
	// calculate co-variance matrices
	////////////////////////////////////////////////////////////

	// CalculateCovarianceMatrices(rbohand2, controller0, directory, 75)

	////////////////////////////////////////////////////////////
	// t-SNE
	////////////////////////////////////////////////////////////

	CalculateTSNE(rbohand2, controller0, directory)
}
