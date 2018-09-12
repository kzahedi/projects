package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	cpus := flag.Int("cpu", 2, "CPUS")
	flag.Parse()
	rand.Seed(time.Now().Unix())

	if *cpus != 2 {
		fmt.Println("woops")
	}
	// collectNewStartingPoints(*cpus)

	// collectReplyTree("https://twitter.com/ArasBacho/status/1031661358444630020")
	collectReplyTree("https://twitter.com/tagesschau/status/1039884797999558656")
}
