package main

import (
	"fmt"
)

func getNewStartingPoints(lst *[]string, account string) []string {
	var newTweets []string

	service, wd := randomLogin()
	defer service.Stop()
	defer wd.Close()

	startingPointURL := fmt.Sprintf("https://twitter.com/%s", account)
	openURL(startingPointURL, &wd)

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
						if add == true && len(newURL) > len("https://twitter.com") {
							newTweets = append(newTweets, newURL)
						}
					}
				}
			}
		}
	}

	err := wd.Quit()
	if err != nil {
		panic(err)
	}

	return newTweets
}
