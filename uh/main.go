package main

import (
	"flag"
)

func main() {
	tPtr := flag.Float64("t", 1.0, "Twistedness")
	nPtr := flag.Int("n", 11, "Size of the world")

	createWorld(*nPtr, *tPtr)
}
