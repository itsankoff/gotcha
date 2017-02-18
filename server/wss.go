package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/itsankoff/gotcha/common"
	"log"
	"net/http"
	"strconv"
	"time"
)

type WebSocketServer struct {
	upgrader     websocket.Upgrader
	connections  map[*common.User]*websocket.Conn
	connected    chan<- *common.User
	disconnected chan<- *common.User
}

func NewWebSocket() WebSocketServer {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return WebSocketServer{
		upgrader:    upgrader,
		connections: make(map[*common.User]*websocket.Conn),
	}
}

func (wss *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("New http connection. Try to upgrade to WebSocket")
	conn, err := wss.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket", err)
		fmt.Fprintf(w, "failed")
	}

	wss.addConnection(conn)
}

func (wss *WebSocketServer) addConnection(conn *websocket.Conn) {
	now := time.Now().UnixNano()
	id := strconv.FormatInt(now, 10)
	user := common.NewUser(id)

	wss.connections[user] = conn
	log.Println("Add websocket connection", user.Id)

	go wss.inputHandler(user, conn)
	go wss.outputHandler(user, conn)
	go wss.closeHandler(user, conn)
	wss.connected <- user
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
			log.Println("Remove websocket connection", user.Id)
			break
		}
	}
}

func (wss *WebSocketServer) inputHandler(user *common.User, conn *websocket.Conn) {
	log.Println("Start websocket input handler for", user.Id)
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Websocket connection read error", err.Error())
			wss.removeConnection(conn)
			return
		}

		message, err := wss.decodeMessage(user, msg, msgType)
		// 		log.Println("Websocket message", user.Id)
		if err != nil {
			log.Println("Failed to decode message", msgType, msg)
			wss.removeConnection(conn)
			return
		}

		user.In <- message
	}
}

func (wss *WebSocketServer) outputHandler(user *common.User, conn *websocket.Conn) {
	log.Println("Start websocket output handler for", user.Id)
	for {
		select {
		case msg := <-user.Out:
			if msg == nil {
				log.Printf("Nil message in output channel for %s.", user.Id)
				log.Printf("Stop websocket output handler %s", user.Id)
				return
			}

			// 			log.Println("Message in output channel for", user.Id)
			message, msgType := wss.encodeMessage(user, msg)
			if err := conn.WriteMessage(msgType, message); err != nil {
				log.Println("Websocket connection write error", err.Error())
				wss.removeConnection(conn)
				return
			}
		}
	}
}

func (wss *WebSocketServer) closeHandler(user *common.User, conn *websocket.Conn) {
	log.Println("Start websocket close handler for", user.Id)
	select {
	case <-user.Done:
		log.Println("Done channel closed for user", user.Id)
		wss.removeConnection(conn)
	}
}

func (wss *WebSocketServer) Start(host string, done <-chan interface{}) {
	subPath := "/websocket"
	http.Handle(subPath, wss)
	defer func() {
		http.Handle(subPath, nil)
	}()

	log.Println("WebSocket Server Listens on:", host+subPath)
	log.Fatal(http.ListenAndServe(host, nil))
}

func (wss *WebSocketServer) OnUserConnected(handler chan<- *common.User) {
	wss.connected = handler
}

func (wss *WebSocketServer) OnUserDisconnected(handler chan<- *common.User) {
	wss.disconnected = handler
}

func (wss WebSocketServer) encodeMessage(u *common.User,
	msg *common.Message) ([]byte, int) {
	json, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode message", err)
		return []byte{}, 0
	}

	return json, int(msg.DataType())
}

func (wss WebSocketServer) decodeMessage(u *common.User,
	data []byte,
	dataType int) (*common.Message, error) {
	var packet map[string]interface{}
	err := json.Unmarshal(data, &packet)
	if err != nil {
		return nil, err
	}

	expire_period, exists := packet["expire_period"]
	var expire_date time.Time
	if exists {
		expire_date = time.Now().Add(time.Duration(expire_period.(float64)) * time.Second)
	}

	messageFrom := packet["from"].(string)
	messageTo := packet["to"].(string)
	messageCmdType := packet["cmd_type"].(string)
	messageCmd := packet["cmd"].(string)
	messageDataType := common.DataType(packet["data_type"].(float64))

	message := common.NewMessage(messageFrom, messageTo,
		messageCmdType, messageCmd,
		expire_date,
		messageDataType, packet["data"])
	return &message, nil
}
