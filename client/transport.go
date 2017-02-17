package client

import "github.com/itsankoff/gotcha/common"

type Transport interface {
	Connect(host string) error
	ConnectAsync(host string) chan bool
	Disconnect()
	Reconnect() error
	ReconnectAsync() chan bool
	SendText(message string) error
	SendBinary(date []byte) error
	SetReceiver(chan *common.Message)
}
