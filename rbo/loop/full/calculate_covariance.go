package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gonum/stat"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func CalculateCovarianceMatrices(input, output string, hands, ctrls []*regexp.Regexp, directory *string, max int) {
	fmt.Println("Calculating covariance matrices")
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
			differenceBehaviours := Select(files, *hand)
			differenceBehaviours = Select(differenceBehaviours, *ctrl)

			for _, s := range differenceBehaviours {
				outfile := strings.Replace(s, input, output, 1)
				if _, err := os.Stat(outfile); os.IsNotExist(err) {
					data := ReadCSVToFloat(s)
					cols := len(data[0])

					start := 10
					stop := 10 + max
					r := make([][]string, cols, cols)
					for i := 0; i < cols; i++ {
						r[i] = make([]string, cols, cols)
						for j := 0; j < cols; j++ {
							di := getColumn(data, i)
							dj := getColumn(data, j)
							di = di[start:stop]
							dj = dj[start:stop]
							r[i][j] = fmt.Sprintf("%f", stat.Covariance(di, dj, nil))
						}
					}
					WriteCSV(outfile, r)
				}
				bar.Increment()
			}
		}
	}
	bar.Finish()
}
