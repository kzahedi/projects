package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/westphae/quaternion"
)

type P3D struct {
	X float64
	Y float64
	Z float64
}

type Pose struct {
	Position    P3D
	Orientation P3D
}

func (p *Pose) String() string {
	return fmt.Sprintf("x:%f y:%f z:%f phi:%f theta:%f psi:%f", p.Position.X, p.Position.Y, p.Position.Z, p.Orientation.X, p.Orientation.Y, p.Orientation.Z)
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
		qw := d[i*7+3]
		qx := d[i*7+4]
		qy := d[i*7+5]
		qz := d[i*7+5]
		q := quaternion.Quaternion{qw, qx, qy, qz}
		phi, theta, psi := quaternion.Euler(q)

		pos := P3D{x, y, z}
		orientation := P3D{phi, theta, psi}

		p := Pose{pos, orientation}
		r = append(r, p)
	}
	return r
}

func ConvertHandSofaSate(parent, input string) {
	output := input
	output = strings.Replace(output, "/raw/", "/", -1)
	output = strings.Replace(output, parent, fmt.Sprintf("%s/%s", parent, "results"), -1)
	output = strings.Replace(output, ".txt", ".csv", -1)
	dir := filepath.Dir(output)
	os.MkdirAll(dir, 0755)

	file, _ := os.Open(input)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, "X=") {
			v := strings.Split(s, " ")[3:]
			d := convert(v)
			fmt.Println(d)
		}
	}

}
