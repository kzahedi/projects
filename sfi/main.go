package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	cpus := flag.Int("cpu", 2, "CPUS")
	flag.Parse()

	rand.Seed(time.Now().Unix())

	startingPoints := readFileToList("data/starting_points.txt")
	accounts := readFileToList("data/accounts.txt")

	// start workers
	var wg sync.WaitGroup
	send := make(chan string, *cpus*2)
	ans := make(chan []string, *cpus*2)

	for i := 0; i < *cpus; i++ {
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
		}
	}

	// newStartingPoints := getNewStartingPoints(&startingPoints, account)

	startingPoints = append(startingPoints, newStartingPoints...)
	writeListToFile(&startingPoints, "data/starting_points.txt")

	// exec.Command("killall", "-9", "firefox")
	// exec.Command("killall", "-9", "geckodriver")

}
