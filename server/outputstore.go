package server

import (
	"errors"
	"github.com/itsankoff/gotcha/common"
	"log"
)

// Store for all authenticated outputs
// Main facility to send a message to the remote user
type OutputStore struct {
	outputs map[string]chan<- *common.Message
}

func NewOutputStore() *OutputStore {
	out := &OutputStore{
		outputs: make(map[string]chan<- *common.Message),
	}

	return out
}

// AddOutput adds user output in the store by id
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

// RemoveOutput removes user output from the store by id
func (store *OutputStore) RemoveOutput(id string) error {
	_, ok := store.outputs[id]
	if ok {
		delete(store.outputs, id)
		log.Println("Remove output from store", id)
		return nil
	}

	return errors.New("Failed to remove output from store " + id)
}

// GetOutput returns a user output from the store by id
func (store OutputStore) GetOutput(id string) chan<- *common.Message {
	output, ok := store.outputs[id]
	if ok {
		return output
	}

	log.Println("Failed to get output for id", id)
	return nil
}

// Send sends a message to user's output
func (store OutputStore) Send(msg *common.Message) {
	output := store.GetOutput(msg.To())
	if output != nil {
		output <- msg
	} else {
		log.Println("Failed to send message", msg)
	}
}
