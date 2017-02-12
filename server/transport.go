package server

import "github.com/itsankoff/gotcha/util"

type UserHandler func(*util.User)

type Transport interface {
    OnUserConnected(UserHandler)
    OnUserDisconnected(UserHandler)
    Start(url string, done <-chan interface{})
}
