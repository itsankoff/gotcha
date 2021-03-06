package server

import (
	"github.com/itsankoff/gotcha/common"
	"log"
)

// Handler for all messages/files communication
// between the clients.
type Messanger struct {
	input       chan *common.Message
	history     *History
	fileStore   *FileStore
	outputStore *OutputStore
}

// NewMessenger creates a new messagener object
// which handlers messages through input channel.
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
					var uri string
					if msg.DataType() == common.TEXT {
						fileContent := msg.String()
						uri = m.fileStore.AddTextFile(fileContent)
					}

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
