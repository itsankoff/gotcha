package main

import (
	"github.com/itsankoff/gotcha/client"
	"log"
)

func main() {
	ws := client.NewWebSocketClient()
	c := client.New(ws)
	err := c.Connect("ws://127.0.0.1:9000/websocket")
	log.Println("connected", err)
	userId, err := c.Register("pesho", "123")
	log.Println("registered", err)

	err = c.Authenticate(userId, "123")
	log.Println("authenticated", err)

	if err == nil {
		c.StartInteractiveMode()
	}
}
