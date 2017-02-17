package server

import (
    "github.com/itsankoff/gotcha/common"
    "log"
    "encoding/json"
)

type Control struct {
    input           chan *common.Message
    groups          []*Group
    outputStore     *OutputStore
}

func NewControl(input chan *common.Message,
                outputStore *OutputStore) *Control {
    c := &Control{
        input: input,
        outputStore: outputStore,
    }

    go c.listen()
    return c
}

func (c Control) listen() {
    for {
        select {
        case msg := <-c.input:
            log.Println("Control received", msg)

            valid := c.validate(msg)
            if valid {
                var payload map[string]interface{}
                err := json.Unmarshal([]byte(msg.String()), &payload)
                if err == nil {
                    cmd := msg.Cmd()
                    switch(cmd) {
                    case "list_contacts":
                    case "add_contact":
                    case "remove_contact":
                    case "create_group":
                        groupId := c.CreateGroup()
                        c.AddToGroup(groupId, msg.From())
                        log.Println("Group created", groupId)
                    case "add_to_group":
                        groupId := payload["group_id"].(string)
                        userId := payload["user_id"].(string)
                        added := c.AddToGroup(groupId, userId)
                        log.Printf("User %s added to group %s %t",
                                   groupId, userId, added)
                    case "remove_from_group":
                        groupId := payload["group_id"].(string)
                        userId := payload["user_id"].(string)
                        removed := c.RemoveFromGroup(groupId, userId)
                        log.Printf("User %s removed from group %s %t",
                                   groupId, userId, removed)
                    case "delete_group":
                        groupId := payload["group_id"].(string)
                        deleted := c.DeleteGroup(groupId)
                        log.Printf("Group %s deleted %t", groupId, deleted)
                    case "list_groups":
                    case "join_group":
                    case "leave_group":
                    default:
                        log.Println("Unknown control command", cmd)
                    }
                } else {
                    log.Println("Failed to decode control message payload", msg)
                }
            } else {
                log.Println("Invalid control message", msg)
            }
        }
    }
}

func (c Control) validate(msg *common.Message) bool {
    return true
}

func (c Control) findGroup(groupId string) *Group {
    for _, g := range c.groups {
        if g.Id == groupId {
            return g
        }
    }

    return nil
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

func (c *Control) CreateGroup() string {
    group := NewGroup()
    c.groups = append(c.groups, group)
    c.outputStore.AddOutput(group.Id, group.Out)
    return group.Id
}

func (c Control) AddToGroup(groupId string, userId string) bool {
    group := c.findGroup(groupId)
    if group == nil {
        log.Println("Failed to add user. No group with id", groupId, userId)
        return false
    }

    userOutput := c.outputStore.GetOutput(userId)
    if userOutput == nil {
        log.Println("Failed to add user to group. No user output",
                    groupId, userId)
        return false
    }

    return group.AddOutput(userId, userOutput)
}

func (c Control) RemoveFromGroup(groupId string, userId string) bool {
    group := c.findGroup(groupId)
    if group == nil {
        log.Printf("Failed to remove user %s from group %s. No group with id",
                   userId, groupId)
        return false
    }

    return group.RemoveOutput(userId)
}

func (c *Control) DeleteGroup(groupId string) bool {
    for i, group:= range c.groups {
        if group.Id == groupId {
            c.groups = append(c.groups[:i], c.groups[i+1:]...)
            c.outputStore.RemoveOutput(group.Id)
            close(group.Out)
            return true
        }
    }

    return false
}

func (c Control) ListGroups(user *common.User) *[]string {
    groupIds := []string{}
    for _, g := range c.groups {
        groupIds = append(groupIds, g.Id)
    }

    return &groupIds
}
