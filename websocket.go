package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// establishWebsocket receives requests from clients to connect to
// the server via a websocket, and establishes the connection.
func establishWebsocket(tweets chan Tweet) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		go disregardClient(conn)
		go serveWebTweets(conn, tweets)
	}
}

func disregardClient(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			break
		}
	}
}

func serveWebTweets(c *websocket.Conn, tweets chan Tweet) {
	for {
		t := <-tweets
		delay(t, time.Second*30)

		msg := fmt.Sprintf("%s\t%s: %s\n", t.CreatedAt, t.User.ScreenName, t.Text)
		if err := c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Fatal(err)
		}
	}
}
