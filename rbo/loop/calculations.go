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
	"github.com/sacado/tsne4go"
)

func CalculateCovarianceMatrices(hand, controller *regexp.Regexp, directory *string, max int) {
	filename := "difference.hand.sofastates.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	differenceBehaviours := Select(files, *hand)
	differenceBehaviours = Select(differenceBehaviours, *controller)

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

		data := make([][]float64, rows, rows)
		for i := 0; i < rows; i++ {
			data[i] = make([]float64, cols, cols)
			for j := 0; j < cols; j++ {
				data[i][j], _ = strconv.ParseFloat(records[i][j], 64)
			}
		}

		start := 10
		stop := 10 + max
		r := make([][]string, cols, cols)
		for i := 0; i < cols; i++ {
			r[i] = make([]string, cols, cols)
			for j := 0; j < cols; j++ {
				di := data[:][i]
				dj := data[:][j]
				di = di[start:stop]
				dj = dj[start:stop]
				r[i][j] = fmt.Sprintf("%f", stat.Covariance(di, dj, nil))
			}
		}

		output := strings.Replace(s, filename, "covariance.csv", 1)
		WriteCSV(output, r)
	}
}

func CalculateTSNE(hand, controller *regexp.Regexp, directory *string) {
	filename := "covariance.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	covariances := Select(files, *hand)
	covariances = Select(covariances, *controller)

	var data tsne4go.VectorDistancer
	data = make([][]float64, len(covariances), len(covariances))
	for i, f := range covariances {
		data[i] = ReadCSVToArray(f)
	}

	tsne := tsne4go.New(data, nil)

	for i := 0; i < 5000; i++ {
		fmt.Println(tsne.Step())
	}

	xy := make([][]float64, len(covariances), len(covariances))
	for i := 0; i < len(covariances); i++ {
		xy[i] = make([]float64, 2, 2)
		xy[i][0] = tsne.Solution[i][0]
		xy[i][1] = tsne.Solution[i][1]
	}

	WriteCSVFloat("/Users/zahedi/Desktop/out.csv", xy)

}
