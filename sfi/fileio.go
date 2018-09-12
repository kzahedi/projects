package main

import (
	"io/ioutil"
	"strings"
)

func readFileToList(file string) []string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	return lines
}

func getLoginPassword(file string) (string, string) {
	lines := readFileToList(file)
	return lines[0], lines[1]
}
