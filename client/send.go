package client

import (
	"errors"
	"github.com/itsankoff/gotcha/common"
	"io/ioutil"
	"log"
	"time"
)

// SendMessage sends a instant message to a user or a group
func (c *Client) SendMessage(userId string, message string) error {
	return c.SendTempMessage(userId, message, time.Time{})
}

// SendTempMessage sends a temporary message
// with expire period to a user or a group
func (c *Client) SendTempMessage(userId string, message string,
	expire time.Time) error {
	if !c.authenticated {
		return errors.New("Not authenticated. Call Authenticate first")
	}

	msg := common.NewMessage(c.userId, userId,
		"message", "send_message",
		expire, common.TEXT, message)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode instant message", err)
		return err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send instant message", err)
		return err
	}

	log.Println("Instant message sent")
	return nil
}

// SendTextFile sends a text file to a user or a group
func (c *Client) SendTextFile(userId string, filePath string) (string, error) {
	var link string
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("Failed to read file", err)
		return link, err
	}

	msg := common.NewMessage(c.userId, userId,
		"file", "send_file", time.Time{},
		common.TEXT, string(fileContent))

	data, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode file", err)
		return link, err
	}

	err = c.transport.SendText(string(data))
	if err != nil {
		log.Println("Failed to send file content", err)
		return link, err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("Send file response error", errMsg)
		return link, errors.New(errMsg)
	}

	link = resp.String()
	log.Println("File sent", link)
	return link, nil
}
