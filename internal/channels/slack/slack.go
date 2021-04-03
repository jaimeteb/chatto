package slack

//go:generate mockgen -source=slack.go -destination=mockslack/mockslack.go -package=mockslack

import (
	"encoding/json"
	"net/http"

	"github.com/jaimeteb/chatto/internal/channels/message"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// MessageIn models a Slack message and/or Slack endpoint challenge
type MessageIn struct {
	Challenge string    `json:"challenge"`
	Type      string    `json:"type"`
	Event     slack.Msg `json:"event"`
	Token     string    `json:"token"`
}

// Config contains the Slack token
type Config struct {
	Token    string `mapstructure:"token"`
	AppToken string `mapstructure:"app_token"`
}

// Client is the Slack client interface
type Client interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

// SocketClient is the Slack socketmode client interface
type SocketClient interface {
	Ack(req socketmode.Request, payload ...interface{})
	Run() error
}

// Channel contains a Slack Channel
type Channel struct {
	Client             Client
	SocketClient       SocketClient
	SocketClientEvents chan socketmode.Event
}

// New returns an initialized slack client
func New(config Config) *Channel {
	var slackOpts []slack.Option

	if config.AppToken != "" {
		slackOpts = append(slackOpts, slack.OptionAppLevelToken(config.AppToken))
	}

	slackClient := slack.New(config.Token, slackOpts...)

	client := &Channel{Client: slackClient}

	if config.AppToken != "" {
		socketclient := socketmode.New(slackClient)
		client.SocketClient = socketclient
		client.SocketClientEvents = socketclient.Events
	}

	log.Info("Added Slack client")

	return client
}

// MessageResponse for Slack. See interface for more details
func (c *Channel) MessageResponse(msgResponse *message.Response) error {
	for _, answer := range msgResponse.Answers {
		slackMsgOptions := []slack.MsgOption{}

		if answer.Image != "" {
			var imageText *slack.TextBlockObject
			if answer.Text != "" {
				imageText = slack.NewTextBlockObject("plain_text", answer.Text, false, false)
			} else {
				imageText = nil
			}
			image := slack.MsgOptionBlocks(slack.NewImageBlock(answer.Image, "image", "1", imageText))
			slackMsgOptions = append(slackMsgOptions, image)
		} else if answer.Text != "" {
			text := slack.MsgOptionText(answer.Text, false)
			slackMsgOptions = append(slackMsgOptions, text)
		}

		if msgResponse.ReplyOpts.Slack.TS != "" {
			slackMsgOptions = append(slackMsgOptions, slack.MsgOptionTS(msgResponse.ReplyOpts.Slack.TS))
		}

		ret, _, err := c.Client.PostMessage(msgResponse.ReplyOpts.Slack.Channel, slackMsgOptions...)
		if err != nil {
			log.Errorf("%s: %+v", err, ret)
			return err
		}
	}

	return nil
}

// MessageRequest for Slack. See interface for more details
func (c *Channel) MessageRequest(body []byte) (*message.Request, error) {
	var slackMsg MessageIn
	err := json.Unmarshal(body, &slackMsg)
	if err != nil {
		return nil, err
	}

	if slackMsg.Type == "url_verification" {
		challenge, err := json.Marshal(map[string]string{"challenge": slackMsg.Challenge})
		if err != nil {
			return nil, err
		}

		return nil, ErrURLVerification{Challenge: challenge}
	}

	if slackMsg.Event.BotID != "" {
		return &message.Request{}, nil
	}

	log.Debug(slackMsg.Type)
	log.Debugf("%+v", slackMsg.Event)

	ts := slackMsg.Event.Timestamp
	if slackMsg.Event.ThreadTimestamp != "" {
		ts = slackMsg.Event.ThreadTimestamp
	}

	msgRequest := &message.Request{
		Question: &query.Question{
			Text:   slackMsg.Event.Text,
			Sender: slackMsg.Event.User,
		},
		ReplyOpts: &message.ReplyOpts{
			Slack: message.SlackReplyOpts{
				Channel: slackMsg.Event.Channel,
				TS:      ts,
			},
		},
		Channel: c.String(),
	}

	return msgRequest, nil
}

// MessageRequestQueue for Slack. See interface for more details
func (c *Channel) MessageRequestQueue(receiveChan chan message.Request) {
	defer close(receiveChan)

	if c.SocketClient == nil {
		return
	}

	go func() {
		for evt := range c.SocketClientEvents {
			switch evt.Type {
			case socketmode.EventTypeHello, socketmode.EventTypeInteractive, socketmode.EventTypeSlashCommand:
				// Ignore
			case socketmode.EventTypeInvalidAuth:
				log.Error("Invalid auth when connecting to Slack...")
			case socketmode.EventTypeIncomingError:
				log.Error("Event type incoming error from Slack...")
			case socketmode.EventTypeErrorWriteFailed:
				log.Error("Writing event message to Slack failed...")
			case socketmode.EventTypeErrorBadMessage:
				log.Error("Bad event message from Slack...")
			case socketmode.EventTypeDisconnect:
				log.Warn("Disconnected from Slack...")
			case socketmode.EventTypeConnecting:
				log.Info("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				log.Error("Connection to Slack failed. Retrying later...")
			case socketmode.EventTypeConnected:
				log.Info("Connected to Slack with Socket Mode")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					log.Warnf("Ignored %+v", evt)
					continue
				}
				c.SocketClient.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						if ev.BotID != "" {
							continue
						}

						ts := ev.TimeStamp
						if ev.ThreadTimeStamp != "" {
							ts = ev.ThreadTimeStamp
						}

						receiveChan <- message.Request{
							Question: &query.Question{
								Text:   ev.Text,
								Sender: ev.User,
							},
							ReplyOpts: &message.ReplyOpts{
								Slack: message.SlackReplyOpts{
									Channel: ev.Channel,
									TS:      ts,
								},
							},
							Channel: c.String(),
						}
					case *slackevents.AppMentionEvent:
						if ev.BotID != "" {
							continue
						}

						ts := ev.TimeStamp
						if ev.ThreadTimeStamp != "" {
							ts = ev.ThreadTimeStamp
						}

						receiveChan <- message.Request{
							Question: &query.Question{
								Text:   ev.Text,
								Sender: ev.User,
							},
							ReplyOpts: &message.ReplyOpts{
								Slack: message.SlackReplyOpts{
									Channel: ev.Channel,
									TS:      ts,
								},
							},
							Channel: c.String(),
						}
					}
				default:
					log.Debugf("Unsupported Events API event received")
				}
			default:
				log.Debugf("Unexpected event type received: %s", evt.Type)
			}
		}
	}()

	err := c.SocketClient.Run()
	if err != nil {
		log.Error(err)
	}
}

// ValidateCallback for Slack not implemented. See interface for more details
func (c *Channel) ValidateCallback(r *http.Request) bool {
	// TODO: Implement callback validation
	return true
}

// String returns Slack channel name. See interface for more details
func (c *Channel) String() string {
	return "slack"
}

// ErrURLVerification raised when an auth challenge is supposed to be performed
type ErrURLVerification struct {
	Challenge []byte
}

// Error message raised when an auth challenge is supposed to be performed
func (e ErrURLVerification) Error() string {
	return "must perform challenge auth verification"
}
