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

func getThreadEntry(we selenium.WebElement, rootNode bool) (Tweet, bool) {
	tweet := makeTweet()
	element := findChildElementByCSS("div.tweet", we)
	class, _ := element.GetAttribute("class")
	if strings.Contains(class, "withheld-tweet") {
		return tweet, false
	}

	tweet.Link, _ = element.GetAttribute("data-permalink-path")
	mentions, _ := element.GetAttribute("data-mentions")
	name, _ := element.GetAttribute("data-name")
	handle, _ := element.GetAttribute("data-screen-name")
	userIDStr, _ := element.GetAttribute("data-user-id")
	linkStr, _ := element.GetAttribute("data-permalink-path")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	var date string
	if rootNode == true {
		dateRoot := findChildElementByCSS("span.metadata", we)
		date, _ = dateRoot.Text()
	} else {
		dateRoot := findChildElementByCSS("a.tweet-timestamp", element)
		date, _ = dateRoot.GetAttribute("data-original-title")
	}

	var text string
	textNode := findChildElementByCSS("p.TweetTextSize", we)
	if textNode != nil {
		text, _ = textNode.Text()
	}

	footer := findChildElementByCSS("div.stream-item-footer", we)

	replyNode := findChildElementByCSS("div.ProfileTweet-action.ProfileTweet-action--reply", footer)
	replyNrNode := findChildElementByCSS("span.ProfileTweet-actionCountForPresentation", replyNode)
	replyStr, _ := replyNrNode.Text()
	nReplies, _ := strconv.ParseInt(replyStr, 10, 64)

	retweetNode := findChildElementByCSS("div.ProfileTweet-action.ProfileTweet-action--retweet.js-toggleState.js-toggleRt", footer)
	retweetNrNode := findChildElementByCSS("span.ProfileTweet-actionCountForPresentation", retweetNode)
	retweetStr, _ := retweetNrNode.Text()
	nRetweets, _ := strconv.ParseInt(retweetStr, 10, 64)

	likesNode := findChildElementByCSS("div.ProfileTweet-action.ProfileTweet-action--favorite.js-toggleState", footer)
	likesNrNode := findChildElementByCSS("span.ProfileTweet-actionCountForPresentation", likesNode)
	likesStr, _ := likesNrNode.Text()
	nLikes, _ := strconv.ParseInt(likesStr, 10, 64)

	tweet.Mentions = strings.Split(mentions, " ")
	tweet.Name = name
	tweet.TwitterHandle = handle
	tweet.Mentions = strings.Split(mentions, " ")
	tweet.Text = text
	tweet.TwitterID = int(userID)
	tweet.Link = fmt.Sprintf("https://twitter.com%s", linkStr)
	tweet.Date = date
	tweet.Likes = int(nLikes)
	tweet.Retweets = int(nRetweets)
	tweet.Replies = int(nReplies)

	return tweet, true
}

func getChildren(root Tweet, wd *selenium.WebDriver) []Tweet {
	var tweets []Tweet
	scrollAllTheWay(wd)

	li := findElementsByCSS("li", wd)
	for _, element := range li {
		class, _ := element.GetAttribute("class")
		if strings.Contains(class, "ThreadedConversation") == true &&
			strings.Contains(class, "moreReplies") == false {
			tweet, r := getThreadEntry(element, false)
			if strings.Contains(class, "loneTweet") == true {
				tweet.Lone = true
			}
			if r == true {
				tweets = append(tweets, tweet)
			}
		}
	}

	return tweets
}

func parseTree(node Tweet, wd *selenium.WebDriver) Tweet {
	var kids []Tweet
	for _, child := range node.Children {
		if child.Lone == false {
			(*wd).Get(child.Link)
			child.Children = getChildren(child, wd)
			child = parseTree(child, wd)
		}
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
			continue
		}

		openURL(url, &wd)
		node := findElementByCSS("div.permalink-inner", &wd)
		if node == nil {
			fmt.Printf("Problem with %s\n", url)
			continue
		}

		we := findElementByCSS("div.permalink-inner.permalink-tweet-container", &wd)
		root, _ := getThreadEntry(we, true)

		if checkTime(root.Date) == false {
			fmt.Printf("The tweet %s is not 48 hours old: %s.\n", url, root.Date)
			continue
		}

		root.Children = getChildren(root, &wd)

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
	spFile := readFileToList("watch/starting_points.txt")
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
