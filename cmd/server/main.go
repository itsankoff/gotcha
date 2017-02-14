package main

import (
    "github.com/itsankoff/gotcha/server"
    "log"
)

func main() {
    s := server.New()
    wss := server.NewWebSocket()
    s.AddTransport("127.0.0.1:9000", &wss)
    done := make(chan interface{})

    err := s.Start(done)
    if err != nil {
        log.Fatal("Failed to start server")
    }
}
