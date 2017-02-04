package main

import (
    "github.com/itsankoff/gotcha/server"
)

func main() {
    server:= server.New()
    server.Start("wss://127.0.0.1:9999")
}
