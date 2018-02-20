package main

import (
	"regexp"

	"github.com/sacado/tsne4go"
)

func AnalyseData(hands, ctrls []*regexp.Regexp, directory *string) {
	filename := "covariance.csv"
	cfiles := ListAllFilesRecursivelyByFilename(*directory, filename)

	covariances := Select(cfiles, *hands)
	covariances = Select(covariances, *ctrls)

	grasps := ReplaceInAll(covariances, filename, "hand.sofastates.csv")

	mcw := make([]float64, len(grasps), len(grasps))

	var data tsne4go.VectorDistancer
	data = make([][]float64, len(covariances), len(covariances))
	for i, f := range covariances {
		data[i] = ReadCSVToArray(f)
	}

	tsne := tsne4go.New(data, nil)

	for i := 0; i < 5000; i++ {
		tsne.Step()
	}

	for _, g := range grasps {
		graspData := ReadCSVToFloat(g)
	}

	xy := make([][]float64, len(covariances), len(covariances))
	for i := 0; i < len(covariances); i++ {
		xy[i] = make([]float64, 2, 2)
		xy[i][0] = tsne.Solution[i][0]
		xy[i][1] = tsne.Solution[i][1]
	}

	WriteCSVFloat("/Users/zahedi/Desktop/out.csv", xy)

}
