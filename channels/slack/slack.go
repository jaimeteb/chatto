package slack

//go:generate mockgen -source=slack.go -destination=mockslack/mockslack.go -package=mockslack

import (
	"encoding/json"

	"github.com/jaimeteb/chatto/channels/messages"
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

// SendMessage to Slack with the bots response
func (c *Channel) SendMessage(response *messages.Response) error {
	for _, answer := range response.Answers {
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

		if response.ReplyOpts.Slack.TS != "" {
			slackMsgOptions = append(slackMsgOptions, slack.MsgOptionTS(response.ReplyOpts.Slack.TS))
		}

		ret, _, err := c.Client.PostMessage(response.ReplyOpts.Slack.Channel, slackMsgOptions...)
		if err != nil {
			log.Errorf("%s: %+v", err, ret)
			return err
		}
	}

	return nil
}

// ReceiveMessage for Slack
func (c *Channel) ReceiveMessage(body []byte) (*messages.Receive, error) {
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
		return &messages.Receive{}, nil
	}

	log.Debug(slackMsg.Type)
	log.Debugf("%+v", slackMsg.Event)

	ts := slackMsg.Event.Timestamp
	if slackMsg.Event.ThreadTimestamp != "" {
		ts = slackMsg.Event.ThreadTimestamp
	}

	receive := &messages.Receive{
		Question: &query.Question{
			Text:   slackMsg.Event.Text,
			Sender: slackMsg.Event.User,
		},
		ReplyOpts: &messages.ReplyOpts{
			Slack: messages.SlackReplyOpts{
				Channel: slackMsg.Event.Channel,
				TS:      ts,
			},
		},
	}

	return receive, nil
}

// ReceiveMessages uses event queues to receive messages. Starts a long running process
func (c *Channel) ReceiveMessages(receiveChan chan messages.Receive) {
	defer close(receiveChan)

	if c.SocketClient == nil {
		return
	}

	go func() {
		for evt := range c.SocketClientEvents {
			switch evt.Type {
			case socketmode.EventTypeHello:
				// Ignore
			case socketmode.EventTypeInteractive:
				// Ignore
			case socketmode.EventTypeSlashCommand:
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
							// Do not interact with bots.
							continue
						}

						ts := ev.TimeStamp
						if ev.ThreadTimeStamp != "" {
							ts = ev.ThreadTimeStamp
						}

						receiveChan <- messages.Receive{
							Question: &query.Question{
								Text:   ev.Text,
								Sender: ev.User,
							},
							ReplyOpts: &messages.ReplyOpts{
								Slack: messages.SlackReplyOpts{
									Channel: ev.Channel,
									TS:      ts,
								},
							},
						}
					case *slackevents.AppMentionEvent:
						if ev.BotID != "" {
							// Do not interact with bots.
							continue
						}

						ts := ev.TimeStamp
						if ev.ThreadTimeStamp != "" {
							ts = ev.ThreadTimeStamp
						}

						receiveChan <- messages.Receive{
							Question: &query.Question{
								Text:   ev.Text,
								Sender: ev.User,
							},
							ReplyOpts: &messages.ReplyOpts{
								Slack: messages.SlackReplyOpts{
									Channel: ev.Channel,
									TS:      ts,
								},
							},
						}
					}
				default:
					log.Debugf("Unsupported Events API event received")
				}
			default:
				// log.Debugf("Unexpected event type received: %s", evt.Type)
			}
		}
	}()

	err := c.SocketClient.Run()
	if err != nil {
		log.Error(err)
	}
}

// ErrURLVerification raised when an auth challenge is supposed to be performed
type ErrURLVerification struct {
	Challenge []byte
}

func (e ErrURLVerification) Error() string {
	return "must perform challenge auth verification"
}
