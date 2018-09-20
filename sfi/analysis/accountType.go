package main

import (
	"strings"

	"github.com/kzahedi/projects/sfi/twitter"
)

func checkForHateAccounts(tweet *twitter.Tweet, hateAccount string) bool {
	if (*tweet).TwitterHandle == hateAccount {
		// fmt.Printf("Found Hate %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		(*tweet).HateAccount = true
		return true
	}
	return false
}

func checkTweet(tweet *twitter.Tweet, hate *[]string, counter *[]string) bool {
	found := false

	for _, h := range *hate {
		h = strings.Replace(h, "@", "", 1)
		if tweet.TwitterHandle == h {
			// fmt.Printf("Found Hate \"%s\" \"%s\"\n", tweet.TwitterHandle, h)
			(*tweet).HateAccount = true
			found = true
		}
	}

	for _, c := range *counter {
		c = strings.Replace(c, "@", "", 1)
		if (*tweet).TwitterHandle == c {
			// fmt.Printf("Found Counter \"%s\" \"%s\"\n", tweet.TwitterHandle, c)
			(*tweet).CounterAccount = true
			found = true
		}
	}

	if strings.Contains((*tweet).Name, "❌") == true {
		// fmt.Printf("Found Hate %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.HateAccount = true
		found = true
	}

	if strings.Contains((*tweet).Name, "QFD") == true {
		// fmt.Printf("Found Hate %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		(*tweet).HateAccount = true
		found = true
	}

	if strings.Contains((*tweet).Name, "Shadowban") == true {
		// fmt.Printf("Found Hate %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		(*tweet).HateAccount = true
		found = true
	}

	if strings.Contains((*tweet).Name, "⭕️") == true {
		// fmt.Printf("Found Counter %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.CounterAccount = true
		found = true
	}

	if strings.Contains((*tweet).Name, "2MInt") == true {
		// fmt.Printf("Found Counter %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.CounterAccount = true
		found = true
	}

	if strings.Contains((*tweet).Name, "#FBPE") == true {
		// fmt.Printf("Found Counter %s \"%s\"\n", tweet.TwitterHandle, tweet.Name)
		tweet.CounterAccount = true
		found = true
	}

	var c []twitter.Tweet
	for _, t := range (*tweet).Children {
		if checkTweet(&t, hate, counter) == true {
			found = true
		}
		c = append(c, t)
	}
	(*tweet).Children = c
	return found
}
