package server

import (
	"log"
	"strconv"
	"time"
)

// Stores authenticated users in memory
type AuthRegistry struct {
	users map[string][]string
}

func NewAuthRegistry() *AuthRegistry {
	return &AuthRegistry{
		users: make(map[string][]string),
	}
}

// Register generates server domain user id which
// for now is current ts in nanoseconds
func (a *AuthRegistry) Register(username string, pass string) (string, bool) {
	now := time.Now().UnixNano()
	userId := strconv.FormatInt(now, 10)
	a.users[userId] = []string{username, pass}
	log.Printf("User %s registered with id %s", username, userId)
	return userId, true
}

// Authenticates authenticates the user id with it's password
func (a *AuthRegistry) Authenticate(userId string, pass string) bool {
	credentials, exists := a.users[userId]
	if !exists {
		log.Println("No user with this id", userId)
		return false
	}

	if credentials[1] != pass {
		log.Println("Failed to autenticate user", userId)
		return false
	}

	log.Printf("User %s authenticated", userId)
	return true
}

// SearchContact searches globally for contact with param username
func (a *AuthRegistry) SearchContact(username string) string {
	for uId, c := range a.users {
		if c[0] == username {
			return uId
		}
	}

	return ""
}
