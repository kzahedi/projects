package main

import "fmt"

func getMeanValueFrameByFrame(x int, data [][]float64) float64 {
	return data[0][x]
}

func getStdValueFrameByFrame(x int, data [][]float64) float64 {
	return data[1][x]
}

func getStringsFrameByFrame(label string, indices []int, data [][]float64) [][]string {
	var r [][]string

	for i, v := range indices {
		s := make([]string, 19, 19)

		s[0] = fmt.Sprintf("%s Frame %d vs. Frame %d", label, i, i+1)
		s[1] = fmt.Sprintf("%.3f", getMeanValueFrameByFrame(9*v+0, data))
		s[2] = fmt.Sprintf("%.3f", getStdValueFrameByFrame(9*v+0, data))
		s[3] = fmt.Sprintf("%.3f", getMeanValueFrameByFrame(9*v+1, data))
		s[4] = fmt.Sprintf("%.3f", getStdValueFrameByFrame(9*v+1, data))
		s[5] = fmt.Sprintf("%.3f", getMeanValueFrameByFrame(9*v+2, data))
		s[6] = fmt.Sprintf("%.3f", getStdValueFrameByFrame(9*v+2, data))
		s[7] = fmt.Sprintf("%.3f", getMeanValueFrameByFrame(9*v+3, data))
		s[8] = fmt.Sprintf("%.3f", getStdValueFrameByFrame(9*v+3, data))
		s[9] = fmt.Sprintf("%.3f", getMeanValueFrameByFrame(9*v+4, data))
		s[10] = fmt.Sprintf("%.3f", getStdValueFrameByFrame(9*v+4, data))
		s[11] = fmt.Sprintf("%.3f", getMeanValueFrameByFrame(9*v+5, data))
		s[12] = fmt.Sprintf("%.3f", getStdValueFrameByFrame(9*v+5, data))
		s[13] = fmt.Sprintf("%.3f", getMeanValueFrameByFrame(9*v+6, data))
		s[14] = fmt.Sprintf("%.3f", getStdValueFrameByFrame(9*v+6, data))
		s[15] = fmt.Sprintf("%.3f", getMeanValueFrameByFrame(9*v+7, data))
		s[16] = fmt.Sprintf("%.3f", getStdValueFrameByFrame(9*v+7, data))
		s[17] = fmt.Sprintf("%.3f", getMeanValueFrameByFrame(9*v+8, data))
		s[18] = fmt.Sprintf("%.3f", getStdValueFrameByFrame(9*v+8, data))
		r = append(r, s)
	}
	return r
}

func getMeanValueIROS(x, y int, data [][]float64) float64 {
	return data[0][x*93+y]
}

////////////////////////////////////////////////////////////////////////////////
// IROS
////////////////////////////////////////////////////////////////////////////////

func getStdValueIROS(x, y int, data [][]float64) float64 {
	return data[1][x*93+y]
}

func getStringsIROS(label string, indices [][]int, data [][]float64) [][]string {
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

		// fmt.Println(fmt.Sprintf("Indices [%d %d %d] [%d %d %d]", xSrcIndex, ySrcIndex, zSrcIndex, xDstIndex, yDstIndex, zDstIndex))

		s := make([]string, 19, 19)

		s[0] = fmt.Sprintf("%s Frame %d vs. Frame %d", label, src, dst)
		s[1] = fmt.Sprintf("%.3f", getMeanValueIROS(xSrcIndex, xDstIndex, data))
		s[2] = fmt.Sprintf("%.3f", getStdValueIROS(xSrcIndex, xDstIndex, data))
		s[3] = fmt.Sprintf("%.3f", getMeanValueIROS(xSrcIndex, yDstIndex, data))
		s[4] = fmt.Sprintf("%.3f", getStdValueIROS(xSrcIndex, yDstIndex, data))
		s[5] = fmt.Sprintf("%.3f", getMeanValueIROS(xSrcIndex, zDstIndex, data))
		s[6] = fmt.Sprintf("%.3f", getStdValueIROS(xSrcIndex, zDstIndex, data))
		s[7] = fmt.Sprintf("%.3f", getMeanValueIROS(ySrcIndex, xDstIndex, data))
		s[8] = fmt.Sprintf("%.3f", getStdValueIROS(ySrcIndex, xDstIndex, data))
		s[9] = fmt.Sprintf("%.3f", getMeanValueIROS(ySrcIndex, yDstIndex, data))
		s[10] = fmt.Sprintf("%.3f", getStdValueIROS(ySrcIndex, yDstIndex, data))
		s[11] = fmt.Sprintf("%.3f", getMeanValueIROS(ySrcIndex, zDstIndex, data))
		s[12] = fmt.Sprintf("%.3f", getStdValueIROS(ySrcIndex, zDstIndex, data))
		s[13] = fmt.Sprintf("%.3f", getMeanValueIROS(zSrcIndex, xDstIndex, data))
		s[14] = fmt.Sprintf("%.3f", getStdValueIROS(zSrcIndex, xDstIndex, data))
		s[15] = fmt.Sprintf("%.3f", getMeanValueIROS(zSrcIndex, yDstIndex, data))
		s[16] = fmt.Sprintf("%.3f", getStdValueIROS(zSrcIndex, yDstIndex, data))
		s[17] = fmt.Sprintf("%.3f", getMeanValueIROS(zSrcIndex, zDstIndex, data))
		s[18] = fmt.Sprintf("%.3f", getStdValueIROS(zSrcIndex, zDstIndex, data))
		r = append(r, s)
	}
	return r
}

func ConvertIROSMatrixResults(dir, input, out string) {

	data := ReadCSVToFloat(fmt.Sprintf("%s/%s", dir, input))

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

	indexStrings := getStringsIROS("Index Finger", index, data)
	middleStrings := getStringsIROS("Middle Finger", middle, data)
	ringStrings := getStringsIROS("Ring Finger", ring, data)
	pinkyStrings := getStringsIROS("Pinky Finger", pinky, data)
	palmStrings := getStringsIROS("Palm", palm, data)
	thumbStrings := getStringsIROS("Thumb", thumb, data)

	for _, v := range indexStrings {
		output = append(output, v)
	}
	for _, v := range middleStrings {
		output = append(output, v)
	}
	for _, v := range ringStrings {
		output = append(output, v)
	}
	for _, v := range pinkyStrings {
		output = append(output, v)
	}
	for _, v := range palmStrings {
		output = append(output, v)
	}
	for _, v := range thumbStrings {
		output = append(output, v)
	}

	WriteCsvMatrix(fmt.Sprintf("%s/%s", dir, out), output)
}

////////////////////////////////////////////////////////////////////////////////
// Segment
////////////////////////////////////////////////////////////////////////////////

func getStringsSegment(label string, index int, data [][]float64) []string {
	s := make([]string, 19, 19)

	s[0] = fmt.Sprintf("%s", label)
	j := 1

	for i := 0; i < 9; i++ {
		s[j] = fmt.Sprintf("%.3f", data[0][index*9+i])
		j++
		s[j] = fmt.Sprintf("%.3f", data[1][index*9+i])
		j++
	}

	return s
}

func ConvertSegmentMatrixResults(dir, input, out string) {

	data := ReadCSVToFloat(fmt.Sprintf("%s/%s", dir, input))

	index := 5
	middle := 4
	ring := 3
	pinky := 2
	palm := 1
	thumb := 0

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

	indexString := getStringsSegment("Index Finger", index, data)
	middleString := getStringsSegment("Middle Finger", middle, data)
	ringString := getStringsSegment("Ring Finger", ring, data)
	pinkyString := getStringsSegment("Pinky Finger", pinky, data)
	palmString := getStringsSegment("Palm", palm, data)
	thumbString := getStringsSegment("Thumb", thumb, data)

	output = append(output, indexString)
	output = append(output, middleString)
	output = append(output, ringString)
	output = append(output, pinkyString)
	output = append(output, palmString)
	output = append(output, thumbString)

	WriteCsvMatrix(fmt.Sprintf("%s/%s", dir, out), output)
}

////////////////////////////////////////////////////////////////////////////////
// Frame by Frame
////////////////////////////////////////////////////////////////////////////////

func ConvertFrameByFrameMatrixResults(dir, input, out string) {

	data := ReadCSVToFloat(fmt.Sprintf("%s/%s", dir, input))

	thumb := []int{0, 1, 2, 3}
	palm := []int{4, 5, 6, 7}
	pinky := []int{8, 9, 10, 11}
	ring := []int{12, 13, 14, 15}
	middle := []int{16, 17, 18, 19}
	index := []int{20, 21, 22, 23}

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

	indexStrings := getStringsFrameByFrame("Index Finger", index, data)
	middleStrings := getStringsFrameByFrame("Middle Finger", middle, data)
	ringStrings := getStringsFrameByFrame("Ring Finger", ring, data)
	pinkyStrings := getStringsFrameByFrame("Pinky Finger", pinky, data)
	palmStrings := getStringsFrameByFrame("Palm", palm, data)
	thumbStrings := getStringsFrameByFrame("Thumb", thumb, data)

	for _, v := range indexStrings {
		output = append(output, v)
	}
	for _, v := range middleStrings {
		output = append(output, v)
	}
	for _, v := range ringStrings {
		output = append(output, v)
	}
	for _, v := range pinkyStrings {
		output = append(output, v)
	}
	for _, v := range palmStrings {
		output = append(output, v)
	}
	for _, v := range thumbStrings {
		output = append(output, v)
	}

	WriteCsvMatrix(fmt.Sprintf("%s/%s", dir, out), output)
}
