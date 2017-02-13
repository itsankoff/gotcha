package server_test

import (
    "github.com/itsankoff/gotcha/server"
    "testing"
    "time"
    "log"
)

func TestServer_AddTrasnport(t *testing.T) {
    s := server.New()
    wss := server.NewWebSocket()
    err := s.AddTransport(":9000", &wss)
    if err != nil {
        t.Error(err)
    }
    // Output:
}

func TestServer_AddTransport_EmptyHost(t *testing.T) {
    s := server.New()
    wss := server.NewWebSocket()
    err := s.AddTransport("", &wss)
    if err == nil {
        t.Error("Need error when adding transport for empty host")
    }
    // Output:
}

func TestServer_AddTransport_NoTransport(t *testing.T) {
    s := server.New()
    err := s.AddTransport(":9000", nil)
    if err == nil {
        t.Error("Need error when adding a nil transport for a host")
    }
    // Output:
}

func TestServer_Start(t *testing.T) {
    s := server.New()
    done := make(chan interface{})
    err := s.Start(done)
    if err == nil {
        t.Error("Need error if trying to start server without any transport")
    }
    // Output:
}

func TestServer_StartAsync(t *testing.T) {
    s := server.New()
    _, err := s.StartAsync()
    if err == nil {
        t.Error("Need error if trying to start server without any transport")
    }
    // Output:
}

func TestServer_RemoveTransport(t *testing.T) {
    s := server.New()
    wss := server.NewWebSocket()
    err := s.AddTransport(":9000", &wss)
    if err != nil {
        t.Error(err)
    }

    err = s.RemoveTransport(":9000")
    if err != nil {
        t.Error(err)
    }
    // Output:
}

func ExampleServer_StartAsync() {
    s := server.New()
    wss := server.NewWebSocket()
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
    // Output:
}
