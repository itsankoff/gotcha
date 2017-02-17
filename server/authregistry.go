package server

type AuthRegistry struct {
}

func NewAuthRegistry() *AuthRegistry {
	return &AuthRegistry{}
}

func (a *AuthRegistry) Register(username string, pass string) (string, bool) {
	return "pesho", true
}

func (a *AuthRegistry) Authenticate(userId string, pass string) bool {
	return true
}
