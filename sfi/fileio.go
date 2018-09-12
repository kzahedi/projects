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

func appendToFile(filename, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func getLoginPassword(file string) (string, string) {
	lines := readFileToList(file)
	return lines[0], lines[1]
}

func writeBytesToFile(filename string, bytes []byte) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	f.Write(bytes)
}
