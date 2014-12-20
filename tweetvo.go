package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	streaming()
}

func streaming() {
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

	respReader := bufio.NewReader(resp.Body)
	for {
		message, err := readMessage(respReader)
		if err != nil {
			log.Fatal(err.Error())
		}

		go processMessage(message)
	}
}

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

func processMessage(message []byte) {
	fmt.Printf("%s\n", string(message))
}
