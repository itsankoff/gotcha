package client

import (
	"errors"
	"github.com/itsankoff/gotcha/common"
	"log"
	"time"
)

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
