package client

import (
	"bufio"
	"fmt"
	"github.com/itsankoff/gotcha/common"
	"os"
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
