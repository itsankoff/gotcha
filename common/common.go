package common

import (
    "log"
    "time"
)

const (
    TEXT   = 1
    BINARY = 2
)

type DataType int

type Message struct {
    from     string
    to       string

    cmdType  string
    cmd      string

    expireDate time.Time

    dataType DataType
    data interface{}
}

func NewMessage(from string, to string,
                cmdType string, cmd string,
                expireDate time.Time,
                dataType DataType, data interface{}) Message {
    return Message{
        from: from,
        to: to,
        cmdType: cmdType,
        cmd: cmd,
        expireDate: expireDate,
        dataType: dataType,
        data: data,
    }
}

func (m Message) String() string {
    if m.dataType != TEXT {
        log.Println("Message data is not a text type")
        return ""
    }

    return m.data.(string)
}

func (m Message) Binary() []byte {
    if m.dataType != BINARY {
        log.Println("Message data is not binary type")
        return []byte{}
    }

    return m.data.([]byte)
}

func (m Message) From() string {
    return m.from
}

func (m Message) To() string {
    return m.to
}

func (m Message) CmdType() string {
    return m.cmdType
}

func (m Message) Cmd() string {
    return m.cmd
}

func (m Message) ExpireDate() time.Time {
    return m.expireDate
}

func (m Message) DataType() DataType {
    return m.dataType
}

type User struct {
    Id  string
    In      chan *Message
    Out     chan *Message
}
