package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/gonum/stat"
)

type Frame struct {
	Index int
	X     int
	Y     int
	Z     int
}

type Segment struct {
	Name   string
	Frames []Frame
}

func (s Segment) String() string {
	str := ""
	str = fmt.Sprintf("%sName:    %s\n", str, s.Name)
	si := ""
	sx := ""
	sy := ""
	sz := ""
	for _, v := range s.Frames {
		si = fmt.Sprintf("%s %d", si, v.Index)
		sx = fmt.Sprintf("%s %d", sx, v.X)
		sy = fmt.Sprintf("%s %d", sy, v.Y)
		sz = fmt.Sprintf("%s %d", sz, v.Z)
	}
	str = fmt.Sprintf("%sIndices:%s\n", str, si)
	str = fmt.Sprintf("%sX:      %s\n", str, sx)
	str = fmt.Sprintf("%sY:      %s\n", str, sy)
	str = fmt.Sprintf("%sZ:      %s\n", str, sz)
	return str
}

type FramePair struct {
	First             Frame
	FirstName         string
	FirstIndex        int
	Second            Frame
	SecondName        string
	SecondIndex       int
	ColumnIndices     []int
	Data              [][]float64
	Mean              []float64
	StandardDeviation []float64
}

func (fp FramePair) String() string {
	str := ""
	// str = fmt.Sprintf("%s%s:", str, fp.FirstName)
	// si := ""
	// for _, v := range fp.ColumnIndices {
	// si = fmt.Sprintf("%s %d", si, v)
	// }
	// str = fmt.Sprintf("%s(%d,%d): ", str, fp.FirstIndex, fp.SecondIndex)
	// str = fmt.Sprintf("%s\n Results:\n", str)
	str = fmt.Sprintf("%s%s Frame %d x-axis vs. Frame %d x-axis:  Mean: %8.3f  Std: %8.3f\n", str, fp.FirstName, fp.FirstIndex, fp.SecondIndex, fp.Mean[0], fp.StandardDeviation[0])
	str = fmt.Sprintf("%s%s Frame %d x-axis vs. frame %d y-axis:  Mean: %8.3f  Std: %8.3f\n", str, fp.FirstName, fp.FirstIndex, fp.SecondIndex, fp.Mean[1], fp.StandardDeviation[1])
	str = fmt.Sprintf("%s%s Frame %d x-axis vs. Frame %d z-axis:  Mean: %8.3f  Std: %8.3f\n", str, fp.FirstName, fp.FirstIndex, fp.SecondIndex, fp.Mean[2], fp.StandardDeviation[2])
	str = fmt.Sprintf("%s%s Frame %d y-axis vs. Frame %d x-axis:  Mean: %8.3f  Std: %8.3f\n", str, fp.FirstName, fp.FirstIndex, fp.SecondIndex, fp.Mean[3], fp.StandardDeviation[3])
	str = fmt.Sprintf("%s%s Frame %d y-axis vs. Frame %d y-axis:  Mean: %8.3f  Std: %8.3f\n", str, fp.FirstName, fp.FirstIndex, fp.SecondIndex, fp.Mean[4], fp.StandardDeviation[4])
	str = fmt.Sprintf("%s%s Frame %d y-axis vs. Frame %d z-axis:  Mean: %8.3f  Std: %8.3f\n", str, fp.FirstName, fp.FirstIndex, fp.SecondIndex, fp.Mean[5], fp.StandardDeviation[5])
	str = fmt.Sprintf("%s%s Frame %d z-axis vs. Frame %d x-axis:  Mean: %8.3f  Std: %8.3f\n", str, fp.FirstName, fp.FirstIndex, fp.SecondIndex, fp.Mean[6], fp.StandardDeviation[6])
	str = fmt.Sprintf("%s%s Frame %d z-axis vs. Frame %d y-axis:  Mean: %8.3f  Std: %8.3f\n", str, fp.FirstName, fp.FirstIndex, fp.SecondIndex, fp.Mean[7], fp.StandardDeviation[7])
	str = fmt.Sprintf("%s%s Frame %d z-axis vs. Frame %d z-axis:  Mean: %8.3f  Std: %8.3f\n", str, fp.FirstName, fp.FirstIndex, fp.SecondIndex, fp.Mean[8], fp.StandardDeviation[8])
	// str = fmt.Sprintf("%sIndices:%s\n", str, si)
	// str = fmt.Sprintf("%s Data:\n", str)
	// for ri := 0; ri < len(fp.Data); ri++ {
	// for ci := 0; ci < len(fp.ColumnIndices); ci++ {
	// str = fmt.Sprintf("%s %f", str, fp.Data[ri][ci])
	// }
	// str = fmt.Sprintf("%s\n", str)
	// }
	// str = fmt.Sprintf("%s Mean:\n", str)
	// for _, v := range fp.Mean {
	// str = fmt.Sprintf("%s %f", str, v)
	// }
	// str = fmt.Sprintf("%s\n Standard Deviation:\n", str)
	// for _, v := range fp.StandardDeviation {
	// str = fmt.Sprintf("%s %f", str, v)
	// }

	return str
}

func extractData(pairs []FramePair, cMatrixData [][]float64) []FramePair {
	var r []FramePair
	for _, fp := range pairs {
		fp.Data = make([][]float64, len(cMatrixData), len(cMatrixData))
		for ri := 0; ri < len(cMatrixData); ri++ {
			fp.Data[ri] = make([]float64, len(fp.ColumnIndices), len(fp.ColumnIndices))
			for i := 0; i < len(fp.ColumnIndices); i++ {
				fp.Data[ri][i] = cMatrixData[ri][fp.ColumnIndices[i]]
			}
		}
		r = append(r, fp)
	}
	return r
}

func calculateMeanStd(pairs []FramePair) []FramePair {
	var r []FramePair
	for _, fp := range pairs {
		fp.Mean = make([]float64, len(fp.ColumnIndices), len(fp.ColumnIndices))
		fp.StandardDeviation = make([]float64, len(fp.ColumnIndices), len(fp.ColumnIndices))
		for ci := 0; ci < len(fp.ColumnIndices); ci++ {
			var d []float64
			for ri := 0; ri < len(fp.Data); ri++ {
				d = append(d, fp.Data[ri][ci])
			}
			mean, std := stat.MeanStdDev(d, nil)
			fp.Mean[ci] = mean
			fp.StandardDeviation[ci] = std
		}
		r = append(r, fp)
	}
	return r
}

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

func createSegment(finderIndex int) Segment {
	start := 0
	end := 0
	name := ""
	switch finderIndex {
	case 0:
		start = 0
		end = 4
		name = "Index finger"
	case 1:
		start = 5
		end = 9
		name = "Middle finger"
	case 2:
		start = 10
		end = 14
		name = "Ring finger"
	case 3:
		start = 15
		end = 19
		name = "Pinkie finger"
	case 4: // palm
		start = 20
		end = 25
		name = "Palm"
	case 5: // thumb
		start = 26
		end = 30
		name = "Thumb"
	}
	var f []Frame
	for i := start; i <= end; i++ {
		f = append(f, Frame{i, 3 * i, 3*i + 1, 3*i + 2})
	}
	s := Segment{Name: name, Frames: f}
	return s
}

func addPair(lst []FramePair, a Frame, aName string, b Frame, bName string) []FramePair {
	indices := make([]int, 9, 9)

	indices[0] = a.X*93 + b.X
	indices[1] = a.X*93 + b.Y
	indices[2] = a.X*93 + b.Z

	indices[3] = a.Y*93 + b.X
	indices[4] = a.Y*93 + b.Y
	indices[5] = a.Y*93 + b.Z

	indices[6] = a.Z*93 + b.X
	indices[7] = a.Z*93 + b.Y
	indices[8] = a.Z*93 + b.Z

	f := FramePair{First: a, FirstName: aName, FirstIndex: a.Index,
		Second: b, SecondName: bName, SecondIndex: b.Index,
		ColumnIndices:     indices,
		Data:              nil,
		Mean:              nil,
		StandardDeviation: nil}

	return append(lst, f)
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

	cMatrixData = selectRows(cMatrixData, bestIndices)

	s0 := createSegment(0)
	s1 := createSegment(1)
	s2 := createSegment(2)
	s3 := createSegment(3)
	s4 := createSegment(4)
	s5 := createSegment(5)

	// fmt.Println(s0)
	// fmt.Println(s1)
	// fmt.Println(s2)
	// fmt.Println(s3)
	// fmt.Println(s4)
	// fmt.Println(s5)

	var s0Pairs []FramePair
	var s1Pairs []FramePair
	var s2Pairs []FramePair
	var s3Pairs []FramePair
	var s4Pairs []FramePair
	var s5Pairs []FramePair

	for i := 0; i < len(s0.Frames)-1; i++ {
		s0Pairs = addPair(s0Pairs, s0.Frames[i], s0.Name, s0.Frames[i+1], s0.Name)
	}

	for i := 0; i < len(s1.Frames)-1; i++ {
		s1Pairs = addPair(s1Pairs, s1.Frames[i], s1.Name, s1.Frames[i+1], s1.Name)
	}

	for i := 0; i < len(s2.Frames)-1; i++ {
		s2Pairs = addPair(s2Pairs, s2.Frames[i], s2.Name, s2.Frames[i+1], s2.Name)
	}

	for i := 0; i < len(s3.Frames)-1; i++ {
		s3Pairs = addPair(s3Pairs, s3.Frames[i], s3.Name, s3.Frames[i+1], s3.Name)
	}

	for i := 0; i < len(s4.Frames)-1; i++ {
		s4Pairs = addPair(s4Pairs, s4.Frames[i], s4.Name, s4.Frames[i+1], s4.Name)
	}

	for i := 0; i < len(s5.Frames)-1; i++ {
		s5Pairs = addPair(s5Pairs, s5.Frames[i], s5.Name, s5.Frames[i+1], s5.Name)
	}

	s0Pairs = extractData(s0Pairs, cMatrixData)
	s1Pairs = extractData(s1Pairs, cMatrixData)
	s2Pairs = extractData(s2Pairs, cMatrixData)
	s3Pairs = extractData(s3Pairs, cMatrixData)
	s4Pairs = extractData(s4Pairs, cMatrixData)
	s5Pairs = extractData(s5Pairs, cMatrixData)

	s0Pairs = calculateMeanStd(s0Pairs)
	s1Pairs = calculateMeanStd(s1Pairs)
	s2Pairs = calculateMeanStd(s2Pairs)
	s3Pairs = calculateMeanStd(s3Pairs)
	s4Pairs = calculateMeanStd(s4Pairs)
	s5Pairs = calculateMeanStd(s5Pairs)

	for _, v := range s0Pairs {
		fmt.Println(v)
		fmt.Println()
	}

	for _, v := range s1Pairs {
		fmt.Println(v)
		fmt.Println()
	}

	for _, v := range s2Pairs {
		fmt.Println(v)
		fmt.Println()
	}

	for _, v := range s3Pairs {
		fmt.Println(v)
		fmt.Println()
	}

	for _, v := range s4Pairs {
		fmt.Println(v)
		fmt.Println()
	}

	for _, v := range s5Pairs {
		fmt.Println(v)
		fmt.Println()
	}

}
