package server

import (
    "github.com/itsankoff/gotcha/common"
    "time"
)

type History struct {

}

func NewHistory() *History {
    return &History{}
}

func (h History) AddMessage(m *common.Message) bool {
    return false
}

func (h History) GetHistory(user *common.User,
                            from time.Time,
                            to time.Time) []*common.Message {
    return []*common.Message{}
}
