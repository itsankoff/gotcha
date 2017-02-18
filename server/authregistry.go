package server

import "log"

type AuthRegistry struct {
}

func NewAuthRegistry() *AuthRegistry {
	return &AuthRegistry{}
}

func (a *AuthRegistry) Register(username string, pass string) (string, bool) {
	userId := username
	log.Printf("User %s registered with id %s", username, userId)
	return userId, true
}

func (a *AuthRegistry) Authenticate(userId string, pass string) bool {
	log.Printf("User %s authenticated", userId)
	return true
}
