package server

import "log"

// Stores contact lists for users
type ContactStore struct {
	contacts map[string][]string
}

func NewContactStore() *ContactStore {
	return &ContactStore{
		contacts: make(map[string][]string),
	}
}

// AddContact adds contact in user's contact list
func (c *ContactStore) AddContact(userId string, contactId string) bool {
	userContacts, ok := c.contacts[userId]
	if !ok {
		c.contacts[userId] = []string{}
		userContacts = c.contacts[userId]
	}

	userContacts = append(userContacts, contactId)
	c.contacts[userId] = userContacts
	return true
}

// RemoveContact removes contact from user's contact list
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

// ListContacts returns user's contact list for userId
func (c *ContactStore) ListContacts(userId string) ([]string, bool) {
	userContacts, ok := c.contacts[userId]
	if !ok {
		return []string{}, false
	}

	return userContacts, true
}
