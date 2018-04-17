package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
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

func ReadResults(filename string) Results {
	data := make(Results)

	file, _ := os.Open(filename)
	defer file.Close()

	reader := csv.NewReader(file)
	record, err := reader.Read()
	for err != io.EOF {

		if strings.HasPrefix(record[0], "#") {
			record, err = reader.Read()
			continue
		}

		// "# Experiment, MC_W, Grasp Distance, t-SNE X, t-SNE Y, Object Type, Object Position, Successful"
		experiment := record[0]
		mcw, _ := strconv.ParseFloat(record[1], 64)
		graspDistance, _ := strconv.ParseFloat(record[2], 64)
		posX, _ := strconv.ParseFloat(record[3], 64)
		posY, _ := strconv.ParseFloat(record[4], 64)
		objectType, _ := strconv.ParseInt(record[5], 10, 64)
		objectPosition, _ := strconv.ParseInt(record[6], 10, 64)

		successfull := false
		if record[7] == "true" {
			successfull = true
		}

		intelligent := false
		if len(record) > 7 {
			if record[7] == "true" {
				intelligent = true
			}
		}

		stupid := false
		if len(record) > 8 {
			if record[8] == "true" {
				stupid = true
			}
		}

		selectedForAnalysis := false
		if len(record) > 9 {
			if record[9] == "true" {
				selectedForAnalysis = true
			}
		}

		distance := -1.0
		if len(record) > 10 {
			distance, _ = strconv.ParseFloat(record[10], 64)
		}

		r := Result{MC_W: mcw, GraspDistance: graspDistance, PosX: posX, PosY: posY, ObjectType: int(objectType), ObjectPosition: int(objectPosition), ClusteredByTSE: true, Successful: successfull, Intelligent: intelligent, Stupid: stupid, SelectedForAnalysis: selectedForAnalysis, Distance: distance}

		data[experiment] = r

		record, err = reader.Read()
	}

	return data
}

func WriteResults(filename string, results Results) {
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)
	s := "# Experiment, MC_W, Grasp Distance, t-SNE X, t-SNE Y, Object Type, Object Position, Successful, Distance"
	w.WriteString(s)
	for key, value := range results {
		if value.ClusteredByTSE == true {
			s = fmt.Sprintf("\n%s,%f,%f,%f,%f,%d,%d,%t,%t,%t,%t,%f", key, value.MC_W, value.GraspDistance, value.PosX, value.PosY, value.ObjectType, value.ObjectPosition, value.Successful, value.Intelligent, value.Stupid, value.SelectedForAnalysis, value.Distance)
			w.WriteString(s)
			w.Flush()
		}
	}
}
