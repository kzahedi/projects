package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Result struct {
	MC_W                float64
	GraspDistance       float64
	PosX                float64
	PosY                float64
	ObjectType          int
	ObjectPosition      int
	ClusteredByTSE      bool
	Successful          bool
	Intelligent         bool
	Stupid              bool
	SelectedForAnalysis bool
	Distance            float64
}

type Results map[string]Result

func PrintResults(r map[string]Result) {
	for key, value := range r {
		fmt.Println(fmt.Sprintf("%s: MC_W: %f, Grasp Distance: %f, Point: (%f,%f), Object Type: %d, Object Position %d, Successful %t, Intelligent %t, Stupid %t, Selected for Analysis %t, Distance %f", key, value.MC_W, value.GraspDistance, value.PosX, value.PosY, value.ObjectType, value.ObjectPosition, value.Successful, value.Intelligent, value.Stupid, value.SelectedForAnalysis, value.Distance))
	}
}

func PrintIntelligent(r map[string]Result) {
	for _, value := range r {
		if value.Intelligent {
			fmt.Println(fmt.Sprintf("MC_W: %f, Grasp Distance: %f,  Intelligent %t, Stupid %t, Selected %t, Distance %f", value.MC_W, value.GraspDistance, value.Intelligent, value.Stupid, value.SelectedForAnalysis, value.Distance))
		}
	}
}

func PrintStupid(r map[string]Result) {
	for _, value := range r {
		if value.Stupid {
			fmt.Println(fmt.Sprintf("MC_W: %f, Grasp Distance: %f,  Intelligent %t, Stupid %t, Selected %t, Distance %f", value.MC_W, value.GraspDistance, value.Intelligent, value.Stupid, value.SelectedForAnalysis, value.Distance))
		}
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
				r := Result{MC_W: 0.0, GraspDistance: 0.0, PosX: 0.0, PosY: 0.0, ObjectType: -1, ObjectPosition: -1, ClusteredByTSE: false, Intelligent: false, Stupid: false, SelectedForAnalysis: false, Distance: -1.0}
				(*results)[key] = r
			}
		}
	}
}

func WriteResults(filename string, results Results, outputDir string) {
	file, err := os.Create(fmt.Sprintf("%s/%s", outputDir, filename))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)
	s := "# Experiment, MC_W, GraspDistance, PosX, PosY, ObjectType, ObjectPosition, Successful, Intelligent, Stupid, SelectedForAnalysis, Distance"

	w.WriteString(s)
	for key, value := range results {
		if value.ClusteredByTSE == true {
			s = fmt.Sprintf("\n%s,%f,%f,%f,%f,%d,%d,%t,%t,%t,%t,%f", key, value.MC_W, value.GraspDistance, value.PosX, value.PosY, value.ObjectType, value.ObjectPosition, value.Successful, value.Intelligent, value.Stupid, value.SelectedForAnalysis, value.Distance)
			w.WriteString(s)
			w.Flush()
		}
	}
}

func ReadResults(filename string) Results {
	results := make(Results)
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range records {
		if strings.Contains(v[0], "#") {
			continue
		}
		key := v[0]
		mcw, _ := strconv.ParseFloat(v[1], 64)
		graspDistance, _ := strconv.ParseFloat(v[2], 64)
		posX, _ := strconv.ParseFloat(v[3], 64)
		posY, _ := strconv.ParseFloat(v[4], 64)
		otype, _ := strconv.ParseInt(v[5], 10, 64)
		opos, _ := strconv.ParseInt(v[6], 10, 64)
		intelligent, _ := strconv.ParseBool(v[7])
		stupid, _ := strconv.ParseBool(v[8])
		selected, _ := strconv.ParseBool(v[9])
		distance, _ := strconv.ParseFloat(v[10], 64)

		results[key] = Result{MC_W: mcw, GraspDistance: graspDistance, PosX: posX, PosY: posY, ObjectType: int(otype), ObjectPosition: int(opos), ClusteredByTSE: true, Successful: selected, Intelligent: intelligent, Stupid: stupid, SelectedForAnalysis: selected, Distance: distance}
	}

	return results
}
