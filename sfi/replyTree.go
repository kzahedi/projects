package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tebeka/selenium"
)

var id int

// Tweet contains the tree
type Tweet struct {
	ID            int
	ParentID      int
	Name          string
	TwitterHandle string
	Text          string
	Date          string
	Replies       int
	Retweets      int
	Likes         int
	Link          string
	Type          string
	Mentions      []string
	Children      []Tweet
}

func makeTweet() Tweet {
	return Tweet{ID: -1,
		Name:          "",
		TwitterHandle: "",
		Type:          "",
		Text:          "",
		Date:          "",
		Link:          "",
		Replies:       0,
		Retweets:      0,
		Likes:         0,
		Children:      nil,
		ParentID:      -1,
		Mentions:      make([]string, 0, 0),
	}
}

func (t Tweet) String() string {
	s := ""
	s = fmt.Sprintf("%sName: %s\n", s, t.Name)
	s = fmt.Sprintf("%sTwitter Handle: %s\n", s, t.TwitterHandle)
	s = fmt.Sprintf("%sText: \"%s\"\n", s, t.Text)
	s = fmt.Sprintf("%sReplies: %d\n", s, t.Replies)
	s = fmt.Sprintf("%sLikes: %d\n", s, t.Likes)
	s = fmt.Sprintf("%sRetweets: %d\n", s, t.Retweets)
	s = fmt.Sprintf("%sParent: %d\n", s, t.ParentID)
	s = fmt.Sprintf("%sID: %d\n", s, t.ID)
	s = fmt.Sprintf("%sMentions:\n", s)
	for _, v := range t.Mentions {
		s = fmt.Sprintf("%s  %s\n", s, v)
	}
	for _, t := range t.Children {
		s = fmt.Sprintf("%s%s", s, t)
	}
	return s
}

func countThreads(wd *selenium.WebDriver) int {
	return len(findElementsByCSS("li.ThreadedConversation", wd)) +
		len(findElementsByCSS("li.ThreadedConversation--loneTweet", wd))
}

func scrollAllTheWay(wd *selenium.WebDriver) {
	n := countThreads(wd)
	m := -1

	for n != m {
		n = countThreads(wd)
		for i := 0; i < 30; i++ {
			(*wd).ExecuteScript("var objDiv = document.getElementById(\"permalink-overlay\"); objDiv.scrollTop = objDiv.scrollHeight;", nil)
			more := findElementByCSS("button.ThreadedConversation-showMoreThreadsButton.u-textUserColor", wd)
			if more != nil {
				more.Click()
			}
		}
		openOffensiveTweets(wd)
		// openMoreReplies(wd)
		time.Sleep(500 * time.Millisecond)
		m = countThreads(wd)
	}
}

func openOffensiveTweets(wd *selenium.WebDriver) {
	offensive := findElementsByCSS("button.Tombstone-action.btn-link.ThreadedConversation-showMoreThreadsPrompt", wd)
	if offensive != nil {
		for _, o := range offensive {
			o.Click()
		}
	}
}

func openMoreReplies(wd *selenium.WebDriver) {
	clicks := findElementsByCSS("a.ThreadedConversation-moreRepliesLink", wd)
	if clicks != nil {
		for _, c := range clicks {
			c.Click()
		}
	}
}

func extractTweetInfo(tweet *Tweet, wd *selenium.WebDriver) bool {
	accountNode := findElementByCSS("div.tweet.permalink-tweet", wd)
	if accountNode == nil {
		return false
	}

	name, _ := accountNode.GetAttribute("data-name")
	handle, _ := accountNode.GetAttribute("data-screen-name")
	mentions, _ := accountNode.GetAttribute("data-mentions")

	var text string
	textNode := findChildElementByCSS("p.TweetTextSize", accountNode)
	if textNode != nil {
		text, _ = textNode.Text()
	}

	var nReplies int64
	repliesRoot := findChildElementByCSS("div.ProfileTweet-action.ProfileTweet-action--reply", accountNode)
	if repliesRoot != nil {
		repliesNode := findChildElementByCSS("span.ProfileTweet-actionCountForPresentation", repliesRoot)
		if repliesNode != nil {
			replies, _ := repliesNode.Text()
			nReplies, _ = strconv.ParseInt(replies, 10, 64)
		}
	}

	var nLikes int64
	likesRoot := findChildElementByCSS("div.ProfileTweet-action.ProfileTweet-action--favorite.js-toggleState", accountNode)
	if likesRoot != nil {
		likesNode := findChildElementByCSS("span.ProfileTweet-actionCountForPresentation", likesRoot)
		if likesNode != nil {
			likes, _ := likesNode.Text()
			nLikes, _ = strconv.ParseInt(likes, 10, 64)
		}
	}

	var nRetweets int64
	retweetsRoot := findChildElementByCSS("div.ProfileTweet-action.ProfileTweet-action--retweet.js-toggleState.js-toggleRt", accountNode)
	if retweetsRoot != nil {
		retweetsNode := findChildElementByCSS("span.ProfileTweet-actionCountForPresentation", retweetsRoot)
		if retweetsNode != nil {
			retweets, _ := retweetsNode.Text()
			nRetweets, _ = strconv.ParseInt(retweets, 10, 64)
		}
	}

	date := ""
	dateRoot := findChildElementByCSS("span.metadata", accountNode)
	if dateRoot != nil {
		dateRoot = findChildElementByCSS("span", dateRoot)
		date, _ = dateRoot.Text()
	}

	linkRoot := findChildElementByCSS("div.tweet.js-stream-tweet.js-actionable-tweet.js-profile-popup-actionable.dismissible-content.descendant.permalink-descendant-tweet", accountNode)
	linkStr := ""
	if linkRoot != nil {
		linkStr, _ = linkRoot.GetAttribute("data-permalink-path")
		linkStr = fmt.Sprintf("https://twitter.com%s", linkStr)
	}

	// fmt.Printf("found link %s\n", linkStr)

	tweet.Name = name
	tweet.TwitterHandle = handle
	tweet.Mentions = strings.Split(mentions, " ")
	tweet.Likes = int(nLikes)
	tweet.Retweets = int(nRetweets)
	tweet.Replies = int(nReplies)
	tweet.Text = text
	tweet.Link = linkStr
	tweet.Date = date
	return true
}

func getChildren(tweet Tweet, wd *selenium.WebDriver) []Tweet {
	var tweets []Tweet
	scrollAllTheWay(wd)

	list := findElementsByCSS("li.ThreadedConversation", wd)
	for _, thread := range list {
		threadTweets := findChildElementsByCSS("div.ThreadedConversation-tweet", thread)
		first := threadTweets[0]
		tweetInfo := findChildElementByCSS("div.tweet.js-stream-tweet", first)
		t := makeTweet()
		t.ID = id
		t.ParentID = tweet.ID
		str, _ := tweetInfo.GetAttribute("data-permalink-path")
		t.Link = fmt.Sprintf("https://twitter.com%s", str)
		id++
		tweets = append(tweets, t)
	}

	singletons := findElementsByCSS("li.ThreadedConversation--loneTweet", wd)

	for _, singleton := range singletons {
		tweetInfo := findChildElementByCSS("div.tweet.js-stream-tweet", singleton)
		t := makeTweet()
		t.ID = id
		t.ParentID = tweet.ID
		str, _ := tweetInfo.GetAttribute("data-permalink-path")
		t.Link = fmt.Sprintf("https://twitter.com%s", str)
		id++
		tweets = append(tweets, t)
	}
	return tweets
}

func parseTree(node Tweet, wd *selenium.WebDriver) Tweet {
	var kids []Tweet
	for _, child := range node.Children {
		(*wd).Get(child.Link)
		extractTweetInfo(&child, wd)
		child.Children = getChildren(child, wd)
		child = parseTree(child, wd)
		kids = append(kids, child)
	}
	node.Children = kids
	return node
}

func checkTime(date string) bool {
	splits := strings.Split(date, " ")
	if len(splits) < 6 {
		return false
	}
	dateStr := fmt.Sprintf("%s %s, %s at %s%s", splits[4], splits[3], splits[5], splits[0], strings.ToLower(splits[1]))

	const longForm = "Jan 2, 2006 at 3:04pm"
	d, err := time.Parse(longForm, dateStr)
	if err != nil {
		panic(err)
	}

	duration := time.Since(d)
	if duration.Hours() > 48 {
		return true
	}
	return false
}

func collectReplyTree(urls []string) []string {

	var r []string
	service, wd := randomLogin()
	defer service.Stop()
	defer wd.Close()

	for _, url := range urls {
		// fmt.Printf("Working on %s\n", url)
		parts := strings.Split(url, "/")
		idStr := parts[len(parts)-1]
		filename := fmt.Sprintf("data/%s.json", idStr)

		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("File %s already exists\n", filename)
		}

		openURL(url, &wd)
		node := findElementByCSS("div.permalink-inner", &wd)
		if node == nil {
			fmt.Printf("Problem with %s\n", url)
			continue
		}

		root := makeTweet()
		r := extractTweetInfo(&root, &wd)
		if r == false {
			fmt.Printf("Url %s has a problem\n", url)
			continue
		}

		if checkTime(root.Date) == false {
			fmt.Printf("The tweet %s is not 48 hours old.\n", url)
			continue
		}

		id = 0
		root.ID = 0
		root.Children = getChildren(root, &wd)
		root.Link = url
		root = parseTree(root, &wd)

		b, _ := json.Marshal(root)

		writeBytesToFile(filename, b)
		fmt.Printf("Wrote %s\n", filename)
	}
	err := wd.Quit()
	if err != nil {
		panic(err)
	}
	return r
}

func collectReplyTrees(cpus int) {
	spFile := readFileToList("data/starting_points.txt")
	jsonFiles := readDirContent("data/*.json")

	var startingPoints []string

	for _, entry := range spFile {
		found := false
		for _, json := range jsonFiles {
			entries := strings.Split(entry, "/")
			e := entries[len(entries)-1]

			js := strings.Split(json, "/")[1]
			js = strings.Split(js, ".")[0]

			if e == js {
				found = true
				break
			}
		}
		if found == false {
			startingPoints = append(startingPoints, entry)
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(startingPoints), func(i, j int) { startingPoints[i], startingPoints[j] = startingPoints[j], startingPoints[i] })

	var wg sync.WaitGroup
	send := make(chan []string, cpus*2)
	ans := make(chan []string, cpus*2)

	for i := 0; i < cpus; i++ {
		wg.Add(1)
		go func(send <-chan []string, ans chan<- []string) {
			defer wg.Done()
			for tweet := range send {
				ans <- collectReplyTree(tweet)
			}
		}(send, ans)
	}

	// start the jobs
	go func(send chan<- []string) {
		for i := 0; i < len(startingPoints)-5; i += 5 {
			send <- startingPoints[i : i+5]
		}
		close(send)
		wg.Wait()
		close(ans)
	}(send)

	for a := range ans {
		for _, v := range a {
			fmt.Printf("%s\n", v)
		}
	}
	// }
}
