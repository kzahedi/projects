package main

import (
	"github.com/kzahedi/projects/sfi/twitter"
)

func getTwitterHandles(tweet twitter.Tweet) []string {
	var r []string

	r = append(r, tweet.TwitterHandle)

	for _, c := range tweet.Children {
		s := getTwitterHandles(c)
		r = append(r, s...)
	}
	return r
}

func countHandles(tweet twitter.Tweet) map[string]int {
	r := make(map[string]int)
	handles := getTwitterHandles(tweet)

	for _, h := range handles {
		r[h] = r[h] + 1
	}

	return r
}
