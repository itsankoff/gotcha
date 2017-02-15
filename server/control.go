package server

import (
    "github.com/itsankoff/gotcha/common"
    "log"
)

type Control struct {
    input chan *common.Message
}

func NewControl(input chan *common.Message) *Control {
    c := &Control{
        input: input,
    }

    go c.listen()
    return c
}

func (c Control) listen() {
    select {
    case msg := <-c.input:
        log.Println("Control received", msg)
        valid := c.validate(msg)
        if valid {
            cmd := msg.Cmd()
            switch(cmd) {
            case "register":
            case "auth":
            case "list_contacts":
            case "add_contact":
            case "remove_contact":
            case "create_group":
            case "add_to_group":
            case "remove_from_group":
            case "delete_group":
            case "list_groups":
            case "join_group":
            case "leave_group":
            default:
                log.Println("Unknown control command", cmd)
            }
        } else {
            log.Println("Invalid control message", msg)
        }
    }
}

func (c Control) validate(msg *common.Message) bool {
    return true
}

func (c Control) RegisterUser(user string, password string) *common.User {
    return nil
}

func (c Control) AuthUser(user string, password string) bool {
    return false
}

func (c Control) ListContacts(user string) []*common.User {
    return nil
}

func (c Control) AddContact(user *common.User, contact *common.User) bool {
    return false
}

func (c Control) RemoveContact(user *common.User, contact *common.User) bool {
    return false
}

func (c Control) CreateGroup(user *common.User) (groupId string) {
    return
}

func (c Control) AddToGroup(user *common.User, groupId string) bool {
    return false
}

func (c Control) RemoveFromGroup(user *common.User, groupId string) bool {
    return false
}

func (c Control) DeleteGroup(groupId string) bool {
    return false
}

func (c Control) ListGroups(user *common.User) *[]string {
    return &[]string{}
}

func (c Control) JoinGroup(user *common.User, groupId string) bool {
    return false
}

func (c Control) LeaveGroup(user *common.User, groupId string) bool {
    return false
}
