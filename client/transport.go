package client

import "github.com/itsankoff/gotcha/common"

type Transport interface {
	Connect(host string) (bool, error)
	ConnectAsync(host string) chan bool
	Disconnect()
	Reconnect() (bool, error)
	ReconnectAsync() chan bool
	SendText(message string)
	SendBinary(date []byte)
	SetReceiver(<-chan common.Message)
}
