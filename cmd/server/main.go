package main

import (
	"flag"
	"github.com/itsankoff/gotcha/server"
	"log"
)

func main() {
	config := server.NewConfig()
	flag.StringVar(&config.ListenHost, "host",
		"0.0.0.0:9000", "host to listen")

	flag.StringVar(&config.FileServerHost, "file_host",
		"http://0.0.0.0:9000", "host to server files")

	flag.StringVar(&config.FileServerPath, "file_path",
		"/", "query path to access files")

	flag.StringVar(&config.FileServerFolder, "file_folder",
		"./", "storage folder")

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
