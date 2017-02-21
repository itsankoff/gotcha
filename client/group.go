package client

import (
	"errors"
	"github.com/itsankoff/gotcha/common"
	"log"
	"time"
)

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
