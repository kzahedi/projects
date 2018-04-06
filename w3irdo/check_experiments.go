package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/kzahedi/utils"
)

func main() {
	directory := flag.String("d", "/Users/zahedi/projects/mpi/experiments/nmode.w3irdo/from_hydra", "Directory")
	search := flag.String("s", "generation-250.xml", "Directory")
	playback := flag.Bool("p", false, "Playback the best")
	// record := flag.Bool("r", false, "Record video")
	flag.Parse()

	directoryPattern := regexp.MustCompile(".*w3irdo_\\d_module.*")
	dirs := utils.ListDirectoriesByRegexp(*directory, *directoryPattern)

	filePattern := regexp.MustCompile("generation.*.xml")

	var complete []string
	var incomplete []string

	for _, dir := range dirs {
		s := fmt.Sprintf("%s/%s", *directory, dir)
		files := utils.ListFilesByExtension(s, filePattern)
		found := false
		for _, file := range files {
			if strings.Contains(file, *search) {
				found = true
				break
			}
		}
		if found == true {
			complete = append(complete, dir)
			if *playback == true {
				fmt.Println(fmt.Sprintf("replay for %s", dir))
				cmd := exec.Command("/usr/local/bin/nmode", "--replayBest", dir)
				cmd.Dir = *directory
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			incomplete = append(incomplete, dir)
		}
	}

	fmt.Println(">>> Completed experiments")
	for _, v := range complete {
		fmt.Println(v)
	}
	fmt.Println(">> Not completed experiments")
	for _, v := range incomplete {
		fmt.Println(v)
	}
}
