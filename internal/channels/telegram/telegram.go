package telegram

//go:generate mockgen -source=telegram.go -destination=mocktelegram/mocktelegram.go -package=mocktelegram

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jaimeteb/chatto/internal/channels/message"
	"github.com/jaimeteb/chatto/query"
	"github.com/kimrgrey/go-telegram"
	log "github.com/sirupsen/logrus"
)

// MessageIn models a Telegram incoming message
type MessageIn struct {
	UpdateID int            `json:"update_id"`
	Message  MessageInInner `json:"message"`
}

// MessageInInner models a Telegram incoming message inner struct
type MessageInInner struct {
	MessageID int                `json:"message_id"`
	From      MessageInInnerFrom `json:"from"`
	Date      int                `json:"date"`
	Text      string             `json:"text"`
}

// MessageInInnerFrom models a Telegram incoming message inner struct
type MessageInInnerFrom struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

// Config models Telegram configuration
type Config struct {
	BotKey string `mapstructure:"bot_key"`
}

// Client is the Telegram client interface
type Client interface {
	Call(method string, params url.Values, v interface{})
}

// Channel contains a Telegram client
type Channel struct {
	Client Client
}

// New returns an initialized Telegram client
func New(config Config) *Channel {
	client := telegram.NewClient(config.BotKey)

	log.Infof("Added Telegram client: %v", client.GetMe().ID)

	return &Channel{Client: client}
}

// MessageResponse for Telegram. See interface for more details
func (c *Channel) MessageResponse(msgResponse *message.Response) error {
	for _, answer := range msgResponse.Answers {
		respValues := url.Values{}
		respValues.Add("chat_id", msgResponse.ReplyOpts.Telegram.Recipient)
		respValues.Add("parse_mode", "Markdown")

		var method string

		if answer.Image != "" {
			respValues.Add("photo", answer.Image)
			respValues.Add("caption", answer.Text)
			method = "SendPhoto"
		} else {
			respValues.Add("text", answer.Text)
			method = "MessageResponse"
		}

		apiResp := new(interface{})
		c.Client.Call(method, respValues, apiResp)

		log.Debug(*apiResp)
	}

	return nil
}

// MessageRequest for Telegram. See interface for more details
func (c *Channel) MessageRequest(body []byte) (*message.Request, error) {
	var messageIn MessageIn
	err := json.Unmarshal(body, &messageIn)
	if err != nil {
		return nil, err
	}

	sender := strconv.Itoa(messageIn.Message.From.ID)

	msgRequest := &message.Request{
		Question: &query.Question{
			Sender: sender,
			Text:   messageIn.Message.Text,
		},
		ReplyOpts: &message.ReplyOpts{
			Telegram: message.TelegramReplyOpts{
				Recipient: sender,
			},
		},
		Channel: c.String(),
	}

	return msgRequest, nil
}

// MessageRequestQueue for Telegram is not implemented. See interface for more details
func (c *Channel) MessageRequestQueue(_ chan message.Request) {
	// Not implemented
}

// ValidateCallback for Telegram is not implemented. See interface for more details
func (c *Channel) ValidateCallback(_ *http.Request) bool {
	// TODO: Implement callback validation
	return true
}

// String returns Telegram channel name. See interface for more details
func (c *Channel) String() string {
	return "telegram"
}
