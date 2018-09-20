package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	cpus := flag.Int("cpu", 2, "CPUS")
	cleanUp := flag.Bool("cu", false, "Clean up starting_points.txt")
	getNew := flag.Bool("gn", false, "Get new tweets")
	all := flag.Bool("all", false, "Get all tweets")
	if *cleanUp == true {
		fmt.Println("Cleaning up")
		cleanUpStartingPoints()
	}

	flag.Parse()
	rand.Seed(time.Now().Unix())
	for true {
		if *getNew == true {
			collectNewStartingPoints(*cpus, *all)
		}
		collectReplyTrees(*cpus)
		time.Sleep(10 * time.Hour)
	}

	// collectReplyTree([]string{"https://twitter.com/GregorGysi/status/818779589644259328"})
}
