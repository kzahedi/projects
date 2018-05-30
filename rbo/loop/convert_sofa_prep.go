package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ConvertSofaStatesPreprocessing(filename string, hands, ctrls []*regexp.Regexp, directory string, convertToWritsFrame bool) {
	fmt.Println("Converting sofa state files:", filename)
	files := ListAllFilesRecursivelyByFilename(directory, filename)

	selectedFiles := SelectFiles(files, hands, ctrls)
	bar := pb.StartNew(len(selectedFiles))

	for _, s := range selectedFiles {
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

	bar.Finish()
}
