package server

type ContactStore struct {
}

func NewContactStore() *ContactStore {
	return &ContactStore{}
}

func (c *ContactStore) AddContact(userId string, contactId string) bool {
	return true
}

func (c *ContactStore) RemoveContact(userId string, contactId string) bool {
	return true
}

func (c *ContactStore) ListContacts(userId string) ([]string, bool) {
	return []string{}, true
}

func (c *ContactStore) SearchContact(userId string, contactName string) string {
	return ""
}
