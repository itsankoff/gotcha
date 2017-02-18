package server

import (
	"github.com/itsankoff/gotcha/common"
	"strconv"
	"time"
)

// Aggregate output. When a message is received in the Out
// channel it will send it to all outputs in the group.
// Like a multicast channel
type Group struct {
	Id      string
	Out     chan *common.Message
	outputs map[string]chan<- *common.Message
}

func NewGroup() *Group {
	group := &Group{
		Id:      strconv.FormatInt(time.Now().Unix(), 10),
		Out:     make(chan *common.Message),
		outputs: make(map[string]chan<- *common.Message),
	}

	go group.listen()
	return group
}

func (g *Group) listen() {
	for {
		select {
		case msg := <-g.Out:
			for _, out := range g.outputs {
				out <- msg
			}
		}
	}
}

// AddOutput adds output in the aggregate
func (g *Group) AddOutput(id string, output chan<- *common.Message) bool {
	g.outputs[id] = output
	return true
}

// RemoveOutput removes output from the aggregate
func (g *Group) RemoveOutput(id string) bool {
	_, ok := g.outputs[id]
	if ok {
		delete(g.outputs, id)
		return true
	}

	return false
}
