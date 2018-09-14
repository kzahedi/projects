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

	// os.Remove("data/1039261599012450305.json")
	// collectReplyTree([]string{"https://twitter.com/MartinaRenner/status/1039261599012450305"})
}
