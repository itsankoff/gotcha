package main

import (
	"github.com/itsankoff/gotcha/client"
	"github.com/itsankoff/gotcha/common"
	"log"
	"time"
)

func main() {
	input := make(chan *common.Message)
	ws := client.NewWebSocketClient()
	ws.SetReceiver(input)

	c := client.New(ws)
	err := c.Connect("ws://127.0.0.1:9000/websocket")
	log.Println("connected", err)
	userId, err := c.Register("pesho", "123")
	log.Println("registered", err)

	err = c.Authenticate(userId, "123")
	log.Println("authenticated", err)

	time.Sleep(10 * time.Second)
}
