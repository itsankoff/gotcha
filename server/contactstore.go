package server

type ContactStore struct {
}

func NewContactStore() *ContactStore {
	return &ContactStore{}
}

func (c *ContactStore) AddContact(userId string, contactId string) bool {
	return false
}

func (c *ContactStore) RemoveContact(userId string, contactId string) bool {
	return false
}

func (c *ContactStore) ListContacts(userId string) ([]string, bool) {
	return []string{}, true
}
