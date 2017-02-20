package server

import (
	"log"
	"strconv"
	"testing"
	"time"
)

func getConfig() *Config {
	return &Config{
		ListenHost:     "127.0.0.1:9000",
		FileServerHost: "127.0.0.0.1:9000",
		FileServerPath: "/tmp" + strconv.FormatInt(time.Now().UnixNano(), 10),
	}
}

func TestServer_AddTrasnport(t *testing.T) {
	s := New(getConfig())
	wss := NewWebSocket(getConfig())
	err := s.AddTransport(":9000", &wss)
	if err != nil {
		t.Error(err)
	}
}

func TestServer_AddTransport_EmptyHost(t *testing.T) {
	s := New(getConfig())
	wss := NewWebSocket(getConfig())
	err := s.AddTransport("", &wss)
	if err == nil {
		t.Error("Need error when adding transport for empty host")
	}
}

func TestServer_AddTransport_NoTransport(t *testing.T) {
	s := New(getConfig())
	err := s.AddTransport(":9000", nil)
	if err == nil {
		t.Error("Need error when adding a nil transport for a host")
	}
}

func TestServer_Start(t *testing.T) {
	s := New(getConfig())
	done := make(chan interface{})
	err := s.Start(done)
	if err == nil {
		t.Error("Need error if trying to start server without any transport")
	}
}

func TestServer_StartAsync(t *testing.T) {
	s := New(getConfig())
	_, err := s.StartAsync()
	if err == nil {
		t.Error("Need error if trying to start server without any transport")
	}
}

func TestServer_RemoveTransport(t *testing.T) {
	s := New(getConfig())
	wss := NewWebSocket(getConfig())
	err := s.AddTransport(":9000", &wss)
	if err != nil {
		t.Error(err)
	}

	err = s.RemoveTransport(":9000")
	if err != nil {
		t.Error(err)
	}
}

func ExampleServer_StartAsync() {
	s := New(getConfig())
	wss := NewWebSocket(getConfig())
	s.AddTransport(":9999", &wss)
	done, err := s.StartAsync()
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		time.Sleep(time.Second * 3)
		close(done)
	}()
	<-done
}
