// Package common declares common type which are used in both
// server and client
package common

import (
	"encoding/json"
	"log"
	"time"
)

const (
	TEXT   = 1
	BINARY = 2
)

const (
	STATUS_OK    = 1
	STATUS_ERROR = 2
)

type DataType int
type Status int

// Message represents protocol message
// For more detailed information check protocol.txt
type Message struct {
	from string
	to   string

	cmdType string
	cmd     string

	status Status // only for server -> client messages

	expireDate time.Time

	dataType DataType
	data     interface{}
}

// NewMessage creates a message from params
func NewMessage(from string, to string,
	cmdType string, cmd string,
	expireDate time.Time,
	dataType DataType, data interface{}) Message {
	return Message{
		from:       from,
		to:         to,
		cmdType:    cmdType,
		cmd:        cmd,
		expireDate: expireDate,
		dataType:   dataType,
		data:       data,
	}
}

// String parses message payload to string if
// message data type is TEXT
func (m Message) String() string {
	if m.dataType != TEXT {
		log.Println("Message data is not a text type")
		return ""
	}

	return m.data.(string)
}

// Binary parses message payload to binary if message
// data type is BINARY
func (m Message) Binary() []byte {
	if m.dataType != BINARY {
		log.Println("Message data is not binary type")
		return []byte{}
	}

	return m.data.([]byte)
}

// Json encodes message in json format
func (m Message) Json() ([]byte, error) {
	msg := make(map[string]interface{})
	msg["from"] = m.from
	msg["to"] = m.to
	msg["cmd_type"] = m.cmdType
	msg["cmd"] = m.cmd
	msg["data_type"] = m.dataType
	msg["data"] = m.data

	return json.Marshal(msg)
}

// ParseJsonData tries to parse message payload as
// map[string]interface{}
func (m Message) ParseJsonData() (map[string]interface{}, error) {
	jdata, ok := m.data.(map[string]interface{})
	if ok {
		return jdata, nil
	}

	return make(map[string]interface{}), nil
}

// GetJsonData tries to extract data from message json payload
// by key
func (m Message) GetJsonData(key string) interface{} {
	jdata, err := m.ParseJsonData()
	if err != nil {
		return nil
	}

	val, _ := jdata[key]
	return val
}

// From returns message sender
func (m Message) From() string {
	return m.from
}

// To returns message receiver
func (m Message) To() string {
	return m.to
}

// CmdType returns command type. See protocol.txt for
// more information
func (m Message) CmdType() string {
	return m.cmdType
}

// Cmd returns the message command. See protocol.txt for
// more information
func (m Message) Cmd() string {
	return m.cmd
}

// Status returns the message status. VALID ONLY FOR MESSAGES FROM SERVER
func (m Message) Status() Status {
	return m.status
}

// Error tries to parse message payload as error message
// if message status is STATUS_ERROR. VALID ONLY FOR MESSAGES FROM SERVER
func (m Message) Error() string {
	jdata, err := m.ParseJsonData()
	if err != nil {
		return ""
	}

	errMsg, _ := jdata["error"].(string)
	return errMsg
}

// ExpireDate returns message expire date
func (m Message) ExpireDate() time.Time {
	return m.expireDate
}

// DataType return message payload data type. (TEXT or BINARY)
func (m Message) DataType() DataType {
	return m.dataType
}

// Helper function to generate responses for a message. USED only in server
func NewResponse(msg *Message, key string, response interface{}) *Message {
	payload := make(map[string]interface{})
	payload[key] = response
	var status Status
	if key == "error" {
		status = STATUS_ERROR
	} else {
		status = STATUS_OK
	}

	responseMsg := NewMessage("server", msg.from, msg.cmdType, msg.cmd,
		time.Time{}, TEXT, payload)
	responseMsg.status = status

	return &responseMsg
}

// User represents a chat user
// In and Out channles are used to send/received messages
// to/from remote user
// Id is transport specific id whicle UserId is authentication
// user id. UserId is used for sending/receiving messages
type User struct {
	Id     string
	UserId string
	In     chan *Message
	Out    chan *Message
	Done   chan interface{}
}

func NewUser(id string) *User {
	return &User{
		Id:   id,
		In:   make(chan *Message),
		Out:  make(chan *Message),
		Done: make(chan interface{}),
	}
}

func (u *User) Disconnect() {
	close(u.Done)
}
