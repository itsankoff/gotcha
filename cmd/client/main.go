package main

import (
	"flag"
	"github.com/itsankoff/gotcha/client"
	"log"
)

func main() {
	var host string
	flag.StringVar(&host, "host",
		"ws://0.0.0.0:9000/websocket", "remote server host")

	flag.Parse()

	ws := client.NewWebSocketClient()
	c := client.New(ws)
	err := c.Connect(host)
	log.Println("connected", err)
	userId, err := c.Register("pesho", "123")
	log.Println("registered", err)

	err = c.Authenticate(userId, "123")
	log.Println("authenticated", err)

	if err == nil {
		c.StartInteractiveMode()
	}
}
