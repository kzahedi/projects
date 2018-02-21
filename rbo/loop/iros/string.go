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
