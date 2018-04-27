package main

import (
	"fmt"
	"regexp"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ExtractObjectPosition(results Results) Results {
	fmt.Println("Extract Object Position")
	bar := pb.StartNew(len(results))

	r := make(map[string]int)

	index := 0

	for key, _ := range results {
		s := extractPositionString(key)
		if _, ok := r[s]; ok == false {
			r[s] = index
			index++
		}
		bar.Increment()
	}
	bar.Finish()

	for key, value := range results {
		s := extractPositionString(key)
		value.ObjectPosition = r[s]
		results[key] = value
	}

	return results
}

func extractPositionString(in string) string {
	re := regexp.MustCompile("-?[0-9]{1,2}.[0-9]{0,2}_-?[0-9]{1,2}.[0-9]{0,2}_-?[0-9]{1,2}.[0-9]{0,2}")
	return re.FindAllString(in, -1)[0]
}
