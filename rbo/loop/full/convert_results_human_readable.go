package main

import "fmt"

func getMeanValue(x, y int, data [][]float64) float64 {
	return data[0][x*93+y]
}

func getStdValue(x, y int, data [][]float64) float64 {
	return data[1][x*93+y]
}

func getStrings(label string, indices [][]int, data [][]float64) [][]string {
	var r [][]string

	for _, v := range indices {
		src := v[0]
		dst := v[1]

		xSrcIndex := src*3 + 0
		ySrcIndex := src*3 + 1
		zSrcIndex := src*3 + 2

		xDstIndex := dst*3 + 0
		yDstIndex := dst*3 + 1
		zDstIndex := dst*3 + 2

		s := make([]string, 19, 19)

		s[0] = fmt.Sprintf("%s Frame %d vs. Frame %d", label, src, dst)
		s[1] = fmt.Sprintf("%.3f", getMeanValue(xSrcIndex, xDstIndex, data))
		s[2] = fmt.Sprintf("%.3f", getStdValue(xSrcIndex, xDstIndex, data))
		s[3] = fmt.Sprintf("%.3f", getMeanValue(xSrcIndex, yDstIndex, data))
		s[4] = fmt.Sprintf("%.3f", getStdValue(xSrcIndex, yDstIndex, data))
		s[5] = fmt.Sprintf("%.3f", getMeanValue(xSrcIndex, zDstIndex, data))
		s[6] = fmt.Sprintf("%.3f", getStdValue(xSrcIndex, zDstIndex, data))
		s[7] = fmt.Sprintf("%.3f", getMeanValue(ySrcIndex, xDstIndex, data))
		s[8] = fmt.Sprintf("%.3f", getStdValue(ySrcIndex, xDstIndex, data))
		s[9] = fmt.Sprintf("%.3f", getMeanValue(ySrcIndex, yDstIndex, data))
		s[10] = fmt.Sprintf("%.3f", getStdValue(ySrcIndex, yDstIndex, data))
		s[11] = fmt.Sprintf("%.3f", getMeanValue(ySrcIndex, zDstIndex, data))
		s[12] = fmt.Sprintf("%.3f", getStdValue(ySrcIndex, zDstIndex, data))
		s[13] = fmt.Sprintf("%.3f", getMeanValue(zSrcIndex, xDstIndex, data))
		s[14] = fmt.Sprintf("%.3f", getStdValue(zSrcIndex, xDstIndex, data))
		s[15] = fmt.Sprintf("%.3f", getMeanValue(zSrcIndex, yDstIndex, data))
		s[16] = fmt.Sprintf("%.3f", getStdValue(zSrcIndex, yDstIndex, data))
		s[17] = fmt.Sprintf("%.3f", getMeanValue(zSrcIndex, zDstIndex, data))
		s[18] = fmt.Sprintf("%.3f", getStdValue(zSrcIndex, zDstIndex, data))
		r = append(r, s)
	}
	return r
}

func ConvertIROSMatrixResults(input string) {

	data := ReadCSVToFloat(input)

	thumb := [][]int{{29, 30}, {28, 29}, {27, 28}, {26, 27}, {25, 30}}
	palm := [][]int{{24, 25}, {23, 24}, {22, 23}, {21, 22}, {0, 21}}
	pinky := [][]int{{19, 20}, {18, 19}, {17, 18}, {16, 17}, {0, 16}}
	ring := [][]int{{14, 15}, {13, 14}, {12, 13}, {11, 12}, {0, 11}}
	middle := [][]int{{9, 10}, {8, 9}, {7, 8}, {6, 7}, {0, 6}}
	index := [][]int{{5, 6}, {4, 5}, {3, 4}, {2, 3}, {0, 1}}

	var output [][]string
	s := make([]string, 19, 19)
	s[1] = "x vs. x"
	s[2] = "x vs. x"
	s[3] = "x vs. y"
	s[4] = "x vs. y"
	s[5] = "x vs. z"
	s[6] = "x vs. z"
	s[7] = "y vs. x"
	s[8] = "y vs. x"
	s[9] = "y vs. y"
	s[10] = "y vs. y"
	s[11] = "y vs. z"
	s[12] = "y vs. z"
	s[13] = "z vs. x"
	s[14] = "z vs. x"
	s[15] = "z vs. y"
	s[16] = "z vs. y"
	s[17] = "z vs. z"
	s[18] = "z vs. z"
	output = append(output, s)

	s = make([]string, 19, 19)
	s[1] = "Mean"
	s[2] = "STD"
	s[3] = "Mean"
	s[4] = "STD"
	s[5] = "Mean"
	s[6] = "STD"
	s[7] = "Mean"
	s[8] = "STD"
	s[9] = "Mean"
	s[10] = "STD"
	s[11] = "Mean"
	s[12] = "STD"
	s[13] = "Mean"
	s[14] = "STD"
	s[15] = "Mean"
	s[16] = "STD"
	s[17] = "Mean"
	s[18] = "STD"
	output = append(output, s)

	indexStrigns := getStrings("Index Finger", index, data)
	middleStrigns := getStrings("Middle Finger", middle, data)
	ringStrigns := getStrings("Ring Finger", ring, data)
	pinkyStrigns := getStrings("Pinky Finger", pinky, data)
	palmStrigns := getStrings("Palm", palm, data)
	thumbStrigns := getStrings("Thumb", thumb, data)

	for _, v := range indexStrigns {
		output = append(output, v)
	}
	for _, v := range middleStrigns {
		output = append(output, v)
	}
	for _, v := range ringStrigns {
		output = append(output, v)
	}
	for _, v := range pinkyStrigns {
		output = append(output, v)
	}
	for _, v := range palmStrigns {
		output = append(output, v)
	}
	for _, v := range thumbStrigns {
		output = append(output, v)
	}

	WriteCsvMatrix(input, output)
}

func ConvertSegmentMatrixResults(input string) {

	data := ReadCSVToFloat(input)

	fmt.Println(len(data))
	fmt.Println(len(data[0]))

}
