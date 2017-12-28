package main

import (
	"flag"
	"fmt"
)

func Main() {
	clusterIndex := flag.Int("c", 0, "Cluster Index in [0,1,2,3]")
	flag.Parse()

	// best := []float64{-17.5545, 34.4454}
	// radius := 10

	var best []float64
	var radius float64

	switch *clusterIndex {
	case 0:
		best = []float64{-65.2778, -3.31772}
		radius = 15.0
	case 1:
		best = []float64{-17.5545, 34.4454}
		radius = 10.0
	case 2:
		best = []float64{-31.0987, -62.2832}
		radius = 20.0
	case 3:
		best = []float64{13.4884, -24.6339}
		radius = 10.0
	}

	parentDir := "/Users/zahedi/projects/TU.Berlin/experiments/run2017011101/results/abort_after_75/rbohand2-controller0"
	cMatrixFilename := fmt.Sprintf("%s/c.plot.data.csv", parentDir)
	tsneFilename := fmt.Sprintf("%s/t-sne.plot.data.csv", parentDir)

	output := fmt.Sprintf("cluster_%d.csv", *clusterIndex)

	extractCluster(output, parentDir, cMatrixFilename, tsneFilename, best, radius)

}
