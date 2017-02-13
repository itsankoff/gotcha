package util

import "errors"

const (
    BINARY = iota
    TEXT
)

type DataType int

type Message struct {
    from *User
    to   *User
    dataType DataType
    data []byte
}

func NewMessage(from *User, to *User, dataType DataType, data []byte) Message {
    return Message{
        from: from,
        to: to,
        dataType: dataType,
        data: data,
    }
}

func (m Message) String() (string, error) {
    if m.dataType != TEXT {
        return "", errors.New("Message data is not a text type")
    }

    return string(m.data), nil
}

func (m Message) Binary() ([]byte, error) {
    if m.dataType != BINARY {
        return []byte{}, errors.New("Message data is not binary type")
    }

    return m.data, nil
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

type User struct {
    Id  string
    In      chan Message
    Out     chan Message
}
