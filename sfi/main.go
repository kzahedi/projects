package main

import (
	"flag"
	"math/rand"
	"os/exec"
	"time"
)

func main() {
	cpus := flag.Int("cpu", 2, "CPUS")
	flag.Parse()
	rand.Seed(time.Now().Unix())

	for true {
		// exec.Command("killall -9 firefox")
		// exec.Command("killall -9 geckodriver")
		// collectNewStartingPoints(*cpus)
		exec.Command("killall -9 firefox")
		exec.Command("killall -9 java")
		exec.Command("killall -9 geckodriver")
		collectReplyTrees(*cpus)
		time.Sleep(10 * time.Hour)
	}
}
