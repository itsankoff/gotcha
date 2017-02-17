package server

import (
	"github.com/itsankoff/gotcha/common"
	"log"
	"time"
)

type History struct {
}

func NewHistory() *History {
	return &History{}
}

func (h History) AddMessage(msg *common.Message) bool {
	log.Println("Add message in history", msg)
	return false
}

func (h History) GetHistory(user *common.User,
	from time.Time,
	to time.Time) []*common.Message {
	return []*common.Message{}
}
