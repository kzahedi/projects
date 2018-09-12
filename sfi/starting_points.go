package main

import (
	"fmt"
	"time"
)

func getNewStartingPoints(lst *[]string, account string) []string {
	var r []string
	service, wd := randomLogin()
	defer service.Stop()
	defer wd.Close()

	newURL := fmt.Sprintf("https://twitter.com/%s", account)
	openURL(newURL, &wd)

	var newTweets []string

	found := false
	for found == false {
		wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight)", nil)
		var r []string
		tweets := findElementsByCSS("div.tweet", &wd)
		for _, t := range tweets {
			s, _ := t.GetAttribute("data-permalink-path")
			r = append(r, fmt.Sprintf("https://twitter.com%s", s))
		}
		for _, newURL := range r {
			for _, oldURL := range *lst {
				fmt.Printf("Comparing %s with %s\n", newURL, oldURL)
				if newURL == oldURL {
					found = true
				} else {
					if len(newTweets) == 0 {
						newTweets = append(newTweets, newURL)
					} else {
						add := true
						for _, url := range newTweets {
							if newURL == url {
								add = false
								break
							}
						}
						if add == true {
							newTweets = append(newTweets, newURL)
						}
					}
				}
			}
		}
	}

	fmt.Println(newTweets)

	time.Sleep(30 * time.Second)

	return r
}
