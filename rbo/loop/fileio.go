package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/westphae/quaternion"
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
		q := quaternion.Quaternion{qw, qx, qy, qz}
		phi, theta, psi := quaternion.Euler(q)

		p := CreatePose(x, y, z, phi, theta, psi)
		r = append(r, p)
	}
	return r
}

func ReadSofaSates(parent, input string) Data {
	data := Data{Trajectories: nil, NrOfTrajectories: 0, NrOfDataPoints: 0}
	fmt.Println("Reading:", input)

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
	fmt.Println(fmt.Sprintf("Writing: %s", filename))
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(stringData)
}

func ReadCSVToData(parent, input string) Data {
	data := Data{Trajectories: nil, NrOfTrajectories: 0, NrOfDataPoints: 0}
	fmt.Println("Reading:", input)

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
			pose := CreatePose(x, y, z, 0.0, 0.0, 0.0)
			data.Trajectories[trajectoryIndex].Frame[frameIndex] = pose
		}
	}
	return data
}

func WriteCSV(filename string, data [][]string) {
	fmt.Println("Writing:", filename)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(file)
	w.WriteAll(data)
}
