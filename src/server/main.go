package main

import (
    "gotcha/server"
)

func main() {
    server:= gotcha.New()
    server.Start("wss://127.0.0.1:9999")
}
