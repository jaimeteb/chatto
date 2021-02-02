package twilio

import (
	"net/http"
	"net/url"

	"github.com/ajg/form"
	"github.com/jaimeteb/chatto/message"
	"github.com/kevinburke/twilio-go"
	log "github.com/sirupsen/logrus"
)

// MessageIn models an incoming Twilio message
type MessageIn struct {
	From             string `form:"From"`
	Body             string `form:"Body"`
	To               string `form:"To"`
	MediaURL         string `form:"MediaUrl"`
	MediaContentType string `form:"MediaContentType"`
	MessageSid       string `form:"MessageSid"`
	SmsStatus        string `form:"SmsStatus"`
	AccountSid       string `form:"AccountSid"`
	Sid              string `form:"Sid"`
	SmsSid           string `form:"SmsSid"`
	SmsMessageSid    string `form:"SmsMessageSid"`
	NumMedia         int    `form:"NumMedia"`
	NumSegments      int    `form:"NumSegments"`
	APIVersion       string `form:"ApiVersion"`
}

// Config models Twilio configuration
type Config struct {
	AccountSid string `mapstructure:"account_sid"`
	AuthToken  string `mapstructure:"auth_token"`
	Number     string `mapstructure:"number"`
}

// Channel contains a Twilio client as well as the Twilio number
type Channel struct {
	client *twilio.Client
	number string
}

// NewChannel returns an initialized telegram client
func NewChannel(config Config) *Channel {
	client := twilio.NewClient(config.AccountSid, config.AuthToken, nil)

	log.Infof("Added Twilio client: %v\n", client.AccountSid)

	return &Channel{client: client}
}

// SendMessage for Twilio
func (c *Channel) SendMessage(msg message.Message, recipient string) error {
	var imageURL []*url.URL

	if msg.Image != "" {
		u, _ := url.Parse(msg.Image)
		imageURL = append(imageURL, u)
	}
	ret, err := c.client.Messages.SendMessage(c.number, recipient, msg.Text, imageURL)
	log.Debug(ret, err)
	return err
}

// ReceiveMessage for Twilio
func (c *Channel) ReceiveMessage(w http.ResponseWriter, r *http.Request) (message.Message, error) {
	decoder := form.NewDecoder(r.Body)
	var messageIn MessageIn
	if err := decoder.Decode(&messageIn); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return message.Message{}, err
	}

	log.Debug(messageIn)
	sender := messageIn.From
	text := messageIn.Body
	mess := message.Message{
		Sender: sender,
		Text:   text,
	}

	return mess, nil
}

// ReceiveMessages uses event queues to receive messages. Starts a long running process
func (c *Channel) ReceiveMessages(messageChan chan message.Message) {
	// Not implemented
}
