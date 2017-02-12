package server

import (
    "net/http"
    "github.com/gorilla/websocket"
    "log"
    "fmt"
)

type WebSocketServer struct {
    upgrader        websocket.Upgrader
    connections     []*websocket.Conn
}

func New() WebSocketServer {
    var upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool { return true },
    }

    return WebSocketServer{upgrader:upgrader}
}

func (wss *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    conn, err := wss.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    go wss.webSocketHandler(conn)
}

func (wss *WebSocketServer) webSocketHandler(conn *websocket.Conn) {
    log.Println("New websocket connection available")
    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Connection read error", err.Error())
            return
        }

        if msgType == websocket.TextMessage {
            log.Println(string(msg))
        }

        if err = conn.WriteMessage(msgType, msg); err != nil {
            log.Println("Connection write error", err.Error())
            return
        }
    }
}

func (wss *WebSocketServer) Start(host string) {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello World")
    })

    http.Handle("/websocket", wss)
    log.Println("Listen on:", host)
    log.Fatal(http.ListenAndServe(host, nil))
}
