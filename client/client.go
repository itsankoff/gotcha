package client

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/itsankoff/gotcha/common"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Main object which represents the chat
// client
type Client struct {
	Out           chan *common.Message
	transport     Transport
	contacts      []string
	userId        string
	username      string
	password      string
	authenticated bool
}

// New creates a new chat client and binds the underlying
// transport for sending/receiving messages and files
func New(transport Transport) *Client {
	client := &Client{
		Out:       make(chan *common.Message),
		transport: transport,
	}

	client.transport.SetReceiver(client.Out)
	return client
}

// Connect calls undelying transport Connect method
func (c *Client) Connect(host string) error {
	return c.transport.Connect(host)
}

// ConnectAsync calls undelying transport ConnectAsync method
func (c *Client) ConnectAsync(host string) chan bool {
	return c.transport.ConnectAsync(host)
}

// Disconnect calls undelying transport Disconnect method
func (c *Client) Disconnect() {
	c.transport.Disconnect()
}

// Reconnect calls undelying transport Reconnect method
func (c *Client) Reconnect() error {
	return c.transport.Reconnect()
}

// ReconnectAsync calls undelying transport ReconnectAsync method
func (c *Client) ReconnectAsync() chan bool {
	return c.transport.ReconnectAsync()
}

// Register registers the client with username and passowrd.
// If registration is successful server specific userId is return
// userId is used with the Athenticate method later
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
			errorMessage := resp.Error()
			log.Println("Failed to register", errorMessage)
			return "", errors.New(errorMessage)
		}

		userId := resp.GetJsonData("user_id").(string)
		log.Println("User registered", userId)
		return userId, nil
	case <-time.After(time.Second * 10):
		log.Println("Register response timeout")
		return "", errors.New("Register response timeout")
	}
}

// Authenticate authenticates the user with is password.
// userId is obtain from Register methods earlier.
// If authentication is successful returns nil error
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
			errorMessage := resp.Error()
			log.Println("Failed to authenticate user", errorMessage)
			return errors.New(errorMessage)
		}

		c.userId = userId
		c.authenticated = true
		authenticated := resp.GetJsonData("authenticated").(bool)
		log.Println("User authenticated", authenticated)
		return nil
	case <-time.After(10 * time.Second):
		log.Println("Authentication response timeout")
		return errors.New("Authentication response timeout")
	}
}

// UserId return the server specific userId
func (c Client) UserId() string {
	return c.userId
}

// ListContacts list the added contacts for this client
func (c *Client) ListContacts() ([]string, error) {
	var contacts []string
	if !c.authenticated {
		return contacts, errors.New("Not authenticated. Call Authenticate first")
	}

	var payload map[string]interface{}
	msg := common.NewMessage(c.userId, "server",
		"control", "list_contacts", time.Time{},
		common.TEXT, payload)

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode auth message", err)
		return contacts, err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send auth message", err)
		return contacts, err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("List contacts response error", errMsg)
		return contacts, errors.New(errMsg)
	}

	contactsData := resp.GetJsonData("contacts")
	rawContacts, ok := contactsData.([]interface{})
	if !ok {
		return contacts, errors.New("Failed to parse contacts response")
	}

	for _, rawContact := range rawContacts {
		contact, ok := rawContact.(string)
		if !ok {
			log.Println("Failed to parse contact", rawContact)
			continue
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// SearchContact searches for user with contactName globally on the server
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
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("Search contact response error", errMsg)
		return "", errors.New(errMsg)
	}

	contactId := resp.GetJsonData("contact_id").(string)
	return contactId, nil
}

// AddContact adds a contact for this client
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
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("Add contact response error", errMsg)
		return errors.New(errMsg)
	}

	return nil
}

// Remove contact removes a contact for this client
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
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("Remove contact response error", errMsg)
		return errors.New(errMsg)
	}

	return nil
}

// CreateGroup creates a chat group
func (c *Client) CreateGroup() (string, error) {
	var groupId string
	if !c.authenticated {
		return groupId, errors.New("Not authenticated. Call Authenticate first")
	}

	msg := common.NewMessage(c.userId, "server",
		"control", "remove_contact", time.Time{},
		common.TEXT, "")

	encoded, err := msg.Json()
	if err != nil {
		log.Println("Failed to encode create group message", err)
		return groupId, err
	}

	err = c.transport.SendText(string(encoded))
	if err != nil {
		log.Println("Failed to send create group message", err)
		return groupId, err
	}

	resp := <-c.Out
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("Create group response error", errMsg)
		return groupId, errors.New(errMsg)
	}

	groupId = resp.GetJsonData("group_id").(string)
	return groupId, nil
}

// AddToGroup adds a participant in already created group
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
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("Add to group response error", errMsg)
		return errors.New(errMsg)
	}

	return nil
}

// RemoveFromGroup removes a participant for already created group
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
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("Remove from group response error", errMsg)
		return errors.New(errMsg)
	}

	return nil
}

// DeleteGroup deletes a group
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
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("Delete group response error", errMsg)
		return errors.New(errMsg)
	}

	return nil
}

// ListGroups lists all available groups
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
	if resp.Status() == common.STATUS_ERROR {
		errMsg := resp.Error()
		log.Println("Delete group response error", errMsg)
		return groups, errors.New(errMsg)
	}

	groups, _ = resp.GetJsonData("groups").([]string)
	return groups, nil
}

// JoinGroup adds the client to the group
func (c *Client) JoinGroup(groupId string) error {
	return c.AddToGroup(c.userId, groupId)
}

// LeaveGroup removes the client from the group
func (c *Client) LeaveGroup(groupId string) error {
	return c.RemoveFromGroup(c.userId, groupId)
}

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

// GetHistory retrieves the conversation history for a user
// or a group
func (c *Client) GetHistory(remote string, from time.Time, to time.Time) (History, error) {
	return History{}, errors.New("Not Implemented")
}

// StartInteractiveMode starts a commandline mode where you chat
// with the other users in the system
func (c *Client) StartInteractiveMode() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		c.SendMessage(c.userId, text)
		resp := <-c.Out
		data := resp.String()
		fmt.Println("Response: ", data)
	}
}
