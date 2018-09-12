package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

func writeListToFile(list *[]string, filename string) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	for _, s := range *list {
		f.WriteString(fmt.Sprintf("%s\n", s))
	}
}

func getLoginPassword(file string) (string, string) {
	lines := readFileToList(file)
	return lines[0], lines[1]
}
