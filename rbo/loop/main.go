package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	directory := flag.String("d", "", "Directory")
	flag.Parse()

	if *directory == "" {
		fmt.Println("Please provide a directory to analyse.")
		os.Exit(0)
	}

	// step 1. get all the raw data hand state files

	hss := "hand.sofastates.txt" //

	hssFiles := ListAllFilesRecursivelyByFilename(*directory, hss)
	rbohand2 := regexp.MustCompile(".*/rbohand2/.*")
	rbohand2Files := Select(hssFiles, *rbohand2)
	for _, s := range rbohand2Files[1:2] {
		ConvertHandSofaSate(*directory, s)
	}

}
