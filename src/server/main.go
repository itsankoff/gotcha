package main

import (
    "fmt"
    "gotcha/server"
)

func main() {
    server:= gotcha.NewServer()
    server.Start()
    fmt.Println("Hello server")
}
