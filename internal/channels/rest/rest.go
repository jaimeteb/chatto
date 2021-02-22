package rest

import (
	"encoding/json"

	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/query"
)

// MessageIn from REST client
type MessageIn struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

// Channel contains a REST client
type Channel struct {
}

// SendMessage for REST
func (c *Channel) SendMessage(response *messages.Response) error {
	// Not implemented
	return nil
}

// ReceiveMessage for REST
func (c *Channel) ReceiveMessage(body []byte) (*messages.Receive, error) {
	var messageIn MessageIn
	err := json.Unmarshal(body, &messageIn)
	if err != nil {
		return nil, err
	}

	receive := &messages.Receive{
		Question: &query.Question{
			Text:   messageIn.Text,
			Sender: messageIn.Sender,
		},
	}

	return receive, nil
}

// ReceiveMessages uses event queues to receive messages. Starts a long running process
func (c *Channel) ReceiveMessages(receiveChan chan messages.Receive) {
	// Not implemented
}
