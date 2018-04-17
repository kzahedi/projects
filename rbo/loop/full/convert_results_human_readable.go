package main

import "fmt"

func ConvertIROSMatrixResults(input string) {

	data := ReadCSVToFloat(input)

	thumb := [][]int{{29, 30}, {28, 29}, {27, 28}, {26, 27}, {25, 30}}
	palm := [][]int{{24, 25}, {23, 24}, {22, 23}, {21, 22}, {0, 21}}
	pinky := [][]int{{19, 20}, {18, 19}, {17, 18}, {16, 17}, {0, 16}}
	ring := [][]int{{14, 15}, {13, 14}, {12, 13}, {11, 12}, {0, 11}}
	middle := [][]int{{9, 10}, {8, 9}, {7, 8}, {6, 7}, {0, 6}}
	index := [][]int{{5, 6}, {4, 5}, {3, 4}, {2, 3}, {0, 1}}

	n := len(thumb) + len(palm) + len(pinky) + len(ring) + len(middle) + len(index)

	output := make([][]string, 3*n+1, 3*n+1)
	output[0] = make([]string, 10, 10)
	output[0][1] = "x vs. x"
	output[0][2] = "x vs. y"
	output[0][3] = "x vs. z"
	output[0][4] = "y vs. x"
	output[0][5] = "y vs. y"
	output[0][6] = "y vs. z"
	output[0][7] = "z vs. x"
	output[0][8] = "z vs. y"
	output[0][9] = "z vs. z"

	for i, v := range thumb {
		src := v[0]
		dst := v[1]

		xSrcIndex := src*3 + 0
		ySrcIndex := src*3 + 1
		zSrcIndex := src*3 + 2

		xDstIndex := dst*3 + 0
		yDstIndex := dst*3 + 1
		zDstIndex := dst*3 + 2

		s := make([]string, 10, 10)

		s[0] = fmt.Sprintf("Thumb Frame %d vs. Frame %d", src, dst)
		s[1] = fmt.Sprintf("%.3f", data[xSrcIndex][xDstIndex])
		s[2] = fmt.Sprintf("%.3f", data[xSrcIndex][yDstIndex])
		s[3] = fmt.Sprintf("%.3f", data[xSrcIndex][zDstIndex])
		s[4] = fmt.Sprintf("%.3f", data[ySrcIndex][xDstIndex])
		s[5] = fmt.Sprintf("%.3f", data[ySrcIndex][yDstIndex])
		s[6] = fmt.Sprintf("%.3f", data[ySrcIndex][zDstIndex])
		s[7] = fmt.Sprintf("%.3f", data[zSrcIndex][xDstIndex])
		s[8] = fmt.Sprintf("%.3f", data[zSrcIndex][yDstIndex])
		s[9] = fmt.Sprintf("%.3f", data[zSrcIndex][zDstIndex])
		output[i+1] = s
	}

	WriteCsvMatrix(input, output)

}
