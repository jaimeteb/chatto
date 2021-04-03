package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jaimeteb/chatto/internal/channels/message"
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

// MessageResponse for REST. See interface for more details
func (c *Channel) MessageResponse(_ *message.Response) error {
	// Not implemented
	return nil
}

// MessageRequest for REST. See interface for more details
func (c *Channel) MessageRequest(body []byte) (*message.Request, error) {
	var messageIn MessageIn
	err := json.Unmarshal(body, &messageIn)
	if err != nil {
		return nil, err
	}

	receive := &message.Request{
		Question: &query.Question{
			Text:   messageIn.Text,
			Sender: messageIn.Sender,
		},
		Channel: c.String(),
	}

	return receive, nil
}

// MessageRequestQueue for REST is not implemented. See interface for more details
func (c *Channel) MessageRequestQueue(_ chan message.Request) {
	// Not implemented
}

// ValidateCallback for REST validates a callback to the channel. See interface for more details
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

// String returns REST channel name. See interface for more details
func (c *Channel) String() string {
	return "rest"
}
