package main

import (
	"flag"
	"math/rand"
	"time"
)

func main() {
	cpus := flag.Int("cpu", 2, "CPUS")
	flag.Parse()
	rand.Seed(time.Now().Unix())
	for true {
		// collectNewStartingPoints(*cpus)
		collectReplyTrees(*cpus)
		time.Sleep(10 * time.Hour)
	}
}
