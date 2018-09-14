package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func getNewStartingPoints(lst *[]string, accounts []string) []string {
	var newTweets []string

	service, wd := randomLogin()
	defer service.Stop()
	defer wd.Close()

	for _, account := range accounts {
		startingPointURL := fmt.Sprintf("https://twitter.com/%s", account)
		fmt.Println(startingPointURL)
		openURL(startingPointURL, &wd)

		continueToLook := true
		for continueToLook == true {
			wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight)", nil)
			var newlyCollected []string
			tweets := findElementsByCSS("div.tweet", &wd)
			for _, t := range tweets {
				s, _ := t.GetAttribute("data-permalink-path")
				newlyCollected = append(newlyCollected, fmt.Sprintf("https://twitter.com%s", s))
			}

			for _, newURL := range newlyCollected {
				newFound := false
				for _, oldURL := range *lst {
					if newURL == oldURL {
						newFound = true
						break
					}
				}
				if newFound == false && len(newURL) > 20 {
					newTweets = append(newTweets, newURL)
				} else {
					continueToLook = false
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

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(accounts), func(i, j int) { accounts[i], accounts[j] = accounts[j], accounts[i] })

	// start workers
	var wg sync.WaitGroup
	send := make(chan []string, cpus*2)
	ans := make(chan []string, cpus*2)

	for i := 0; i < cpus; i++ {
		wg.Add(1)
		go func(send <-chan []string, ans chan<- []string) {
			defer wg.Done()
			for p := range send {
				ans <- getNewStartingPoints(&startingPoints, p)
			}
		}(send, ans)
	}

	// start the jobs
	go func(send chan<- []string) {
		for i := 0; i < len(accounts)-5; i += 5 {
			send <- accounts[i : i+5]
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

	// startingPoints = append(startingPoints, newStartingPoints...)
	// writeListToFile(&startingPoints, "data/starting_points.txt")
}
