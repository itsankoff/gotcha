package main

import (
	"flag"
	"github.com/itsankoff/gotcha/server"
	"log"
)

func main() {
	config := server.NewConfig()
	config.ListenHost = *flag.String("host", ":9000", "host to listen")
	config.FileServerHost = *flag.String("file_host", ":9000", "host to server files")
	config.FileServerPath = *flag.String("file_path", "/files", "query file path to access files")
	config.FileServerFolder = *flag.String("file_folder", "./", "storage folder")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 && args[0] == "--help" {
		flag.PrintDefaults()
		return
	}

	srv := server.New(config)
	wss := server.NewWebSocket()
	srv.AddTransport("127.0.0.1:9000", &wss)
	done := make(chan interface{})

	err := srv.Start(done)
	if err != nil {
		log.Fatal("Failed to start server")
	}
}
