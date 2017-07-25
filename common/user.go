package common

// User represents a chat user
// In and Out channles are used to send/received messages
// to/from remote user
// Id is transport specific id whicle UserId is authentication
// user id. UserId is used for sending/receiving messages
type User struct {
	Id     string
	UserId string
	In     chan *Message
	Out    chan *Message
	Done   chan interface{}
}

func NewUser(id string) *User {
	return &User{
		Id:   id,
		In:   make(chan *Message),
		Out:  make(chan *Message),
		Done: make(chan interface{}),
	}
}

func (u *User) Disconnect() {
	close(u.Done)
}
