package main

import "net/http"

// TwitterClient is the means through which one interacts with the Twitter API
type TwitterClient struct {
	client *http.Client
	config *OAuthConfig
}

// NewTwitterClient creates a ready to use TwitterClient that will utilize the given
// OAuth tokens to authenticate its requests
func NewTwitterClient(c *OAuthConfig) *TwitterClient {
	return &TwitterClient{
		client: &http.Client{},
		config: c,
	}
}
