package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	startingPoints := readFileToList("data/starting_points.txt")
	accounts := readFileToList("data/accounts.txt")
	account := accounts[0]
	newStartingPoints := getNewStartingPoints(&startingPoints, account)

	fmt.Println(newStartingPoints)

	time.Sleep(30 * time.Second)

}
