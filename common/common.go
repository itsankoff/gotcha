package common

import "log"

const (
    TEXT   = 1
    BINARY = 2
)

type DataType int

type Message struct {
    from *User
    to   *User

    // Possible values:
    //  control
    //  message
    //  file
    messageType string
    dataType DataType
    data []byte
}

func NewMessage(from *User, to *User, messageType string,
                dataType DataType, data []byte) Message {
    return Message{
        from: from,
        to: to,
        messageType: messageType,
        dataType: dataType,
        data: data,
    }
}

func (m Message) String() string {
    if m.dataType != TEXT {
        log.Println("Message data is not a text type")
        return ""
    }

    return string(m.data)
}

func (m Message) Binary() []byte {
    if m.dataType != BINARY {
        log.Println("Message data is not binary type")
        return []byte{}
    }

    return m.data
}

func (m Message) From() *User {
    return m.from
}

func (m Message) To() *User {
    return m.to
}

func (m Message) DataType() DataType {
    return m.dataType
}

func (m Message) Raw() []byte {
    return m.data
}

type User struct {
    Id  string
    In      chan Message
    Out     chan Message
}
