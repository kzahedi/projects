package main

import (
	"fmt"
	"regexp"

	"github.com/sacado/tsne4go"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func CalculateTSNE(filename string, hand, controller *regexp.Regexp, directory *string, iterations int, successfulOnly bool, results *Results) {
	fmt.Println("Calculating TSNE")
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	covariances := Select(files, *hand)
	covariances = Select(covariances, *controller)

	var selected []string
	if successfulOnly == false {
		selected = covariances
	} else {
		fmt.Println("number:", len(covariances))
		for _, v := range covariances {
			key := GetKey(v)
			elem := (*results)[key]
			if elem.Successful {
				selected = append(selected, v)
			}
		}
		fmt.Println("number:", len(selected))
	}

	var data tsne4go.VectorDistancer
	data = make([][]float64, len(selected), len(selected))
	for i, f := range selected {
		data[i] = ReadCSVToArray(f)
	}

	tsne := tsne4go.New(data, nil)

	bar := pb.StartNew(iterations)

	for i := 0; i < iterations; i++ {
		tsne.Step()
		bar.Increment()
	}
	bar.Finish()

	for i := 0; i < len(selected); i++ {
		key := GetKey(selected[i])
		v := (*results)[key]
		v.Point[0] = tsne.Solution[i][0]
		v.Point[1] = tsne.Solution[i][1]
		v.ClusteredByTSE = true
		(*results)[key] = v
	}
}
