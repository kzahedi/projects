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
	First         Frame
	FirstName     string
	FirstIndex    int
	Second        Frame
	SecondName    string
	SecondIndex   int
	ColumnIndices []int
}

func (fp FramePair) String() string {
	str := ""
	str = fmt.Sprintf("%s%s vs. %s ", str, fp.FirstName, fp.SecondName)
	si := ""
	for _, v := range fp.ColumnIndices {
		si = fmt.Sprintf("%s %d", si, v)
	}
	str = fmt.Sprintf("%s(%d,%d): ", str, fp.FirstIndex, fp.SecondIndex)
	str = fmt.Sprintf("%sIndices:%s\n", str, si)
	return str
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
		ColumnIndices: indices}
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

	fmt.Println(fmt.Sprintf("Number of matrices found = %d", len(bestIndices)))

	cMatrixData = selectRows(cMatrixData, bestIndices)

	fmt.Println(cMatrixData[0][1:3])

	s0 := createSegment(0)
	s1 := createSegment(1)
	s2 := createSegment(2)
	s3 := createSegment(3)
	s4 := createSegment(4)
	s5 := createSegment(5)

	fmt.Println(s0)
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(s3)
	fmt.Println(s4)
	fmt.Println(s5)

	var s0Pairs []FramePair

	s0Pairs = addPair(s0Pairs, s0.Frames[0], s0.Name, s0.Frames[1], s0.Name)
	s0Pairs = addPair(s0Pairs, s0.Frames[1], s0.Name, s0.Frames[2], s0.Name)
	s0Pairs = addPair(s0Pairs, s0.Frames[2], s0.Name, s0.Frames[3], s0.Name)
	s0Pairs = addPair(s0Pairs, s0.Frames[3], s0.Name, s0.Frames[4], s0.Name)

	fmt.Println(s0Pairs)

}
