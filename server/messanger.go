package server

import (
    "github.com/itsankoff/gotcha/common"
)

type Messanger struct {

}

func (m Messanger) SendMessage(msg *common.Message) bool {
    return false
}

func (m Messanger) SendTmpMessage(msg *common.Message) bool {
    return false
}

func (m Messanger) SendFile(msg *common.Message) bool {
    return false
}

func (m Messanger) SendGroupFile(msg *common.Message) bool {
    return false
}

func (m Messanger) SendTmpFile(msg *common.Message) bool {
    return false
}
