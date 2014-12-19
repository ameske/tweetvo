package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

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
	log.Printf("%v", api)
}
