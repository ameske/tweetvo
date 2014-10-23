package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

type OAuthConfig struct {
	consumerKey string `yaml:"CONSUMER_KEY"`
	accesToken  string `yaml:"ACCESS_TOKEN"`
	version     string `yaml:"VERSION"`
}

type OAuthHeader struct {
	ConsumerKey     string
	Nonce           string
	Signature       string
	SignatureMethod string
	Timestamp       string
	AccessToken     string
	Version         string
}

func OAuthNonce() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Int31()
	var b bytes.Buffer
	binary.PutVariant(b, n)

	return base64.URLEncoding.EncodeToString(b)
}

type TwitterClient struct {
	client *http.Client
	config *OAuthConfig
}

func NewTwitterClient(c *OAuthConfig) *TwitterClient {
	return &TwitterClient{
		client: &http.Client{},
		config: c,
	}
}

func (t *TwitterClient) GetOneTweet() {
	url := "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=Kyle_Ames_CS&count=1"
	// Include screen_name = "Kyle_Ames_CS"
	// Include count = 1

	r := http.NewRequest("GET", url, nil)
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Connection", "close")
	r.Header.Add("User-Agent", "Tweetvo")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Host", "api.twitter.com")

	h := OAuthHeader{
		ConsumerKey:     config.consumerKey,
		Nonce:           OAuthNonce(),
		Signature:       "",
		SignatureMethod: "HMAC-SHA1",
		Timestamp:       string(time.Now().Unix()),
		Token:           config.accessToken,
		Version:         config.version,
	}

}

func main() {
	// Read in OAuth information from a file to prevent leaking it on github ;-)
	c, err := ioutil.ReadFile("tweetvo.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}

	config := &OAuthConfig{}
	err = yaml.Unmarshal(c, config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	api := NewTwitterClient(config)

}
