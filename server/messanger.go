package server

import (
	"github.com/itsankoff/gotcha/common"
	"log"
)

type Messanger struct {
	input       chan *common.Message
	history     *History
	fileStore   *FileStore
	outputStore *OutputStore
}

func NewMessanger(input chan *common.Message,
	history *History, outputStore *OutputStore, fileStore *FileStore) *Messanger {
	m := &Messanger{
		input:       input,
		fileStore:   fileStore,
		history:     history,
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
				if msg.CmdType() == "file" {
					var fileContent string
					if msg.DataType() == common.TEXT {
						fileContent = msg.String()
					}

					uri := m.fileStore.AddFile(fileContent)
					newMsg := common.NewMessage(msg.From(), msg.To(),
						msg.CmdType(), msg.Cmd(),
						msg.ExpireDate(), common.TEXT, uri)
					msg = &newMsg
				}

				if msg.ExpireDate().IsZero() {
					m.history.AddMessage(msg)
				} else {
					log.Println("Temporary Message")
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
