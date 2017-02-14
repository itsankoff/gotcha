package server

import (
    "github.com/itsankoff/gotcha/common"
)

type Control struct {

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
