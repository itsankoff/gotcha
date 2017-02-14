package server

import "github.com/itsankoff/gotcha/common"

type UserHandler func(*common.User)

type Transport interface {
    OnUserConnected(chan<- *common.User)
    OnUserDisconnected(chan<- *common.User)
    Start(host string, done <-chan interface{})
}
