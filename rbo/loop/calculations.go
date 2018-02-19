package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gonum/stat"
)

func CalculateCovarianceMatrices(pattern *regexp.Regexp, directory *string) {
	filename := "difference.hand.sofastates.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	differenceBehaviours := Select(files, *pattern)

	for _, s := range differenceBehaviours {

		file, _ := os.Open(s)
		defer file.Close()
		d := csv.NewReader(file)
		records, err := d.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		rows := len(records)
		cols := len(records[0])

		fmt.Println(rows, " ", cols)

		data := make([][]float64, rows, rows)
		for i := 0; i < rows; i++ {
			data[i] = make([]float64, cols, cols)
			for j := 0; j < cols; j++ {
				data[i][j], _ = strconv.ParseFloat(records[i][j], 64)
			}
		}

		r := make([][]string, cols, cols)
		for i := 0; i < cols; i++ {
			r[i] = make([]string, cols, cols)
			for j := 0; j < cols; j++ {
				r[i][j] = fmt.Sprintf("%f", stat.Covariance(data[:][i], data[:][j], nil))
			}
		}

		output := strings.Replace(s, filename, "covariance.csv", 1)
		WriteCSV(output, r)
	}
}
