package server

import (
	"github.com/itsankoff/gotcha/common"
	"log"
	"time"
)

type History struct {
	input       chan *common.Message
	outputStore *OutputStore
}

func NewHistory(input chan *common.Message, outputStore *OutputStore) *History {
	history := &History{
		input:       input,
		outputStore: outputStore,
	}

	go history.listen()
	return history
}

func (h History) listen() {
	for {
		select {
		case msg := <-h.input:
			valid := h.validate(msg)
			if !valid {
				log.Println("Failed to validate history message")
				response := common.NewResponse(msg, "error", "invalid message")
				h.outputStore.Send(response)
				return
			}

			payload, err := msg.ParseJsonData()
			if err != nil {
				log.Printf("Failed to parse history message from %s", msg.From())
				response := common.NewResponse(msg, "error", "bad message")
				h.outputStore.Send(response)
				return
			}

			cmd := msg.Cmd()
			switch cmd {
			case "get_history":
				userId := msg.From()
				forUserId := payload["for_user_id"].(string)
				fromDate := time.Unix(int64(payload["from_date"].(float64)), 0)
				toDate := time.Unix(int64(payload["to_date"].(float64)), 0)
				messages := h.GetHistory(userId, forUserId, fromDate, toDate)
				accumulated := h.accumulate(messages)
				response := common.NewResponse(msg, "history", accumulated)
				h.outputStore.Send(response)
			}
		}
	}
}

func (h History) validate(msg *common.Message) bool {
	return true
}

func (h History) accumulate(messages []*common.Message) string {
	return ""
}

func (h History) AddMessage(msg *common.Message) bool {
	log.Println("Add message in history", msg)
	return false
}

func (h History) GetHistory(userId string, forUserId string,
	from time.Time, to time.Time) []*common.Message {
	return []*common.Message{}
}
