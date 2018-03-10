package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStates(filename string, hands, ctrls []*regexp.Regexp, directory *string, convertToWritsFrame bool) {
	fmt.Println("Converting sofa state files:", filename)
	files := ListAllFilesRecursivelyByFilename(*directory, filename)
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
			rbohand2Files := Select(files, *hand)
			rbohand2Files = Select(rbohand2Files, *ctrl)
			for _, s := range rbohand2Files {
				outfile := strings.Replace(s, "raw", "analysis", 1)
				outfile = strings.Replace(outfile, "txt", "csv", 1)
				if _, err := os.Stat(outfile); os.IsNotExist(err) {
					data := ReadSofaSates(s) // returns 2d-array of pose
					if convertToWritsFrame {
						data = transformIntoWristFramePositionOnly(data)
					}
					CreateDir(outfile)
					WritePositions(outfile, data)
				}
				bar.Increment()
			}
		}
	}
	bar.Finish()
}
