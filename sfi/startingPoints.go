package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func getNewStartingPoints(lst []string, accounts []string) []string {
	var newTweets []string

	service, wd := randomLogin()
	defer service.Stop()
	defer wd.Close()

	for _, account := range accounts {
		startingPointURL := fmt.Sprintf("https://twitter.com/%s", account)
		fmt.Printf("Checking %s\n", startingPointURL)
		openURL(startingPointURL, &wd)

		for true {
			for i := 0; i < 10; i++ {
				wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight)", nil)
			}
			tweets := findElementsByCSS("div.tweet", &wd)
			var newlyCollected []string
			for _, t := range tweets {
				s, _ := t.GetAttribute("data-permalink-path")
				if len(s) > 5 {
					link := fmt.Sprintf("https://twitter.com%s", s)
					newlyCollected = append(newlyCollected, link)
				}
			}

			found := false
			for _, newURL := range newlyCollected {
				if contains(newURL, lst) == true {
					found = true
				} else {
					newTweets = append(newTweets, newURL)
				}
			}
			if found == true {
				break
			}
		}
	}

	err := wd.Quit()
	if err != nil {
		panic(err)
	}

	return newTweets
}

func contains(entry string, lst []string) bool {
	for _, v := range lst {
		if v == entry {
			return true
		}
	}
	return false
}

func getAllStartingPoints(lst []string, accounts []string) []string {
	var newTweets []string

	service, wd := randomLogin()
	defer service.Stop()
	defer wd.Close()

	for _, account := range accounts {
		startingPointURL := fmt.Sprintf("https://twitter.com/%s", account)
		fmt.Printf("Checking %s\n", startingPointURL)
		openURL(startingPointURL, &wd)

		m := 0
		for true {
			for i := 0; i < 10; i++ {
				wd.ExecuteScript("window.scrollTo(0, document.body.scrollHeight)", nil)
			}
			tweets := findElementsByCSS("div.tweet", &wd)
			n := len(tweets)
			fmt.Println(m, "=", n)
			if m == n {
				break
			}
			m = n

			for _, t := range tweets {
				s, _ := t.GetAttribute("data-permalink-path")
				if len(s) > 1 {
					s = fmt.Sprintf("https://twitter.com%s", s)
					fmt.Println("checking", s)
					if contains(s, lst) == false {
						newTweets = append(newTweets, s)
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

func cleanUpStartingPoints() {
	startingPoints := readFileToList("watch/starting_points.txt")
	var r []string
	for _, s := range startingPoints {
		if s == "https://twitter.com" {
			continue
		}
		if contains(s, r) == false {
			r = append(r, s)
		}
	}
	writeListToFile(&r, "watch/starting_points.txt")
}

func collectNewStartingPoints(cpus int, all bool) {
	cleanUpStartingPoints()
	startingPoints := readFileToList("watch/starting_points.txt")
	accounts := readFileToList("watch/accounts.txt")

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
				if all == true {
					ans <- getAllStartingPoints(startingPoints, p)
				} else {
					ans <- getNewStartingPoints(startingPoints, p)
				}
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
			appendToFile("watch/starting_points.txt", v)
		}
	}

	// newStartingPoints := getNewStartingPoints(&startingPoints, account)

	// startingPoints = append(startingPoints, newStartingPoints...)
	// writeListToFile(&startingPoints, "data/starting_points.txt")
}
