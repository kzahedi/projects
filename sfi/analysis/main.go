package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kzahedi/projects/sfi/tweet"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func checkForHateAccounts(tweet *tweet.Tweet, hate string) bool {
	if tweet.TwitterHandle == hate {
		// fmt.Printf("Found Hate %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.HateAccount = true
		return true
	}

	if strings.Contains(tweet.Name, "❌") == true {
		// fmt.Printf("Found Hate %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.HateAccount = true
		return true
	}

	if strings.Contains(tweet.Name, "QFD") == true {
		// fmt.Printf("Found Hate %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.HateAccount = true
		return true
	}

	if strings.Contains(tweet.Name, "⭕️") == true {
		// fmt.Printf("Found Counter %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.HateAccount = true
		return true
	}

	if strings.Contains(tweet.Name, "2MInt") == true {
		// fmt.Printf("Found Counter %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.HateAccount = true
		return true
	}

	if strings.Contains(tweet.Name, "#FBPE") == true {
		// fmt.Printf("Found Counter %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.HateAccount = true
		return true
	}

	found := false
	for _, t := range (*tweet).Children {
		f := checkForHateAccounts(&t, hate)
		found = found || f
	}
	return found
}

func main() {
	dir := flag.String("d", "", "Input dir")
	hate := flag.String("h", "", "hate accounts")
	flag.Parse()

	files := readDirContent(fmt.Sprintf("%s/*.json", *dir))
	hateAccounts := readFileToList(*hate)

	bar := pb.StartNew(len(files))

	for _, f := range files {
		tweet := tweet.ReadTweetJSON(f)
		found := false
		for _, h := range hateAccounts {
			h := strings.Replace(h, "@", "", 1)
			found = checkForHateAccounts(&tweet, h)
		}
		bar.Increment()
		if found == true {
			tweet.ExportJSON(f)
		}
	}
	bar.Finish()
}
