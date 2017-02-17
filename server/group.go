package server

import (
    "time"
    "github.com/itsankoff/gotcha/common"
    "strconv"
)

type Group struct {
    Id              string
    Out             chan *common.Message
    outputs         map[string]chan<- *common.Message
}

func NewGroup() *Group {
    group := &Group{
        Id: strconv.FormatInt(time.Now().Unix(), 10),
        outputs: make(map[string]chan<- *common.Message),
    }

    go group.listen()
    return group
}

func (g *Group) listen() {
    select {
    case msg := <-g.Out:
        for _, out := range g.outputs {
            out<- msg
        }
    }
}

func (g *Group) AddOutput(id string, output chan<- *common.Message) bool{
    g.outputs[id] = output
    return true
}

func (g *Group) RemoveOutput(id string) bool {
    _, ok := g.outputs[id]
    if ok {
        delete(g.outputs, id)
        return true
    }

    return false
}