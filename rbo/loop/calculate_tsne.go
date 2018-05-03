package main

import (
	"fmt"
	"regexp"

	"github.com/sacado/tsne4go"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func CalculateTSNE(input string, hands, ctrls []*regexp.Regexp, directory *string, iterations int, successfulOnly bool, results Results) Results {
	fmt.Println("Calculating TSNE")
	files := ListAllFilesRecursivelyByFilename(*directory, input)

	covariances := SelectFiles(files, hands, ctrls)

	var selected []string
	if successfulOnly == false {
		selected = covariances
	} else {
		for _, v := range covariances {
			key := GetKey(v)
			elem := results[key]
			if elem.Successful {
				// fmt.Println("found selected")
				selected = append(selected, v)
			}
		}
		fmt.Println("number:", len(selected))
	}

	fmt.Println("Clustering on:", len(covariances))

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
		v := results[key]
		v.PosX = tsne.Solution[i][0]
		v.PosY = tsne.Solution[i][1]
		v.ClusteredByTSE = true
		results[key] = v
	}

	return results
}
