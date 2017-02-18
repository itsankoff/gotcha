package server

import "log"

type ContactStore struct {
	contacts map[string][]string
}

func NewContactStore() *ContactStore {
	return &ContactStore{
		contacts: make(map[string][]string),
	}
}

func (c *ContactStore) AddContact(userId string, contactId string) bool {
	userContacts, ok := c.contacts[userId]
	if !ok {
		c.contacts[userId] = []string{}
		userContacts = c.contacts[userId]
	}

	userContacts = append(userContacts, contactId)
	return true
}

func (c *ContactStore) RemoveContact(userId string, contactId string) bool {
	userContacts, ok := c.contacts[userId]
	if !ok {
		log.Println("No contacts for user", userId)
		return false
	}

	for i, uc := range userContacts {
		if uc == contactId {
			userContacts = append(userContacts[:i], userContacts[i+1:]...)
			return true
		}
	}

	return false
}

func (c *ContactStore) ListContacts(userId string) ([]string, bool) {
	userContacts, ok := c.contacts[userId]
	if !ok {
		return []string{}, false
	}

	return userContacts, true
}
