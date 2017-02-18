package server

import "github.com/itsankoff/gotcha/common"

type UserHandler func(*common.User)

// Defines the interface for the transport implementations
type Transport interface {
	OnUserConnected(chan<- *common.User)
	OnUserDisconnected(chan<- *common.User)
	Start(host string, done <-chan interface{})
}
