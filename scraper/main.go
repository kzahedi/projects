package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	yaml "gopkg.in/yaml.v2"
)

type twitterConfig struct {
	APIKey         string `yaml:"API key"`
	APISecretKey   string `yaml:"API secret key"`
	TokenKey       string `yaml:"Access token"`
	TokenSecretKey string `yaml:"Access token secret"`
}

func readLines(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	return lines
}

func main() {
	idPtr := flag.Int("i", 950882597881737216, "Twitter ID")
	configFilePtr := flag.String("c", "", "Config file")
	flag.Parse()

	fmt.Println("Reading " + (*configFilePtr))
	data, err := ioutil.ReadFile(*configFilePtr)
	if err != nil {
		panic(err)
	}

	cfg := twitterConfig{}
	err = yaml.Unmarshal(data, &cfg)
	fmt.Println(cfg)
	if err != nil {
		panic(err)
	}

	config := oauth1.NewConfig(cfg.APIKey, cfg.APISecretKey)
	token := oauth1.NewToken(cfg.TokenKey, cfg.TokenSecretKey)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Home Timeline
	// tweets, resp, err := client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
	// 	Count: 20,
	// })

	lastCount := -1
	var timeline []twitter.Tweet

	for lastCount != len(timeline) && len(timeline) < 3200 {
		// User Timeline
		lastCount = len(timeline)
		var tweets []twitter.Tweet
		if len(timeline) == 0 {
			tweets, _, _ = client.Timelines.UserTimeline(&twitter.UserTimelineParams{
				UserID: int64(*idPtr),
				Count:  3200,
			})
		} else {
			tweet := timeline[len(timeline)-1]
			tweets, _, _ = client.Timelines.UserTimeline(&twitter.UserTimelineParams{
				UserID: int64(*idPtr),
				Count:  3200,
				MaxID:  tweet.ID,
			})

		}
		timeline = append(timeline, tweets...)
		fmt.Println(len(timeline))
	}

	// fmt.Println(tweets)
	// fmt.Println(len(tweets))
	// fmt.Println(resp)
	// fmt.Println(err)

	bytes, _ := json.Marshal(timeline)
	f, err := os.Create(fmt.Sprintf("%d.json", *idPtr))
	defer f.Close()
	if err != nil {
		panic(err)
	}
	f.Write(bytes)
}
