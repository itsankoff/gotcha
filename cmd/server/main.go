package main

import (
    "github.com/itsankoff/gotcha/server"
    "time"
    "log"
)

func main() {
    s := server.New()
    wss := server.NewWebSocket()
    s.AddTransport("127.0.0.1:9000", &wss)
    done := make(chan interface{})
    go func() {
        log.Println("Will close done channel")
        time.Sleep(10 * time.Second)
        log.Println("Close done channel")
        close(done)
    }()

    err := s.Start(done)
    if err != nil {
        log.Fatal("Failed to start server")
    }
}
