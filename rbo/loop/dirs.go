package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func ListAllSubdirectories(parent string) []string {
	var r []string

	files, err := ioutil.ReadDir(parent)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			r = append(r, file.Name())
		}
	}
	return r
}

func ListAllFiles(parent string) []string {
	var r []string

	files, err := ioutil.ReadDir(parent)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() == false {
			r = append(r, file.Name())
		}
	}
	return r
}

func ListAllFilesRecursivelyByExtension(parent, extension string) []string {
	var r []string

	subdirs := ListAllSubdirectories(parent)
	files := ListAllFiles(parent)

	for _, file := range files {
		if filepath.Ext(file) == extension {
			r = append(r, fmt.Sprintf("%s/%s", parent, file))
		}
	}

	for _, dir := range subdirs {
		s := ListAllFilesRecursivelyByExtension(fmt.Sprintf("%s/%s", parent, dir), extension)
		r = append(r, s...)
	}
	return r
}

func ListAllFilesRecursivelyByFilename(parent, filename string) []string {
	var r []string

	subdirs := ListAllSubdirectories(parent)
	files := ListAllFiles(parent)

	for _, file := range files {
		if file == filename {
			r = append(r, fmt.Sprintf("%s/%s", parent, file))
		}
	}

	for _, dir := range subdirs {
		s := ListAllFilesRecursivelyByFilename(fmt.Sprintf("%s/%s", parent, dir), filename)
		r = append(r, s...)
	}
	return r
}

func CreateDir(dir string) {
	os.MkdirAll(path.Dir(dir), os.ModePerm)
}
