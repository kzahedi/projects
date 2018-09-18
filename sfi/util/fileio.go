package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ReadFileToList return list of string from file
func ReadFileToList(file string) []string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	return lines
}

// WriteListToFile writes a string list to file
func WriteListToFile(filename string, list *[]string) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	f.WriteString(fmt.Sprintf("%s", (*list)[0]))
	for i := 1; i < len(*list); i++ {
		f.WriteString(fmt.Sprintf("\n%s", (*list)[i]))
	}
}

// AppendToFile appends a single line to a file
func AppendToFile(filename, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			f, _ = os.Create(filename)
		} else {
			panic(err)
		}
	}

	defer f.Close()

	line := fmt.Sprintf("\n%s", text)

	if _, err = f.WriteString(line); err != nil {
		panic(err)
	}
}

// WriteBytesToFile write bytes to a file
func WriteBytesToFile(filename string, bytes []byte) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	f.Write(bytes)
}

// ReadDirContent returns string list from dir
func ReadDirContent(pattern string) []string {
	files, err := filepath.Glob(pattern)
	if err != nil {
		panic(err)
	}
	return files
}
