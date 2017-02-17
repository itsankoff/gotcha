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

func (m Message) String() string {
	if m.dataType != TEXT {
		log.Println("Message data is not a text type")
		return ""
	}

	return m.data.(string)
}

func (m Message) Binary() []byte {
	if m.dataType != BINARY {
		log.Println("Message data is not binary type")
		return []byte{}
	}

	return m.data.([]byte)
}

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

func (m Message) ParseJsonData() (map[string]interface{}, error) {
	return m.data.(map[string]interface{}), nil
}

func (m Message) From() string {
	return m.from
}

func (m Message) To() string {
	return m.to
}

func (m Message) CmdType() string {
	return m.cmdType
}

func (m Message) Cmd() string {
	return m.cmd
}

func (m Message) Status() Status {
	return m.status
}

func (m Message) ExpireDate() time.Time {
	return m.expireDate
}

func (m Message) DataType() DataType {
	return m.dataType
}

func NewResponse(msg *Message, response interface{}) *Message {
	responseMsg := NewMessage("server", msg.from, msg.cmdType, msg.cmd,
		time.Time{}, TEXT, response)
	return &responseMsg
}

type User struct {
	Id   string
	In   chan *Message
	Out  chan *Message
	Done chan interface{}
}

func NewUser(userId string) *User {
	return &User{
		Id:   userId,
		In:   make(chan *Message),
		Out:  make(chan *Message),
		Done: make(chan interface{}),
	}
}

func (u *User) Disconnect() {
	close(u.Done)
}
