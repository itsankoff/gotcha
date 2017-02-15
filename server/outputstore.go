package server

import (
    "github.com/itsankoff/gotcha/common"
    "log"
    "errors"
)


type command struct {
    cmd         string
    outputId    string
}

type OutputStore struct {
    manager chan *command
    outputs map[string]chan<- *common.Message
}

func NewOutputStore() *OutputStore {
    out := &OutputStore{
        manager: make(chan *command),
        outputs: make(map[string]chan<- *common.Message),
    }

    return out
}

func (store *OutputStore) AddOutput(id string,
                                    output chan<- *common.Message) error {
    _, ok := store.outputs[id]
    if ok {
        return errors.New("Output for already exists for id " + id)
    }

    store.outputs[id] = output
    log.Println("Add output in store", id)
    return nil
}

func (store *OutputStore) RemoveOutput(id string) error {
    _, ok := store.outputs[id]
    if ok {
        delete(store.outputs, id)
        return nil
    }

    return errors.New("Failed to remove output from store " + id)
}

func (store OutputStore) GetOutput(id string) chan<- *common.Message {
    output, ok := store.outputs[id]
    if ok {
        return output
    }

    log.Println("Failed to get output for id", id)
    return nil
}

func (store OutputStore) Send(msg *common.Message) {
    output := store.GetOutput(msg.To())
    if output != nil {
        output <- msg
    } else {
        log.Println("Failed to send message", msg)
    }
}
