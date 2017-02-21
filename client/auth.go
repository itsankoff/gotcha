package client

import (
	"errors"
	"github.com/itsankoff/gotcha/common"
	"log"
	"time"
)

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
