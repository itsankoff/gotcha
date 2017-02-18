package client

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/itsankoff/gotcha/common"
	"log"
	"time"
)

type WebSocketClient struct {
	In         chan *common.Message
	serverHost string
	conn       *websocket.Conn
}

func NewWebSocketClient() *WebSocketClient {
	return &WebSocketClient{}
}

func (ws *WebSocketClient) inputHandler() {
	defer ws.conn.Close()
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message", err)
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(msg, &data)
		if err != nil {
			log.Println("Failed to decode message", err)
			return
		}

		from := data["from"].(string)
		to := data["to"].(string)
		cmdType := data["cmd_type"].(string)
		cmd := data["cmd"].(string)
		dataType := common.DataType(data["data_type"].(float64))
		// TODO: add expire_date

		message := common.NewMessage(from, to, cmdType, cmd,
			time.Time{}, dataType, data["data"])
		ws.In <- &message
	}
}

func (ws *WebSocketClient) Connect(host string) error {
	if ws.In == nil {
		return errors.New("In channled not initialized. Call SetReceiver")
	}
	c, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		log.Println("Failed to connect to server", err)
		return err
	}
	ws.serverHost = host
	ws.conn = c
	go ws.inputHandler()
	return nil
}

func (ws *WebSocketClient) ConnectAsync(host string) chan bool {
	connected := make(chan bool)
	go func() {
		err := ws.Connect(host)
		if err != nil {
			connected <- false
		}

		connected <- true
	}()

	return connected
}

func (ws *WebSocketClient) Disconnect() {
	if ws.conn != nil {
		ws.conn.Close()
	}
}

func (ws *WebSocketClient) Reconnect() error {
	if ws.serverHost == "" {
		return errors.New("No server host provided")
	}

	ws.Disconnect()
	return ws.Connect(ws.serverHost)
}

func (ws *WebSocketClient) ReconnectAsync() chan bool {
	reconnected := make(chan bool)
	go func() {
		ws.Disconnect()
		err := ws.Connect(ws.serverHost)
		if err != nil {
			log.Println("Failed to reconnect", err)
			reconnected <- false
		} else {
			log.Println("Reconnected to", ws.serverHost)
			reconnected <- true
		}
	}()

	return reconnected
}

func (ws *WebSocketClient) SendText(message string) error {
	err := ws.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Failed to send text message", err)
	}

	return err
}

func (ws *WebSocketClient) SendBinary(data []byte) error {
	err := ws.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("Failed to send binary message", err)
	}

	return err
}

func (ws *WebSocketClient) SetReceiver(input chan *common.Message) {
	ws.In = input
}
