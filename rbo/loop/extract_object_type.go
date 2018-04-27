package main

import (
	"fmt"
	"regexp"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ExtractObjectType(results Results) Results {
	fmt.Println("Extract Object Type")
	bar := pb.StartNew(len(results))

	r := make(map[string]int)

	index := 0

	for key, _ := range results {
		s := extractObjectString(key)
		if _, ok := r[s]; ok == false {
			r[s] = index
			index++
		}
		bar.Increment()
	}
	bar.Finish()

	for key, value := range results {
		s := extractObjectString(key)
		value.ObjectType = r[s]
		results[key] = value
	}

	return results
}

func extractObjectString(in string) string {
	re := regexp.MustCompile("object[a-zA-Z]+")
	return re.FindAllString(in, -1)[0]
}
