package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/kzahedi/projects/sfi/twitter"
	"github.com/kzahedi/projects/sfi/util"
	pb "gopkg.in/cheggaaa/pb.v1"
)

func main() {
	dir := flag.String("d", "", "Input dir")
	hate := flag.String("h", "", "hate accounts")
	counter := flag.String("c", "", "counter accounts")
	min := flag.Int("min", 10, "Minimum number of interactions")
	flag.Parse()

	files := readDirContent(fmt.Sprintf("%s/*.json", *dir))

	hateAccounts := util.ReadFileToList(*hate)
	counterAccounts := util.ReadFileToList(*counter)

	fmt.Println("Marking Hate and Counter Accounts")
	bar := pb.StartNew(len(files))
	for _, f := range files {
		// fmt.Println(f)
		tweet := twitter.ReadTweetJSON(f)
		checkTweet(&tweet, &hateAccounts, &counterAccounts)
		tweet.ExportJSON(f)
		bar.Increment()
	}
	bar.Finish()

	fmt.Println("Twitter handle histogram")
	bar = pb.StartNew(len(files))
	counts := make(map[string]int)
	for _, f := range files {
		// fmt.Println(f)
		tweet := twitter.ReadTweetJSON(f)
		r := countHandles(tweet)
		for k, v := range r {
			counts[k] = counts[k] + v
		}
		bar.Increment()
	}
	bar.Finish()

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range counts {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	watchedAccounts := util.ReadFileToList("../scraper/watch/accounts.txt")

	var list []string

	for _, kv := range ss {
		if kv.Value >= *min &&
			util.ListContains(&watchedAccounts, kv.Key) == false &&
			util.ListContains(&hateAccounts, kv.Key) == false {
			list = append(list, fmt.Sprintf("%s,%d", kv.Key, kv.Value))
		}
	}
	util.WriteListToFile("../scraper/watch/account_counts.csv", &list)
}
