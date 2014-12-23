package main

import (
	"fmt"
	"time"
)

func serveStdoutTweets(tweets chan Tweet) {
	t := <-tweets
	delay(t, time.Second*30)
	fmt.Printf("%s\t%s: %s\n", t.CreatedAt, t.User.ScreenName, t.Text)
}
