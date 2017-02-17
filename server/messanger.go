package server

import (
    "github.com/itsankoff/gotcha/common"
    "log"
)

type Messanger struct {
    input           chan *common.Message
    history         *History
    outputStore     *OutputStore
}

func NewMessanger(input chan *common.Message,
                  history *History,
                  outputStore *OutputStore) *Messanger {
    m := &Messanger {
        input: input,
        history: history,
        outputStore: outputStore,
    }

    go m.listen()
    return m
}

func (m *Messanger) listen() {
    for {
        select {
        case msg := <-m.input:
            log.Println("Message received")
            valid := m.validate(msg)
            if valid {
                if msg.Cmd() == "file" {
                    // add file to file store
                    // change msg content to file url (plus token)
                }

                if !(msg.ExpireDate().IsZero()) {
                    m.history.AddMessage(msg)
                } else {
                    log.Println("Message expire date is not zero", msg.ExpireDate())
                }

                m.outputStore.Send(msg)
            } else {
                log.Println("Invalid instant message", msg)
            }
        }
    }
}

func (m Messanger) validate(msg *common.Message) bool {
    return true
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
