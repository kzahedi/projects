package main

import (
	"fmt"
	"regexp"
)

type Result struct {
	MC_W           float64
	GraspDistance  float64
	Point          []float64
	ObjectType     int
	ObjectPosition int
	ClusteredByTSE bool
	Successful     bool
}

type Results map[string]Result

func PrintResults(r map[string]Result) {
	for key, value := range r {
		fmt.Println(fmt.Sprintf("%s: MC_W: %f, Grasp Distance: %f, Point: (%f,%f), Object Type: %d, Object Position %d", key, value.MC_W, value.GraspDistance, value.Point[0], value.Point[1], value.ObjectType, value.ObjectPosition))
	}
}

func CreateResultsContainer(hands, ctrls []*regexp.Regexp, directory *string, results *Results) {
	filename := "hand.sofastates.csv"
	files := ListAllFilesRecursivelyByFilename(*directory, filename)

	for _, hand := range hands {
		for _, ctrl := range ctrls {
			hfiles := Select(files, *hand)
			hfiles = Select(hfiles, *ctrl)
			for _, s := range hfiles {
				key := GetKey(s)
				r := Result{MC_W: 0.0, GraspDistance: 0.0, Point: []float64{0.0, 0.0}, ObjectType: -1, ObjectPosition: -1, ClusteredByTSE: false}
				(*results)[key] = r
			}
		}
	}

}
