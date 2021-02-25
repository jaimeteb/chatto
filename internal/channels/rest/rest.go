package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/query"
)

// MessageIn from REST client
type MessageIn struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

// Config models REST channel configuration
type Config struct {
	CallbackToken string `mapstructure:"callback_token"`
}

// Channel contains a REST client
type Channel struct {
	token string
}

// New returns an initialized REST client/channel
func New(config Config) *Channel {
	return &Channel{config.CallbackToken}
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

// ValidateCallback validates a callback to the channel
func (c *Channel) ValidateCallback(r *http.Request) bool {
	if c.token != "" {
		reqToken := r.Header.Get("Authorization")
		reqToken = strings.TrimPrefix(reqToken, "Bearer ")

		if c.token != reqToken {
			return false
		}
	}
	return true
}
