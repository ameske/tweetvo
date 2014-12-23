package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	config OAuthConfig
)

func init() {
	c, err := ioutil.ReadFile("tweetvo.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = yaml.Unmarshal(c, &config)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func main() {
	delayedTweets := make(chan Tweet, 1024)
	go delay(delayedTweets)

	req, err := http.NewRequest("GET", "https://userstream.twitter.com/1.1/user.json?with=true", nil)
	if err != nil {
		log.Fatalf("NewRequest: %s", err.Error())
	}

	err = config.Sign(req, config.header())
	if err != nil {
		log.Fatalf("Sign: %s", err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Do: %s", err.Error())
	}
	defer resp.Body.Close()

	// Continually parse messages from the Twitter streaming endpoint
	respReader := bufio.NewReader(resp.Body)
	for {
		message, err := readMessage(respReader)
		if err != nil {
			log.Fatal(err.Error())
		}

		processMessage(message, delayedTweets)
	}
}

// readMessage tokenizes messages from the Twitter streaming endpoint.
func readMessage(r *bufio.Reader) ([]byte, error) {
	var msg bytes.Buffer

	isPrefix := true
	for isPrefix {
		var err error
		var bytes []byte
		bytes, isPrefix, err = r.ReadLine()
		if err != nil {
			return nil, err
		}

		msg.Write(bytes)
	}

	return msg.Bytes(), nil
}

// processMessage determines whether or not the message is a tweet or some other
// informration. If the message is a tweet it is sent to the delay function, otherwise
// it is printed out.
func processMessage(message []byte, delay chan Tweet) {
	if len(message) == 0 {
		return
	}

	if !strings.Contains(string(message), "text") {
		fmt.Printf("%s\n", string(message))
		return
	}

	var t Tweet
	err := json.Unmarshal(message, &t)
	if err != nil {
		log.Fatalf(err.Error())
	}

	delay <- t
}

// delay applies a delay to incoming tweets
func delay(tweets chan Tweet) {
	for {
		t := <-tweets

		tweetTime, err := time.Parse("Mon Jan 2 15:04:05 +0000 2006", t.CreatedAt)
		if err != nil {
			log.Fatalf(err.Error())
		}
		fmt.Printf("Live: %d\n", tweetTime.Unix())

		adjTweetTime := tweetTime.Add(-time.Second * 30)
		if adjTweetTime.Before(time.Now()) {
			<-time.After(time.Now().Sub(adjTweetTime))
		}

		fmt.Printf("Delayed: %d\n", time.Now().Unix())
		fmt.Printf("%s\t%s: %s\n", t.CreatedAt, t.User.ScreenName, t.Text)
	}
}
