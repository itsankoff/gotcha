package main

import (
    "github.com/itsankoff/gotcha/client"
)

func main() {
    c := client.New()
    c.Connect("wss://127.0.0.1:9999")
}
