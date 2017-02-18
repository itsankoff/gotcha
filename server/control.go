package server

import (
	"github.com/itsankoff/gotcha/common"
	"log"
)

type Control struct {
	input        chan *common.Message
	groups       []*Group
	outputStore  *OutputStore
	contactStore *ContactStore
}

func NewControl(input chan *common.Message,
	outputStore *OutputStore, contactStore *ContactStore) *Control {
	c := &Control{
		input:        input,
		outputStore:  outputStore,
		contactStore: contactStore,
	}

	go c.listen()
	return c
}

func (c Control) listen() {
	for {
		select {
		case msg := <-c.input:
			log.Println("Control received", msg.From(), msg.Cmd())
			valid := c.validate(msg)
			if valid {
				payload, err := msg.ParseJsonData()
				if err == nil {
					cmd := msg.Cmd()
					switch cmd {
					case "list_contacts":
						contacts, _ := c.contactStore.ListContacts(msg.From())
						response := common.NewResponse(msg, "contacts", contacts)
						c.outputStore.Send(response)
					case "add_contact":
						contactId := payload["contact_id"].(string)
						added := c.AddContact(msg.From(), contactId)
						response := common.NewResponse(msg, "added", added)
						c.outputStore.Send(response)
					case "remove_contact":
						contactId := payload["contact_id"].(string)
						removed := c.RemoveContact(msg.From(), contactId)
						response := common.NewResponse(msg, "removed", removed)
						c.outputStore.Send(response)
					case "search_contact":
						contactName := payload["contact_name"].(string)
						contactId := c.SearchContact(msg.From(), contactName)
						response := common.NewResponse(msg, "contact_id", contactId)
						c.outputStore.Send(response)
					case "create_group":
						groupId := c.CreateGroup()
						c.AddToGroup(groupId, msg.From())
						response := common.NewResponse(msg, "group_id", groupId)
						c.outputStore.Send(response)
					case "add_to_group":
						groupId := payload["group_id"].(string)
						userId := payload["user_id"].(string)
						added := c.AddToGroup(groupId, userId)
						response := common.NewResponse(msg, "added", added)
						c.outputStore.Send(response)
					case "remove_from_group":
						groupId := payload["group_id"].(string)
						userId := payload["user_id"].(string)
						removed := c.RemoveFromGroup(groupId, userId)
						response := common.NewResponse(msg, "removed", removed)
						c.outputStore.Send(response)
					case "delete_group":
						groupId := payload["group_id"].(string)
						deleted := c.DeleteGroup(groupId)
						response := common.NewResponse(msg, "deleted", deleted)
						c.outputStore.Send(response)
					case "list_groups":
						groups := c.ListGroups(msg.From())
						response := common.NewResponse(msg, "groups", groups)
						c.outputStore.Send(response)
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

func (c Control) ListContacts(user string) ([]string, bool) {
	contacts, listed := c.contactStore.ListContacts(user)
	log.Printf("List contacts for user %s %t", user, listed)
	return contacts, listed
}

func (c Control) AddContact(user string, contactId string) bool {
	added := c.contactStore.AddContact(user, contactId)
	log.Printf("Add contact %s for user %s %t", contactId, user, added)
	return added
}

func (c Control) RemoveContact(user string, contactId string) bool {
	removed := c.contactStore.RemoveContact(user, contactId)
	log.Printf("Remove contact %s for user %s %t", contactId, user, removed)
	return removed
}

func (c Control) SearchContact(user string, contactName string) string {
	contactId := c.contactStore.SearchContact(user, contactName)
	log.Printf("Seach contact %s for user %s contactId: %s",
		contactName, user, contactId)
	return contactId
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
		log.Printf("Failed to add user. No group with id", groupId, userId)
		return false
	}

	userOutput := c.outputStore.GetOutput(userId)
	if userOutput == nil {
		log.Printf("Failed to add user to group. No user output",
			groupId, userId)
		return false
	}

	added := group.AddOutput(userId, userOutput)
	log.Printf("Group created %s from user %s %t", groupId, userId, added)
	return added
}

func (c Control) RemoveFromGroup(groupId string, userId string) bool {
	group := c.findGroup(groupId)
	if group == nil {
		log.Printf("Failed to remove user %s from group %s. No group with id",
			userId, groupId)
		return false
	}

	removed := group.RemoveOutput(userId)
	log.Printf("User %s removed from group %s %t", groupId, userId, removed)
	return removed
}

func (c *Control) DeleteGroup(groupId string) bool {
	var deleted bool
	for i, group := range c.groups {
		if group.Id == groupId {
			c.groups = append(c.groups[:i], c.groups[i+1:]...)
			c.outputStore.RemoveOutput(group.Id)
			close(group.Out)
			deleted = true
			break
		}
	}

	log.Printf("Group %s deleted %t", groupId, deleted)
	return deleted
}

func (c Control) ListGroups(userId string) *[]string {
	groupIds := []string{}
	for _, g := range c.groups {
		groupIds = append(groupIds, g.Id)
	}

	log.Printf("List groups for user %s", userId)
	return &groupIds
}
