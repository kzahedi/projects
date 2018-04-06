package main

import (
	"fmt"
	"regexp"
)

type Result struct {
	MC_W           float64
	GraspDistance  float64
	PosX           float64
	PosY           float64
	ObjectType     int
	ObjectPosition int
	ClusteredByTSE bool
	Successful     bool
	Intelligent    bool
	Stupid         bool
}

type Results map[string]Result

func PrintResults(r map[string]Result) {
	for key, value := range r {
		fmt.Println(fmt.Sprintf("%s: MC_W: %f, Grasp Distance: %f, Point: (%f,%f), Object Type: %d, Object Position %d, Successful %t, Intelligent %t, Stupid %t", key, value.MC_W, value.GraspDistance, value.PosX, value.PosY, value.ObjectType, value.ObjectPosition, value.Successful, value.Intelligent, value.Stupid))
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
				r := Result{MC_W: 0.0, GraspDistance: 0.0, PosX: 0.0, PosY: 0.0, ObjectType: -1, ObjectPosition: -1, ClusteredByTSE: false, Intelligent: false, Stupid: false}
				(*results)[key] = r
			}
		}
	}

}
