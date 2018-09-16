package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func cleanUp() {
	var delete []string

	files := readDirContent("data/*.json")

	for _, f := range files {
		bytes, _ := ioutil.ReadFile(f)
		tweet := makeTweet()
		json.Unmarshal(bytes, &tweet)
		if tweet.Valid() == false {
			delete = append(delete, f)
		}
	}
	for _, v := range delete {
		fmt.Println("Deleting", v)
		os.Remove(v)
	}
}
