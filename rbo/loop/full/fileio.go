package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func convert(data []string) []Pose {
	var r []Pose
	n := len(data) / 7
	d := make([]float64, len(data), len(data))
	for i, v := range data {
		d[i], _ = strconv.ParseFloat(v, 64)
	}

	for i := 0; i < n; i++ {
		x := d[i*7]
		y := d[i*7+1]
		z := d[i*7+2]
		qx := d[i*7+3]
		qy := d[i*7+4]
		qz := d[i*7+5]
		qw := d[i*7+6]

		p := CreatePose(x, y, z, qx, qy, qz, qw)
		r = append(r, p)
	}
	return r
}

func ReadSofaSates(input string) Data {
	data := Data{Trajectories: nil, NrOfTrajectories: 0, NrOfDataPoints: 0}

	file, _ := os.Open(input)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "X=") {
			v := strings.Split(s, " ")[3:]
			d := convert(v)
			if data.Trajectories == nil {
				data.Trajectories = make([]Trajectory, len(d), len(d))
			}
			for i, v := range d {
				data.Trajectories[i].Frame = append(data.Trajectories[i].Frame, v)
			}
		}
	}

	data.NrOfDataPoints = len(data.Trajectories[0].Frame)
	data.NrOfTrajectories = len(data.Trajectories)

	return data
}

func framesToStringSlice(trajectories []Trajectory) [][]string {
	nrOfFrames := len(trajectories[0].Frame)
	r := make([][]string, nrOfFrames, nrOfFrames)
	n := len(trajectories) * 3
	for rowIndex := 0; rowIndex < nrOfFrames; rowIndex++ {
		r[rowIndex] = make([]string, n, n)
		for trajectoryIndex := 0; trajectoryIndex < len(trajectories); trajectoryIndex++ {
			p := trajectories[trajectoryIndex].Frame[rowIndex].Position
			r[rowIndex][trajectoryIndex*3] = fmt.Sprintf("%.3f", p.X)
			r[rowIndex][trajectoryIndex*3+1] = fmt.Sprintf("%.3f", p.Y)
			r[rowIndex][trajectoryIndex*3+2] = fmt.Sprintf("%.3f", p.Z)
		}
	}
	return r
}

func WritePositions(filename string, data Data) {
	stringData := framesToStringSlice(data.Trajectories)
	// fmt.Println(fmt.Sprintf("Writing: %s", filename))
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(stringData)
}

func ReadCSVToData(input string) Data {
	data := Data{Trajectories: nil, NrOfTrajectories: 0, NrOfDataPoints: 0}
	// fmt.Println("Reading:", input)

	file, _ := os.Open(input)
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	data.NrOfDataPoints = len(records)
	data.NrOfTrajectories = len(records[0]) / 3

	data.Trajectories = make([]Trajectory, data.NrOfTrajectories, data.NrOfTrajectories)

	for trajectoryIndex := 0; trajectoryIndex < data.NrOfTrajectories; trajectoryIndex++ {
		data.Trajectories[trajectoryIndex].Frame = make([]Pose, data.NrOfDataPoints, data.NrOfDataPoints)
		for frameIndex := 0; frameIndex < data.NrOfDataPoints; frameIndex++ {
			xs := records[frameIndex][trajectoryIndex*3+0]
			ys := records[frameIndex][trajectoryIndex*3+1]
			zs := records[frameIndex][trajectoryIndex*3+2]
			x, _ := strconv.ParseFloat(xs, 64)
			y, _ := strconv.ParseFloat(ys, 64)
			z, _ := strconv.ParseFloat(zs, 64)
			pose := CreatePose(x, y, z, 0.0, 0.0, 0.0, 1.0)
			data.Trajectories[trajectoryIndex].Frame[frameIndex] = pose
		}
	}
	return data
}

func ReadCSVToArray(input string) []float64 {
	// fmt.Println("Reading:", input)

	file, _ := os.Open(input)
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	data := make([]float64, len(records)*len(records[0]), len(records)*len(records[0]))

	index := 0
	for i := 0; i < len(records); i++ {
		for j := 0; j < len(records[0]); j++ {
			data[index], _ = strconv.ParseFloat(records[i][j], 64)
			index++
		}
	}
	return data
}

// ignores the header line
func ReadControlFile(input string) [][]float64 {
	// fmt.Println("Reading:", input)

	file, _ := os.Open(input)
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	data := make([][]float64, len(records)-1, len(records)-1)

	for i := 0; i < len(records)-1; i++ {
		data[i] = make([]float64, len(records[i+1]), len(records[i+1]))
		for j := 0; j < len(records[i]); j++ {
			data[i][j], _ = strconv.ParseFloat(records[i+1][j], 64)
		}
	}
	return data
}

func ReadCSVToFloat(input string) [][]float64 {
	// fmt.Println("Reading:", input)

	file, _ := os.Open(input)
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	data := make([][]float64, len(records), len(records))

	for i := 0; i < len(records); i++ {
		data[i] = make([]float64, len(records[i]), len(records[i]))
		for j := 0; j < len(records[i]); j++ {
			data[i][j], _ = strconv.ParseFloat(records[i][j], 64)
		}
	}
	return data
}

func WriteCsvFloatMatrix(filename string, data [][]float64) {
	d := make([][]string, len(data), len(data))
	for i := 0; i < len(data); i++ {
		d[i] = make([]string, len(data[i]), len(data[i]))
		for j := 0; j < len(data[i]); j++ {
			d[i][j] = fmt.Sprintf("%f", data[i][j])
		}
	}
	WriteCsvMatrix(filename, d)
}

func WriteCsvFloatVector(filename string, data []float64) {
	d := make([]string, len(data), len(data))
	for i := 0; i < len(data); i++ {
		d[i] = fmt.Sprintf("%f", data[i])
	}
	WriteCsvVector(filename, d)
}

func WriteCsvVector(filename string, data []string) {
	// fmt.Println("Writing:", filename)
	// fmt.Println("Data:", data)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	w.Write(data)
}

func WriteCsvMatrix(filename string, data [][]string) {
	// fmt.Println("Writing:", filename)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	w.WriteAll(data)
}
