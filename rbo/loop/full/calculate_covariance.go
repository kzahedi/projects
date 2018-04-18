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
						covarianceMatrix = calculateCovarianceSegment(data, start, stop)
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

	r := make([]string, 3*len(indices), 3*len(indices))
	for i, v := range indices {
		dix := getColumn(data, 3*v[0]+0)
		djx := getColumn(data, 3*v[1]+0)
		dix = dix[start:stop]
		djx = djx[start:stop]

		diy := getColumn(data, 3*v[0]+1)
		djy := getColumn(data, 3*v[1]+1)
		diy = diy[start:stop]
		djy = djy[start:stop]

		diz := getColumn(data, 3*v[0]+2)
		djz := getColumn(data, 3*v[1]+2)
		diz = diz[start:stop]
		djz = djz[start:stop]

		r[3*i] = fmt.Sprintf("%f", stat.Covariance(dix, djx, nil))
		r[3*i+1] = fmt.Sprintf("%f", stat.Covariance(diy, djy, nil))
		r[3*i+3] = fmt.Sprintf("%f", stat.Covariance(diz, djz, nil))
	}
	return r
}

func calculateCovarianceSegment(data [][]float64, start, stop int) []string {
	// data is pruned in ConvertSofaStatesSegment
	indices := [][]int{{0, 1}, {2, 3}, {4, 5}, {6, 7}, {8, 9}, {10, 11}}

	// fmt.Println("Data: ", len(data))
	// fmt.Println("Data[0]: ", len(data[0]))
	// fmt.Println("Indices: ", len(indices))

	r := make([]string, 9*len(indices), 9*len(indices))
	for i, v := range indices {
		// fmt.Println(3*v[0]+0, 3*v[1]+0)
		dix := getColumn(data, 3*v[0]+0)
		djx := getColumn(data, 3*v[1]+0)
		dix = dix[start:stop]
		djx = djx[start:stop]

		// fmt.Println(3*v[0]+1, 3*v[1]+1)
		diy := getColumn(data, 3*v[0]+1)
		djy := getColumn(data, 3*v[1]+1)
		diy = diy[start:stop]
		djy = djy[start:stop]

		// fmt.Println(3*v[0]+2, 3*v[1]+2)
		diz := getColumn(data, 3*v[0]+2)
		djz := getColumn(data, 3*v[1]+2)
		diz = diz[start:stop]
		djz = djz[start:stop]

		r[9*i+0] = fmt.Sprintf("%f", stat.Covariance(dix, djx, nil))
		r[9*i+1] = fmt.Sprintf("%f", stat.Covariance(dix, djy, nil))
		r[9*i+2] = fmt.Sprintf("%f", stat.Covariance(dix, djz, nil))
		r[9*i+3] = fmt.Sprintf("%f", stat.Covariance(diy, djx, nil))
		r[9*i+4] = fmt.Sprintf("%f", stat.Covariance(diy, djy, nil))
		r[9*i+5] = fmt.Sprintf("%f", stat.Covariance(diy, djz, nil))
		r[9*i+6] = fmt.Sprintf("%f", stat.Covariance(diz, djx, nil))
		r[9*i+7] = fmt.Sprintf("%f", stat.Covariance(diz, djy, nil))
		r[9*i+8] = fmt.Sprintf("%f", stat.Covariance(diz, djz, nil))
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
