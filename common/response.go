package common

import (
	"time"
)

// NewResponse is helper which generates responses for Message. Used only from
// server.
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
