package twilio

//go:generate mockgen -source=twilio.go -destination=mocktwilio/mocktwilio.go -package=mocktwilio

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/ajg/form"
	"github.com/jaimeteb/chatto/internal/channels/message"
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
}

// New returns an initialized telegram client
func New(config Config) *Channel {
	client := twilio.NewClient(config.AccountSid, config.AuthToken, nil)

	log.Infof("Added Twilio client: %v", client.AccountSid)

	return &Channel{Client: client.Messages, Number: config.Number, token: config.AuthToken}
}

// MessageResponse for Twilio. See interface for more details
func (c *Channel) MessageResponse(msgResponse *message.Response) error {
	for _, answer := range msgResponse.Answers {
		var imageURL []*url.URL

		if answer.Image != "" {
			u, _ := url.Parse(answer.Image)
			imageURL = append(imageURL, u)
		}

		_, err := c.Client.SendMessage(c.Number, msgResponse.ReplyOpts.Twilio.Recipient, answer.Text, imageURL)
		if err != nil {
			return err
		}
	}

	return nil
}

// MessageRequest for Twilio. See interface for more details
func (c *Channel) MessageRequest(body []byte) (*message.Request, error) {
	byteReader := bytes.NewReader(body)

	decoder := form.NewDecoder(byteReader)

	var messageIn MessageIn
	if err := decoder.Decode(&messageIn); err != nil {
		return nil, err
	}

	msgRequest := &message.Request{
		Question: &query.Question{
			Sender: messageIn.From,
			Text:   messageIn.Body,
		},
		ReplyOpts: &message.ReplyOpts{
			Twilio: message.TwilioReplyOpts{
				Recipient: messageIn.From,
			},
		},
		Channel: c.String(),
	}

	return msgRequest, nil
}

// MessageRequestQueue for Twilio is not implemented. See interface for more details
func (c *Channel) MessageRequestQueue(receiveChan chan message.Request) {
	// Not implemented
}

// ValidateCallback for Twilio not implemented. See interface for more details
func (c *Channel) ValidateCallback(r *http.Request) bool {
	// TODO: Implement callback validation
	return true
}

// String returns Twilio channel name. See interface for more details
func (c *Channel) String() string {
	return "twilio"
}
