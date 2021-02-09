package telegram

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/jaimeteb/chatto/channels/messages"
	"github.com/jaimeteb/chatto/query"
	"github.com/kimrgrey/go-telegram"
	log "github.com/sirupsen/logrus"
)

// MessageIn models a telegram incoming message
type MessageIn struct {
	UpdateID int            `json:"update_id"`
	Message  MessageInInner `json:"message"`
}

// MessageInInner models a telegram incoming message inner struct
type MessageInInner struct {
	MessageID int                `json:"message_id"`
	From      MessageInInnerFrom `json:"from"`
	Date      int                `json:"date"`
	Text      string             `json:"text"`
}

// MessageInInnerFrom models a telegram incoming message inner struct
type MessageInInnerFrom struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

// Config models Telegram configuration
type Config struct {
	BotKey string `mapstructure:"bot_key"`
}

// Channel contains a Telegram client
type Channel struct {
	client *telegram.Client
}

// NewChannel returns an initialized telegram client
func NewChannel(config Config) *Channel {
	client := telegram.NewClient(config.BotKey)

	log.Infof("Added Telegram client: %v", client.GetMe().ID)

	return &Channel{client: client}
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
		c.client.Call(method, respValues, apiResp)

		log.Debug(*apiResp)
	}

	return nil
}

// ReceiveMessage for Telegram
func (c *Channel) ReceiveMessage(w http.ResponseWriter, r *http.Request) (*messages.Receive, error) {
	decoder := json.NewDecoder(r.Body)
	var messageIn MessageIn

	err := decoder.Decode(&messageIn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	log.Debug(messageIn)

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
	}

	return receive, nil
}

// ReceiveMessages uses event queues to receive messages. Starts a long running process
func (c *Channel) ReceiveMessages(receiveChan chan messages.Receive) {
	// Not implemented
}
