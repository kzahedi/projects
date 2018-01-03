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

	hss := "hand.sofastates.txt"     //
	oss := "obstacle.sofastates.txt" //
	rbohand2 := regexp.MustCompile(".*/rbohand2/.*")

	hssFiles := ListAllFilesRecursivelyByFilename(*directory, hss)
	rbohand2HssFiles := Select(hssFiles, *rbohand2)
	for _, s := range rbohand2HssFiles[1:2] {
		data := ReadSofaSates(*directory, s) // returns 2d-array of pose
		data = CalculateGlobalVelocities(data)
		// for _, t := range data.Trajectories {
		// for i, p := range t.Frame {
		// fmt.Println("Frame", i, p)
		// }
		// fmt.Println("*******************")
		// }
	}

	ossFiles := ListAllFilesRecursivelyByFilename(*directory, oss)
	rbohand2OssFiles := Select(ossFiles, *rbohand2)
	for _, s := range rbohand2OssFiles[1:2] {
		ReadSofaSates(*directory, s)
	}

}
