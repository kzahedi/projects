package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"regexp"
)

func CalculateCovarianceMatrices(grasp, pattern *regexp.Regexp, directory *string) {
	filename := "difference.hand.sofastates.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	differenceBehaviours := Select(files, *grasp)

	for _, s := range differenceBehaviours {

		d := csv.NewReader(file)
		srecords, err := d.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		data := make([][]float64, len(srecords), len(srecords))
		for i := 0; i < len(srecords); i++ {
			data[i] = make([]float64, len(srecords[0]), len(srecords[0]))
		}

		r := make([][]float64, data.NrOfTrajectories, data.NrOfTrajectories)
		for i := 0; i < data.NrOfTrajectories; i++ {
			r[i] = make([]float64, data.NrOfTrajectories, data.NrOfTrajectories)
			for j := 0; j < data.NrOfTrajectories; j++ {
				r[i][j] = gonum.Covariance(data[:][i], data[:][j], nil)
			}
		}
	}
	fmt.Println(r)
}
