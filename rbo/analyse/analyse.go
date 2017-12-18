package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func readCsvData(filename string) [][]float64 {
	f, _ := os.Open(filename)
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	sdata, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	data := make([][]float64, len(sdata), len(sdata))
	for i := 0; i < len(sdata); i++ {
		data[i] = make([]float64, len(sdata[i]), len(sdata[i]))
	}

	for i := 0; i < len(sdata); i++ {
		for j := 0; j < len(sdata[i]); j++ {
			v, converr := strconv.ParseFloat(sdata[i][j], 64)
			data[i][j] = v
			if converr != nil {
				log.Fatal(converr)
			}
		}
	}

	return data
}

func dist(a, b []float64) float64 {
	distx := a[0] - b[0]
	disty := a[1] - b[1]
	return math.Sqrt(distx*distx + disty*disty)
}

func selectRows(data [][]float64, indices []int) [][]float64 {
	r := make([][]float64, len(indices), len(indices))
	for i, v := range indices {
		r[i] = data[v]
	}
	return r
}

func generateFingerIndices(finderIndex int) []int {
	start := 0
	end := 0
	switch finderIndex {
	case 1:
		start = 1
		end = 5
	case 2:
		start = 6
		end = 10
	case 3:
		start = 11
		end = 15
	case 4:
		start = 16
		end = 20
	case 5: // palm
		start = 21
		end = 26
	case 6: // thumb
		start = 27
		end = 31
	}
	var r []int
	for i := start; i <= end; i++ {
		r = append(r, i)
	}
	return r
}

func main() {
	parentDir := "/Users/zahedi/projects/TU.Berlin/experiments/run2017011101/results/abort_after_75/rbohand2-controller0"
	cMatrixFilename := fmt.Sprintf("%s/c.plot.data.csv", parentDir)
	cMatrixData := readCsvData(cMatrixFilename)

	tsneFilename := fmt.Sprintf("%s/t-sne.plot.data.csv", parentDir)
	tsneData := readCsvData(tsneFilename)

	// best := []float64{-17.5545, 34.4454}
	// radius := 10

	best := []float64{-65.2778, -3.31772}
	radius := 15.0

	var bestIndices []int

	for i, v := range tsneData {
		if dist(best, v) <= radius {
			bestIndices = append(bestIndices, i)
		}
	}

	fmt.Println(fmt.Sprintf("Number of matrices found = %d", len(bestIndices)))

	cMatrixData = selectRows(cMatrixData, bestIndices)

	fmt.Println(cMatrixData[0][1:3])

	fmt.Println(generateFingerIndices(1))
	fmt.Println(generateFingerIndices(2))
	fmt.Println(generateFingerIndices(3))
	fmt.Println(generateFingerIndices(4))
	fmt.Println(generateFingerIndices(5))
	fmt.Println(generateFingerIndices(6))

}
