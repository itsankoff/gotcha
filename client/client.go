package client

import (
	"errors"
	"time"
)

type AsyncHandler interface {
	Connected(host string, connected bool, err error)
}

type Client struct {
	transport Transport
	contacts  []string
}

type ClientHandler interface {
	OnMessage(userId int64, message string)
	OnTempMessage(userId int64, message string, expire time.Duration)
	OnGroupMessage(userId int64, message string)
	OnGroupTempMessage(groupId int64, message string, expire time.Duration)
	OnFile(userId int64, file File)
	OnGroupFile(userId int64, file File)
}

func New(transport Transport) *Client {
	return &Client{
		transport: transport,
	}
}

func (c *Client) Connect(host string) (bool, error) {
	return c.transport.Connect(host)
}

func (c *Client) ConnectAsync(host string) chan bool {
	return c.transport.ConnectAsync(host)
}

func (c *Client) Disconnect() {
	c.transport.Disconnect()
}

func (c *Client) Reconnect() (bool, error) {
	return c.transport.Reconnect()
}

func (c *Client) ReconnectAsync() chan bool {
	return c.transport.ReconnectAsync()
}

func (c *Client) Register(username, password string) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) ListContacts() ([]string, error) {
	return []string{}, errors.New("Not Implemented")
}

func (c *Client) SearchContact(username string) (int64, error) {
	return -1, errors.New("Not Implemented")
}

func (c *Client) AddContact(int64) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) RemoveContact(int64) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) CreateGroup() int64 {
	return -1
}

func (c *Client) AddToGroup(groupId, userId int64) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) RemoveFromGroup(groupId, userId int64) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) DeleteGroup(groupId int64) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) ListGroups() ([]*Group, error) {
	return []*Group{}, errors.New("Not Implemented")
}

func (c *Client) JoinGroup(groupId int64) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) LeaveGroup(groupId int64) {

}

func (c *Client) SendMessage(userId int64, message string) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) SendTempMessage(userId int64, message string, expire time.Duration) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) SendGroupMessage(groupId int64, message string) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) SendFile(userId int64, filePath string) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) SendGroupFile(userId int64, filePath string) (bool, error) {
	return false, errors.New("Not Implemented")
}

func (c *Client) GetHistory(from time.Time, to time.Time) (History, error) {
	return History{}, errors.New("Not Implemented")
}

func (c *Client) SetHandler(handler ClientHandler) {

}

func (c *Client) PrintHelp() {

}

func (c *Client) StartInteractiveMode() {

}
