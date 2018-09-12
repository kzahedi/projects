package main

import (
	"fmt"
	"sync"
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

func collectNewStartingPoints(cpus int) {
	startingPoints := readFileToList("data/starting_points.txt")
	accounts := readFileToList("data/accounts.txt")

	// start workers
	var wg sync.WaitGroup
	send := make(chan string, cpus*2)
	ans := make(chan []string, cpus*2)

	for i := 0; i < cpus; i++ {
		wg.Add(1)
		go func(send <-chan string, ans chan<- []string) {
			defer wg.Done()
			for p := range send {
				ans <- getNewStartingPoints(&startingPoints, p)
			}
		}(send, ans)
	}

	// start the jobs
	go func(send chan<- string) {
		for _, account := range accounts {
			send <- account
		}
		close(send)
		wg.Wait()
		close(ans)
	}(send)

	var newStartingPoints []string
	for a := range ans {
		for _, v := range a {
			fmt.Printf("Found new starting point %s\n", v)
			newStartingPoints = append(newStartingPoints, v)
			appendToFile("data/starting_points.txt", v)
		}
	}

	// newStartingPoints := getNewStartingPoints(&startingPoints, account)

	startingPoints = append(startingPoints, newStartingPoints...)
	writeListToFile(&startingPoints, "data/starting_points.txt")
}
