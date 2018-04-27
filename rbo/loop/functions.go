package main

import (
	"regexp"
	"strings"
)

func Select(lst []string, pattern regexp.Regexp) []string {
	var r []string

	for _, f := range lst {
		if pattern.MatchString(f) == true {
			r = append(r, f)
		}
	}

	return r
}

func Exclude(lst []string, pattern regexp.Regexp) []string {
	var r []string

	for _, f := range lst {
		if pattern.MatchString(f) == false {
			r = append(r, f)
		}
	}

	return r
}

func ReplaceInAll(lst []string, a, b string) []string {
	var r []string

	for _, f := range lst {
		r = append(r, strings.Replace(f, a, b, -1))
	}

	return r
}

func GetKey(s string) string {
	re := regexp.MustCompile("rbo[a-zA-Z0-9-]+/[a-zA-Z0-9_.-]+")
	return re.FindAllString(s, -1)[0]
}

func GetObjectName(s string) string {
	re := regexp.MustCompile("object[a-zA-Z0-9-]+")
	return re.FindAllString(s, -1)[0]
}

func getColumn(data [][]float64, col int) []float64 {
	r := make([]float64, len(data), len(data))
	for row := 0; row < len(data); row++ {
		r[row] = data[row][col]
	}
	return r
}

func SelectFiles(files []string, hands, ctrls []*regexp.Regexp) []string {
	var r []string
	for _, hand := range hands {
		for _, ctrl := range ctrls {
			f := Select(files, *hand)
			f = Select(f, *ctrl)
			r = append(r, f...)
		}
	}
	return r
}
