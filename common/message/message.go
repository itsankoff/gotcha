package message

import (
	"time"
)

// Message represents protocol message. For more detailed information check
// protocol.txt
type Message struct {
	From       string      `json:"from"`
	To         string      `json:"to"`
	CmdType    string      `json:"cmd_type"`
	Cmd        string      `json:"cmd"`
	Status     Status      `json:"status"` // only for server -> client messages
	ExpiryDate time.Time   `json:"expiry_date"`
	DataType   DataType    `json:"data_type"`
	Data       interface{} `json:"data"`
}

// New creates a message from given arguments
func New(from, to, cmdType, cmd string, expireDate time.Time,
	dataType DataType, data interface{}) Message {

	return Message{
		From:       from,
		To:         to,
		CmdType:    cmdType,
		Cmd:        cmd,
		ExpiryDate: expireDate,
		DataType:   dataType,
		Data:       data,
	}
}

// String returns message payload in string
func (m *Message) String() string {
	return string(m.Data)
}

// Binary returns message payload in binary
func (m *Message) Binary() []byte {
	return []byte(m.Data)
}
