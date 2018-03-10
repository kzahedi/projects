package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func CalculateDifferenceBehaviour(input, output string, hands, ctrls []*regexp.Regexp, prescriptive *regexp.Regexp, directory *string) {
	fmt.Println("Calculating difference behaviour")
	files := ListAllFilesRecursivelyByFilename(*directory, input)

	iterations := 0
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Files := Select(files, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			iterations += len(rbohand2Files)
		}
	}

	bar := pb.StartNew(iterations)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			rbohand2Grasps := Select(files, *hand)
			rbohand2Grasps = Select(rbohand2Grasps, *ctrl)

			rbohand2Prescriptives := Select(files, *prescriptive)
			rbohand2Prescriptives = Select(rbohand2Prescriptives, *ctrl)
			prescritiveData := ReadCSVToData(rbohand2Prescriptives[0])

			for _, s := range rbohand2Grasps {
				outfile := strings.Replace(s, input, output, 1)
				if _, err := os.Stat(outfile); os.IsNotExist(err) {
					data := ReadCSVToData(s) // returns 2d-array of pose
					diff := calculateDifferencePositionOnlyIROS(data, prescritiveData)
					WritePositions(outfile, diff)
				}
				bar.Increment()
			}
		}
	}
	bar.Finish()
}
