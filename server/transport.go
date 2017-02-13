package server

import "github.com/itsankoff/gotcha/util"

type UserHandler func(*util.User)

type Transport interface {
    OnUserConnected(chan<- *util.User)
    OnUserDisconnected(chan<- *util.User)
    Start(host string, done <-chan interface{})
}
