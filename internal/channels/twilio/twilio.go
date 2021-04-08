package twilio

//go:generate mockgen -source=twilio.go -destination=mocktwilio/mocktwilio.go -package=mocktwilio

import (
	"bytes"
	"net/http"
	"net/url"
	"time"

	"github.com/ajg/form"
	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/query"
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
	ProfileName      string `form:"ProfileName"`
	WaID             string `form:"WaId"`
}

// Config models Twilio configuration
type Config struct {
	AccountSid string `mapstructure:"account_sid"`
	AuthToken  string `mapstructure:"auth_token"`
	Number     string `mapstructure:"number"`
	Delay      int    `mapstructure:"delay"`
}

// Client is the twilio client interface
type Client interface {
	SendMessage(from string, to string, body string, mediaURLs []*url.URL) (*twilio.Message, error)
}

// Channel contains a Twilio client and number
type Channel struct {
	Client Client
	Number string
	token  string
	delay  int
}

// New returns an initialized telegram client
func New(config Config) *Channel {
	client := twilio.NewClient(config.AccountSid, config.AuthToken, nil)

	log.Infof("Added Twilio client: %v", client.AccountSid)

	return &Channel{Client: client.Messages, Number: config.Number, token: config.AuthToken}
}

// SendMessage for Twilio
func (c *Channel) SendMessage(response *messages.Response) error {
	for _, answer := range response.Answers {
		var imageURL []*url.URL

		if answer.Image != "" {
			u, _ := url.Parse(answer.Image)
			imageURL = append(imageURL, u)
		}

		log.Debugf("Sending Twilio message: %+v", answer)
		apiResp, err := c.Client.SendMessage(c.Number, response.ReplyOpts.Twilio.Recipient, answer.Text, imageURL)
		if err != nil {
			return err
		}
		log.Debugf("Twilio response: %+v", apiResp)

		time.Sleep(time.Duration(c.delay) * time.Second)
	}

	return nil
}

// ReceiveMessage for Twilio
func (c *Channel) ReceiveMessage(body []byte) (*messages.Receive, error) {
	byteReader := bytes.NewReader(body)

	decoder := form.NewDecoder(byteReader)

	var messageIn MessageIn
	if err := decoder.Decode(&messageIn); err != nil {
		return nil, err
	}

	receive := &messages.Receive{
		Question: &query.Question{
			Sender: messageIn.From,
			Text:   messageIn.Body,
		},
		ReplyOpts: &messages.ReplyOpts{
			Twilio: messages.TwilioReplyOpts{
				Recipient: messageIn.From,
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
	return "twilio"
}
