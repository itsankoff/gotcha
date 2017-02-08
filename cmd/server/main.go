package main

import (
    "github.com/itsankoff/gotcha/server"
)

func main() {
    server:= server.New()
    server.Start("127.0.0.1:9000")
}
