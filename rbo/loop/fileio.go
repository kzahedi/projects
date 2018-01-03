package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/westphae/quaternion"
)

// data.Trajectories[trajectoryIndex].Frame[frameIndex].X = convertX(data, trajectoryIndex, frameIndex)

func convertX(data Data, trajectoryIndex, frameIndex int) float64 {
	previous := data.Trajectories[trajectoryIndex].Frame[frameIndex-1].Orientation.X
	current := data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.X
	if previous < 0.0 && current > 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + current - math.Pi
	}
	if previous > 0.0 && current < 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + math.Pi + current
	}
	return current
}

func convertY(data Data, trajectoryIndex, frameIndex int) float64 {
	previous := data.Trajectories[trajectoryIndex].Frame[frameIndex-1].Orientation.Y
	current := data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.Y
	if previous < 0.0 && current > 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + current - math.Pi
	}
	if previous > 0.0 && current < 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + math.Pi + current
	}
	return current
}

func convertZ(data Data, trajectoryIndex, frameIndex int) float64 {
	previous := data.Trajectories[trajectoryIndex].Frame[frameIndex-1].Orientation.Z
	current := data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.Z
	if previous < 0.0 && current > 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + current - math.Pi
	}
	if previous > 0.0 && current < 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + math.Pi + current
	}
	return current
}

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

		pos := P3D{x, y, z}
		orientation := P3D{phi, theta, psi}

		p := Pose{Position: pos, Orientation: orientation}
		r = append(r, p)
	}
	return r
}

func ReadSofaSates(parent, input string) Data {
	data := Data{Trajectories: nil, NrOfTrajectories: 0, NrOfDataPoints: 0}
	fmt.Println(input)
	dir := input
	dir = strings.Replace(dir, "/raw/", "/", -1)
	dir = strings.Replace(dir, parent, fmt.Sprintf("%s/%s", parent, "results"), -1)
	dir = filepath.Dir(dir)
	os.MkdirAll(dir, 0755)

	file, _ := os.Open(input)

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

	// make angles continuous without jump from +pi to -pi and vice versa
	for trajectoryIndex := 0; trajectoryIndex < data.NrOfTrajectories; trajectoryIndex++ {
		for frameIndex := 1; frameIndex < data.NrOfDataPoints; frameIndex++ {
			data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.X = convertX(data, trajectoryIndex, frameIndex)
			data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.Y = convertY(data, trajectoryIndex, frameIndex)
			data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.Z = convertZ(data, trajectoryIndex, frameIndex)
		}
	}

	return data
}
