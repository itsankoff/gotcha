package user

// User represents a chat user.
// In and Out channels are used to send/received messages to/from a remote user.
// Id is transport specific id while UserId is authentication user id.
// UserId is used for sending/receiving messages.
type User struct {
	Id     string
	UserId string
	In     chan *Message
	Out    chan *Message
	Done   chan struct{}
}

// New creates a new user with given id
func New(id string) *User {
	return &User{
		Id:   id,
		In:   make(chan *Message),
		Out:  make(chan *Message),
		Done: make(chan struct{}),
	}
}

// Disconnect closes user's done channel
func (u *User) Disconnect() {
	close(u.Done)
}
