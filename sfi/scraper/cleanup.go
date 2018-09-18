package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kzahedi/projects/sfi/io"
	"github.com/kzahedi/projects/sfi/twitter"
)

func cleanUp() {
	var delete []string

	files := io.ReadDirContent("data/*.json")

	for _, f := range files {
		bytes, _ := ioutil.ReadFile(f)
		tweet := twitter.MakeTweet()
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
