package client

import (
	"encoding/json"
	"errors"
	"github.com/itsankoff/gotcha/common"
	"log"
	"time"
)

type Client struct {
	Out           chan *common.Message
	transport     Transport
	contacts      []string
	userId        string
	username      string
	password      string
	authenticated bool
}

func New(transport Transport) *Client {
	client := &Client{
		Out:       make(chan *common.Message),
		transport: transport,
	}

	client.transport.SetReceiver(client.Out)
	return client
}

func (c *Client) Connect(host string) error {
	return c.transport.Connect(host)
}

func (c *Client) ConnectAsync(host string) chan bool {
	return c.transport.ConnectAsync(host)
}

func (c *Client) Disconnect() {
	c.transport.Disconnect()
}

func (c *Client) Reconnect() error {
	return c.transport.Reconnect()
}

func (c *Client) ReconnectAsync() chan bool {
	return c.transport.ReconnectAsync()
}

func (c *Client) Register(username, password string) (string, error) {
	c.username = username
	c.password = password

	payload := make(map[string]interface{})
	payload["username"] = username
	payload["password"] = password
	msg := common.NewMessage(username, "server",
		"auth", "register", time.Time{},
		common.TEXT, payload)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode register message", err)
		return "", err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send register message", err)
		return "", err
	}

	select {
	case resp := <-c.Out:
		if resp.Status() == common.STATUS_ERROR {
			errorMessage := resp.String()
			log.Println("Failed to register", errorMessage)
			return "", errors.New(errorMessage)
		}

		userId := resp.String()
		log.Println("User registered", userId)
		return userId, nil
	case <-time.After(time.Second * 10):
		log.Println("Register response timeout")
		return "", errors.New("Register response timeout")
	}
}

func (c *Client) Authenticate(userId, password string) error {
	payload := make(map[string]interface{})
	payload["user_id"] = userId
	payload["password"] = password
	msg := common.NewMessage(userId, "server",
		"auth", "auth", time.Time{},
		common.TEXT, payload)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode auth message", err)
		return err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send auth message", err)
		return err
	}

	select {
	case resp := <-c.Out:
		if resp.Status() == common.STATUS_ERROR {
			errorMessage := resp.String()
			log.Println("Failed to authenticate user", errorMessage)
			return errors.New(errorMessage)
		}

		respMessage := resp.String()
		c.userId = userId
		c.authenticated = true
		log.Println("User authenticated", respMessage)
		return nil
	case <-time.After(10 * time.Second):
		log.Println("Authentication response timeout")
		return errors.New("Authentication response timeout")
	}
}

func (c *Client) ListContacts() ([]string, error) {
	if !c.authenticated {
		return []string{}, errors.New("Not authenticated. Call Authenticate first")
	}

	msg := common.NewMessage(c.userId, "server",
		"control", "list_contacts", time.Time{},
		common.TEXT, "")

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode auth message", err)
		return []string{}, err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send auth message", err)
		return []string{}, err
	}

	resp := <-c.Out
	data := resp.String()
	var contacts []string
	err = json.Unmarshal([]byte(data), &contacts)
	if err != nil {
		log.Println("Failed to parse contact response", err)
		return []string{}, err
	}

	return contacts, nil
}

func (c *Client) SearchContact(contactName string) (string, error) {
	if !c.authenticated {
		return "", errors.New("Not authenticated. Call Authenticate first")
	}

	payload := make(map[string]interface{})
	payload["contact_name"] = contactName
	msg := common.NewMessage(c.userId, "server",
		"control", "seach_contact", time.Time{},
		common.TEXT, payload)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode search contact message", err)
		return "", err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send search contact message", err)
		return "", err
	}

	resp := <-c.Out
	contactId := resp.String()
	return contactId, nil
}

func (c *Client) AddContact(contactId string) error {
	if !c.authenticated {
		return errors.New("Not authenticated. Call Authenticate first")
	}

	payload := make(map[string]interface{})
	payload["contact_id"] = contactId
	msg := common.NewMessage(c.userId, "server",
		"control", "add_contact", time.Time{},
		common.TEXT, payload)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode add contact message", err)
		return err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send add contact message", err)
		return err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_OK {
		return nil
	}

	return errors.New(resp.String())
}

func (c *Client) RemoveContact(contactId string) error {
	if !c.authenticated {
		return errors.New("Not authenticated. Call Authenticate first")
	}

	payload := make(map[string]interface{})
	payload["contact_id"] = contactId
	msg := common.NewMessage(c.userId, "server",
		"control", "remove_contact", time.Time{},
		common.TEXT, payload)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode remove contact message", err)
		return err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send remove contact message", err)
		return err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_OK {
		return nil
	}

	return errors.New(resp.String())
}

func (c *Client) CreateGroup() (string, error) {
	if !c.authenticated {
		return "", errors.New("Not authenticated. Call Authenticate first")
	}

	msg := common.NewMessage(c.userId, "server",
		"control", "remove_contact", time.Time{},
		common.TEXT, "")

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode create group message", err)
		return "", err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send create group message", err)
		return "", err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_OK {
		return resp.String(), nil
	}

	return "", errors.New(resp.String())
}

func (c *Client) AddToGroup(userId, groupId string) error {
	if !c.authenticated {
		return errors.New("Not authenticated. Call Authenticate first")
	}

	payload := make(map[string]interface{})
	payload["group_id"] = groupId
	payload["user_id"] = userId
	msg := common.NewMessage(c.userId, "server",
		"control", "add_to_group", time.Time{},
		common.TEXT, payload)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode add to group message", err)
		return err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send add to group message", err)
		return err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_OK {
		return nil
	}

	return errors.New(resp.String())
}

func (c *Client) RemoveFromGroup(userId, groupId string) error {
	if !c.authenticated {
		return errors.New("Not authenticated. Call Authenticate first")
	}

	payload := make(map[string]interface{})
	payload["group_id"] = groupId
	payload["user_id"] = userId
	msg := common.NewMessage(c.userId, "server",
		"control", "remove_to_group", time.Time{},
		common.TEXT, payload)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode remove group message", err)
		return err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send remove from group message", err)
		return err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_OK {
		return nil
	}

	return errors.New(resp.String())
}

func (c *Client) DeleteGroup(groupId string) error {
	if !c.authenticated {
		return errors.New("Not authenticated. Call Authenticate first")
	}

	payload := make(map[string]interface{})
	payload["group_id"] = groupId
	msg := common.NewMessage(c.userId, "server",
		"control", "remove_to_group", time.Time{},
		common.TEXT, payload)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode delete group message", err)
		return err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send delete group message", err)
		return err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_OK {
		return nil
	}

	return errors.New(resp.String())
}

func (c *Client) ListGroups() ([]string, error) {
	var groups []string
	if !c.authenticated {
		return groups, errors.New("Not authenticated. Call Authenticate first")
	}

	msg := common.NewMessage(c.userId, "server",
		"control", "remove_to_group", time.Time{},
		common.TEXT, "")

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode list groups message", err)
		return groups, err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send list groups message", err)
		return groups, err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_OK {
		data := resp.String()
		err := json.Unmarshal([]byte(data), &groups)
		if err != nil {
			log.Println("Failed to parse list groups response", err)
			return []string{}, err
		}

		return groups, nil
	}

	return groups, errors.New(resp.String())
}

func (c *Client) JoinGroup(groupId string) error {
	return c.AddToGroup(c.userId, groupId)
}

func (c *Client) LeaveGroup(groupId string) error {
	return c.RemoveFromGroup(c.userId, groupId)
}

func (c *Client) SendMessage(userId int64, message string) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) SendTempMessage(userId int64, message string, expire time.Duration) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) SendFile(userId int64, filePath string) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) GetHistory(from time.Time, to time.Time) (History, error) {
	return History{}, errors.New("Not Implemented")
}

func (c *Client) PrintHelp() {

}

func (c *Client) StartInteractiveMode() {

}
