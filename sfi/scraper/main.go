package main

import (
	"flag"
	"math/rand"
	"time"
)

func main() {
	cpus := flag.Int("cpu", 2, "CPUS")
	//all := flag.Bool("all", false, "Get all tweets")
	// cleanUpStartingPoints()

	flag.Parse()
	rand.Seed(time.Now().Unix())
	for true {
		//		collectNewStartingPoints(*cpus, *all)
		collectReplyTrees(*cpus)
		time.Sleep(10 * time.Hour)
	}

	// collectReplyTree([]string{"https://twitter.com/GregorGysi/status/818779589644259328"})
}
