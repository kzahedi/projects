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

func CalculateCovarianceMatrices(input, output string, hands, ctrls []*regexp.Regexp, directory string, max, mode int) {
	fmt.Println("Calculating covariance matrices")
	files := ListAllFilesRecursivelyByFilename(directory, input)

	selectedFiles := SelectFiles(files, hands, ctrls)
	bar := pb.StartNew(len(selectedFiles))

	for _, s := range selectedFiles {
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
	bar.Finish()
}

func calculateCovarianceFrameByFrame(data [][]float64, start, stop int) []string {
	indices := [][]int{
		{0, 1}, // thumb
		{1, 2},
		{2, 3},
		{3, 4},
		{5, 6}, // palm
		{6, 7},
		{7, 8},
		{8, 9},
		{10, 11}, // pinky
		{11, 12},
		{12, 13},
		{13, 14},
		{15, 16}, // ring
		{16, 17},
		{17, 18},
		{18, 19},
		{20, 21}, // middle
		{21, 22},
		{22, 23},
		{23, 24},
		{25, 26}, // index
		{26, 27},
		{27, 28},
		{28, 29}}

	r := make([]string, 9*len(indices), 9*len(indices))
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

func calculateCovarianceSegment(data [][]float64, start, stop int) []string { // checked
	// From ConvertSofaStatesSegment:
	// index finder  coordinate frame 5 in local coordinates of frame 3, then
	// index finder  coordinate frame 3 in local coordinates of frame 1
	// middle finder coordinate frame 10 in local coordinates of frame 8, then
	// middle finder coordinate frame 8 in local coordinates of frame 6
	// ring finder   coordinate frame 15 in local coordinates of frame 13, then
	// ring finder   coordinate frame 13 in local coordinates of frame 11
	// pinky finder  coordinate frame 20 in local coordinates of frame 18, then
	// pinky finder  coordinate frame 18 in local coordinates of frame 16
	// palm          coordinate frame 25 in local coordinate of frame 23, then
	// palm          coordinate frame 23 in local coordinate of frame 21
	// thumb         coordinate frame 30 in local coordinate of frame 28, then
	// thumb         coordinate frame 28 in local coordinate of frame 26
	// hence, as we only want to calculate correlations along segments (e.g. fingers),
	// we calculate the covariance between the x,y,z coordinates of the following
	// coordinate frames
	// index finger : 0, 1
	// middle finger: 2, 3
	// ring finger:   4, 5
	// pinky finger:  6, 7
	// palm:          8, 9
	// thumb:         10, 11

	indices := [][]int{{0, 1}, {2, 3}, {4, 5}, {6, 7}, {8, 9}, {10, 11}}

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
	// rows := len(data)
	cols := len(data[0])

	// fmt.Println(fmt.Sprintf("Rows: %d Cols: %d", rows, cols))

	r := make([]string, cols*cols, cols*cols)
	index := 0
	for col1 := 0; col1 < cols; col1++ {
		for col2 := 0; col2 < cols; col2++ {
			di := getColumn(data, col1)
			dj := getColumn(data, col2)
			di = di[start:stop]
			dj = dj[start:stop]
			r[index] = fmt.Sprintf("%f", stat.Covariance(di, dj, nil))
			index++
		}
	}
	return r
}
