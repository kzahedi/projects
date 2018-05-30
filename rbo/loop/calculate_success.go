package main

import (
	"fmt"
	"math"
	"regexp"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func CalculateSuccess(input string, hands, ctrls []*regexp.Regexp, directory string, height float64, results Results) Results {
	fmt.Println("Calculating Success")
	objectFiles := ListAllFilesRecursivelyByFilename(directory, input)

	osizes := make(map[string]float64)
	osizes["objectcylinder"] = 20.0
	osizes["objectcylinderB"] = 40.0
	osizes["objectbox"] = 35.0
	osizes["objectboxB"] = 20.0
	osizes["objectsphere"] = 35.0
	osizes["objectsphereB"] = 20.0
	osizes["objectegg"] = 35.0
	osizes["objecteggB"] = 20.0

	selectedFiles := SelectFiles(objectFiles, hands, ctrls)
	bar := pb.StartNew(len(selectedFiles))

	for _, s := range selectedFiles {
		data := ReadCSVToFloat(s)
		maxHeight := data[20][1]
		for i := 20; i < len(data); i++ {
			maxHeight = math.Max(maxHeight, data[i][1])
		}

		key := GetKey(s)
		objectName := GetObjectName(s)

		v := results[key]
		v.Successful = ((maxHeight - osizes[objectName]) > height)
		results[key] = v

		bar.Increment()
	}
	bar.Finish()

	return results
}
