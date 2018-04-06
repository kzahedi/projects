package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gonum/stat"
	pb "gopkg.in/cheggaaa/pb.v1"
)

const ( // iota is reset to 0
	MODE_FULL           = iota // c0 == 0
	MODE_FRAME_BY_FRAME = iota // c1 == 1
	MODE_SEGMENT        = iota // c2 == 2
)

func CalculateCovarianceMatrices(input, output string, hands, ctrls []*regexp.Regexp, directory *string, max, mode int) {
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

					start := 10
					stop := 10 + max
					var covarianceMatrix []string
					switch mode {
					case MODE_FRAME_BY_FRAME:
						covarianceMatrix = calculateCovarianceFrameByFrame(data, start, stop)
					case MODE_SEGMENT:
						covarianceMatrix = calculateCovarianceTipToRoot(data, start, stop)
					case MODE_FULL:
						covarianceMatrix = calculateCovarianceFull(data, start, stop)
					default:
						panic("Unknown mode")
					}
					WriteCsvVector(outfile, covarianceMatrix)
				}
				bar.Increment()
			}
		}
	}
	bar.Finish()
}

func calculateCovarianceFrameByFrame(data [][]float64, start, stop int) []string {
	cols := len(data[0])
	indices := [][]int{
		{28, 29}, // thumb and palm
		{27, 28},
		{26, 27},
		{25, 26},
		{24, 29},
		{23, 24},
		{22, 23},
		{21, 22},
		{20, 21},
		{18, 19}, // pinky finger
		{17, 18},
		{16, 17},
		{15, 16},
		{13, 14}, // ring finger
		{12, 13},
		{11, 12},
		{10, 11},
		{8, 9}, // middle finger
		{7, 8},
		{6, 7},
		{5, 6},
		{4, 5}, // thumb
		{3, 4},
		{2, 3},
		{1, 2}}

	r := make([]string, len(indices), len(indices))
	for i, v := range indices {
		for j := 0; j < cols; j++ {
			di := getColumn(data, v[0])
			dj := getColumn(data, v[1])
			di = di[start:stop]
			dj = dj[start:stop]
			r[i] = fmt.Sprintf("%f", stat.Covariance(di, dj, nil))
		}
	}
	return r
}

func calculateCovarianceTipToRoot(data [][]float64, start, stop int) []string {
	cols := len(data[0])
	indices := [][]int{{24, 29}, {20, 24}, {15, 19}, {10, 14}, {5, 9}, {0, 4}}

	r := make([]string, len(indices), len(indices))
	for i, v := range indices {
		for j := 0; j < cols; j++ {
			di := getColumn(data, v[0])
			dj := getColumn(data, v[1])
			di = di[start:stop]
			dj = dj[start:stop]
			r[i] = fmt.Sprintf("%f", stat.Covariance(di, dj, nil))
		}
	}
	return r
}

func calculateCovarianceFull(data [][]float64, start, stop int) []string {
	rows := len(data)
	cols := len(data[0])

	r := make([]string, rows*cols, rows*cols)
	index := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			di := getColumn(data, row)
			dj := getColumn(data, col)
			di = di[start:stop]
			dj = dj[start:stop]
			r[index] = fmt.Sprintf("%f", stat.Covariance(di, dj, nil))
			index++
		}
	}
	return r
}
