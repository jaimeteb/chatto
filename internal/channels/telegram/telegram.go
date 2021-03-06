package telegram

//go:generate mockgen -source=telegram.go -destination=mocktelegram/mocktelegram.go -package=mocktelegram

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jaimeteb/chatto/internal/channels/messages"
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
	BotKey string        `mapstructure:"bot_key"`
	Delay  time.Duration `mapstructure:"delay"`
}

// Client is the Telegram client interface
type Client interface {
	Call(method string, params url.Values, v interface{})
}

// Channel contains a Telegram client
type Channel struct {
	Client Client
	delay  time.Duration
}

// New returns an initialized Telegram client
func New(config Config) *Channel {
	client := telegram.NewClient(config.BotKey)

	log.Infof("Added Telegram client: %v", client.GetMe().ID)

	return &Channel{Client: client, delay: config.Delay}
}

// SendMessage for Telegram
func (c *Channel) SendMessage(response *messages.Response) error {
	for _, answer := range response.Answers {
		respValues := url.Values{}
		respValues.Add("chat_id", response.ReplyOpts.Telegram.Recipient)
		respValues.Add("parse_mode", "Markdown")

		var method string

		if answer.Image != "" {
			respValues.Add("photo", answer.Image)
			respValues.Add("caption", answer.Text)
			method = "SendPhoto"
		} else {
			respValues.Add("text", answer.Text)
			method = "SendMessage"
		}

		apiResp := new(interface{})
		log.Debugf("Sending Telegram message: %+v", answer)
		c.Client.Call(method, respValues, apiResp)
		log.Debugf("Telegram response: %+v", apiResp)

		time.Sleep(c.delay)
	}

	return nil
}

// ReceiveMessage for Telegram
func (c *Channel) ReceiveMessage(body []byte) (*messages.Receive, error) {
	var messageIn MessageIn
	err := json.Unmarshal(body, &messageIn)
	if err != nil {
		return nil, err
	}

	sender := strconv.Itoa(messageIn.Message.From.ID)

	receive := &messages.Receive{
		Question: &query.Question{
			Sender: sender,
			Text:   messageIn.Message.Text,
		},
		ReplyOpts: &messages.ReplyOpts{
			Telegram: messages.TelegramReplyOpts{
				Recipient: sender,
			},
		},
		Channel: c.String(),
	}

	return receive, nil
}

// ReceiveMessages uses event queues to receive messages. Starts a long running process
func (c *Channel) ReceiveMessages(receiveChan chan messages.Receive) {
	// Not implemented
}

// ValidateCallback validates a callback to the channel
func (c *Channel) ValidateCallback(r *http.Request) bool {
	// Not implemented
	return true
}

func (c *Channel) String() string {
	return "telegram"
}
