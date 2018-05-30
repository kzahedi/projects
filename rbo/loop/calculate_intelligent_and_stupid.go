package main

import (
	"fmt"

	"github.com/gonum/stat"
)

func AnalyseIntelligent(data Results, parent, input, output string) [][]float64 {
	var intelligent []string

	for key, value := range data {
		if value.Intelligent && value.SelectedForAnalysis {
			intelligent = append(intelligent, key)
		}
	}

	filename := fmt.Sprintf("%s/%s/analysis/%s", parent, intelligent[0], input)
	matrix := ReadCSVToArray(filename)

	nrOfEntries := len(matrix)
	nrOfMatrices := len(intelligent)

	matrices := make([][]float64, nrOfMatrices, nrOfMatrices)

	meanAndStd := make([][]float64, 2, 2)
	meanAndStd[0] = make([]float64, nrOfEntries, nrOfEntries)
	meanAndStd[1] = make([]float64, nrOfEntries, nrOfEntries)

	for i := 0; i < nrOfMatrices; i++ {
		filename = fmt.Sprintf("%s/%s/analysis/%s", parent, intelligent[i], input)
		matrices[i] = make([]float64, nrOfEntries, nrOfEntries)
		matrices[i] = ReadCSVToArray(filename)
	}

	for col := 0; col < nrOfEntries; col++ {
		d := make([]float64, nrOfMatrices, nrOfMatrices)
		for row := 0; row < nrOfMatrices; row++ {
			d[row] = matrices[row][col]
		}
		m, s := stat.MeanStdDev(d, nil)
		meanAndStd[0][col] = m
		meanAndStd[1][col] = s
	}

	WriteCsvFloatMatrix(output, meanAndStd)

	return meanAndStd
}

func AnalyseStupid(data Results, parent, input, output string) [][]float64 {
	var stupid []string

	for key, value := range data {
		if value.Stupid && value.SelectedForAnalysis {
			stupid = append(stupid, key)
		}
	}

	filename := fmt.Sprintf("%s/%s/analysis/%s", parent, stupid[0], input)
	matrix := ReadCSVToArray(filename)

	nrOfEntries := len(matrix)
	nrOfMatrices := len(stupid)

	matrices := make([][]float64, nrOfMatrices, nrOfMatrices)

	meanAndStd := make([][]float64, 2, 2)
	meanAndStd[0] = make([]float64, nrOfEntries, nrOfEntries)
	meanAndStd[1] = make([]float64, nrOfEntries, nrOfEntries)

	for i := 0; i < nrOfMatrices; i++ {
		filename = fmt.Sprintf("%s/%s/analysis/%s", parent, stupid[i], input)
		matrices[i] = make([]float64, nrOfEntries, nrOfEntries)
		matrices[i] = ReadCSVToArray(filename)
	}

	for col := 0; col < nrOfEntries; col++ {
		d := make([]float64, nrOfMatrices, nrOfMatrices)
		for row := 0; row < nrOfMatrices; row++ {
			d[row] = matrices[row][col]
		}
		m, s := stat.MeanStdDev(d, nil)
		meanAndStd[0][col] = m
		meanAndStd[1][col] = s
	}

	WriteCsvFloatMatrix(output, meanAndStd)

	return meanAndStd
}
