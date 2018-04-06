package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/kzahedi/utils"
	pb "gopkg.in/cheggaaa/pb.v1"
)

type Experiment struct {
	Name string
	Avg  float64
	Best float64
}

type XmlIndividual struct {
	XMLName xml.Name `xml:"individual"`
	Id      int      `xml:"id,attr"`
	Fitness float64  `xml:"fitness,attr"`
}

type XmlPopulation struct {
	XMLName     xml.Name        `xml:"population"`
	Generation  int             `xml:"generation,attr"`
	Individuals []XmlIndividual `xml:"individual"`
}

type XmlNMODE struct {
	XMLName    xml.Name      `xml:"nmode"`
	Population XmlPopulation `xml:"population"`
}

func ReadNmode(filename string) (*[]XmlIndividual, error) {
	var xmlNmode XmlNMODE

	filePath, err := filepath.Abs(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	if err := xml.NewDecoder(file).Decode(&xmlNmode); err != nil {
		return nil, err
	}

	return &(xmlNmode.Population.Individuals), nil
}

func main() {
	// directory := flag.String("d", "", "Directory")
	directory := flag.String("d", "/Volumes/LaCie/nmode.paper/nmode.w3irdo", "Directory")
	verbose := flag.Bool("v", false, "Verbose")
	// generation := flag.Int("g", 99, "Generation")
	flag.Parse()

	if *directory == "" {
		fmt.Println("Please provide a directory to analyse.")
		os.Exit(0)
	}

	pattern := regexp.MustCompile(`generation-[0-9]+.xml`)
	r := utils.ListAllFilesRecursivelyByRegularExpression(*directory, pattern, *verbose)

	var bar *pb.ProgressBar
	if *verbose == true {
		fmt.Println(">> Collecting Data")
		bar = pb.StartNew(len(r))
	}

	var experiments []Experiment
	for _, s := range r {
		t, _ := ReadNmode(s)
		best, avg := fitnessValues(t)
		e := Experiment{Name: s, Avg: avg, Best: best}
		experiments = append(experiments, e)
		if *verbose == true {
			bar.Increment()
		}
	}
	if *verbose == true {
		bar.FinishPrint("Finished")

		fmt.Println(">> Sorting")
	}

	sort.Slice(experiments, func(i, j int) bool {
		return experiments[i].Best > experiments[j].Best
	})

	if *verbose == true {
		fmt.Println(">> Done")
	}

	fmt.Println(len(experiments))

	for _, e := range experiments[0:9] {
		fmt.Println(fmt.Sprintf("%s, %f, %f", e.Name, e.Best, e.Avg))
	}

}

func fitnessValues(individuals *[]XmlIndividual) (float64, float64) {
	max := 0.0
	sum := 0.0
	n := float64(len(*individuals))

	for i, v := range *individuals {
		if i == 0 {
			max = v.Fitness
		}

		if max < v.Fitness {
			max = v.Fitness
		}
		sum += v.Fitness
	}

	return max, (sum / n)

}
