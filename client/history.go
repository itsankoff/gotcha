package client

import (
	"errors"
	"github.com/itsankoff/gotcha/common"
	"time"
)

// GetHistory retrieves the conversation history for a user
// or a group
func (c *Client) GetHistory(remote string, from time.Time, to time.Time) (*[]common.Message, error) {
	var messages *[]common.Message
	return messages, errors.New("Not Implemented")
}
