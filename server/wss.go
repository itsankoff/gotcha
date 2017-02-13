package server

import (
    "log"
    "time"
    "strconv"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/itsankoff/gotcha/util"
)

type WebSocketServer struct {
    upgrader            websocket.Upgrader
    connections         map[*util.User]*websocket.Conn
    connected           chan<- *util.User
    disconnected        chan<- *util.User
}

func NewWebSocket() WebSocketServer {
    var upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool { return true },
    }

    return WebSocketServer{
        upgrader:upgrader,
        connections: make(map[*util.User]*websocket.Conn),
    }
}

func (wss *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    conn, err := wss.upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    wss.addConnection(conn)
}

func (wss *WebSocketServer) addConnection(conn *websocket.Conn) {
    now := time.Now().Unix()
    id := strconv.FormatInt(now, 10)
    user := &util.User{
        Id: id,
        In: make(chan util.Message),
        Out: make(chan util.Message),
    }

    wss.connections[user] = conn
    go wss.inputHandler(user, conn)
    go wss.outputHandler(user, conn)
    wss.connected <- user
    log.Println("Add connections", user.Id)
}

func (wss *WebSocketServer) removeConnection(conn *websocket.Conn) {
    for user, c := range wss.connections {
        if c == conn {
            if wss.disconnected != nil {
                wss.disconnected <- user
            }

            close(user.In)
            close(user.Out)
            conn.Close()
            delete(wss.connections, user)

            log.Println("Remove connection", user.Id)
            break
        }
    }
}

func (wss *WebSocketServer) inputHandler(user *util.User, conn *websocket.Conn) {
    log.Println("Start websocket input handler for", user.Id)
    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Connection read error", err.Error())
            wss.removeConnection(conn)
            return
        }

        message, err := wss.decodeMessage(user, msg, msgType)
        if err != nil {
            log.Println("Failed to decode message", msgType, msg)
            wss.removeConnection(conn)
            return
        }

        user.In <- message
    }
}

func (wss *WebSocketServer) outputHandler(user *util.User, conn *websocket.Conn) {
    log.Println("Start websocket output handler for", user.Id)
    select {
    case msg := <-user.Out:
        message, msgType := wss.encodeMessage(user, msg)
        if err := conn.WriteMessage(msgType, message); err != nil {
            log.Println("Connection write error", err.Error())
            wss.removeConnection(conn)
            return
        }
    }
}

func (wss *WebSocketServer) Start(host string, done <-chan interface{}) {
    subPath := "/websocket"
    http.Handle(subPath, wss)
    defer func() {
        http.Handle(subPath, nil)
    }()

    log.Println("Listen on:", host + subPath)
    log.Fatal(http.ListenAndServe(host, nil))
}

func (wss *WebSocketServer) OnUserConnected(handler chan<- *util.User) {
    wss.connected = handler
}

func (wss *WebSocketServer) OnUserDisconnected(handler chan<- *util.User) {
    wss.disconnected = handler
}

func (wss WebSocketServer) encodeMessage(u *util.User,
                                         msg util.Message) ([]byte, int) {
    return msg.Raw(), int(msg.DataType())
}

func (wss WebSocketServer) decodeMessage(u *util.User,
                                         data []byte,
                                         dataType int) (util.Message, error) {
    // TODO: Parse message data
    //       And find who is the destination
    message := util.NewMessage(u, u, "message", util.DataType(dataType), data)
    return message, nil
}
