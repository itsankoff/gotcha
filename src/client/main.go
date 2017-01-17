package main

import (
    "fmt"
    "gotcha/client"
)

func main() {
    client := gotcha.NewClient()
    client.Connect("wss://127.0.0.1:9999")
    fmt.Println("Hello client")
}
