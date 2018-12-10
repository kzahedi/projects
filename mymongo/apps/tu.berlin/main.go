package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	mymongo "github.com/kzahedi/projects/mymongo"
)

func appendIfMissing(slice []string, str string) []string {
	for _, ele := range slice {
		if ele == str {
			return slice
		}
	}
	return append(slice, str)
}

func main() {
	client := mymongo.Connect()

	fmt.Println("connected")
	fmt.Println(client)

	// db := client.Database("RBOHand").Collection("IROS Data")

	// collection := client.Database("testing").Collection("numbers")

	var files []string
	err := filepath.Walk("/Users/zahedi/projects/TU.Berlin/experiments/run2017011101",
		func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.Contains(file, "/raw/") {
				p := path.Dir(file)
				files = appendIfMissing(files, p)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	for _, v := range files {
		fmt.Println(split(v))
	}

}

type experiment struct {
	Hand       string
	Controller int
	Object     string
	X          float64
	Y          float64
	Theta      float64
}

// /Users/zahedi/projects/TU.Berlin/experiments/run2017011101/rbohandkz2/hand0_0-controller0-objectsphereB_-30.0_0.0_-45.0

func split(str string) experiment {

	handRE := regexp.MustCompile("/(rbo[a-z0-9-]*)/")
	ctrlRE := regexp.MustCompile("-controller([0-9])-")
	objRE := regexp.MustCompile("-object([a-zA-Z]+)_")
	// coordRE := regexp.MustCompile("_(-*[0-9]+.[0-9]*_-*[0-9]+.[0-9]*_-*[0-9]+.[0-9]*)")

	r := handRE.FindAllStringSubmatch(str, -1)[0][1]
	c := ctrlRE.FindAllStringSubmatch(str, 1)[0][1]
	o := objRE.FindAllStringSubmatch(str, 1)[0][1]
	cc := coordRE.FindAllStringSubmatch(str, 1)
	v := strings.Split(cc[0][1], "_")

	x, _ := strconv.ParseFloat(v[0], 64)
	y, _ := strconv.ParseFloat(v[1], 64)
	theta, _ := strconv.ParseFloat(v[2], 64)

	var x, y, theta float64

	ci, _ := strconv.ParseInt(c, 10, 64)

	return experiment{Hand: r, Controller: int(ci), Object: o, X: x, Y: y, Theta: theta}
}
