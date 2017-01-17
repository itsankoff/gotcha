package main

import (
    "gotcha/client"
)

func main() {
    c := gotcha.New()
    c.Connect("wss://127.0.0.1:9999")
}
