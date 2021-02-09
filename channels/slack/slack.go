package slack

import (
	"encoding/json"
	"net/http"

	"github.com/jaimeteb/chatto/channels/messages"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// Message models a Slack message and/or Slack endpoint challenge
type Message struct {
	Challenge string    `json:"challenge"`
	Type      string    `json:"type"`
	Event     slack.Msg `json:"event"`
}

// Config contains the Slack token
type Config struct {
	Token    string `mapstructure:"token"`
	AppToken string `mapstructure:"app_token"`
}

// Channel contains a Slack Channel
type Channel struct {
	client       *slack.Client
	socketclient *socketmode.Client
}

// NewChannel returns an initialized slack client
func NewChannel(config Config) *Channel {
	var slackOpts []slack.Option

	if config.AppToken != "" {
		slackOpts = append(slackOpts, slack.OptionAppLevelToken(config.AppToken))
	}

	slackClient := slack.New(config.Token, slackOpts...)

	client := &Channel{client: slackClient}

	if config.AppToken != "" {
		client.socketclient = socketmode.New(slackClient)
	}

	log.Infof("Added Slack client: %v...", config.Token[:10])

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

		ret, _, err := c.client.PostMessage(response.ReplyOpts.Slack.Channel, slackMsgOptions...)
		if err != nil {
			log.Errorf("%s: %+v", err, ret)
			return err
		}
	}

	return nil
}

// ReceiveMessage for Slack
func (c *Channel) ReceiveMessage(w http.ResponseWriter, r *http.Request) (*messages.Receive, error) {
	log.Debug(r.Body)

	var slackMsg Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&slackMsg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	if slackMsg.Type == "url_verification" {
		js, err := json.Marshal(map[string]string{"challenge": slackMsg.Challenge})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

		return &messages.Receive{}, nil
	}

	if slackMsg.Event.BotID != "" {
		return &messages.Receive{}, nil
	}

	log.Debug(slackMsg.Type)
	log.Debugf("%+v", slackMsg.Event)

	receive := &messages.Receive{
		Question: &query.Question{
			Text:   slackMsg.Event.Text,
			Sender: slackMsg.Event.User,
		},
		ReplyOpts: &messages.ReplyOpts{
			Slack: messages.SlackReplyOpts{
				Channel: slackMsg.Event.Channel,
				TS:      slackMsg.Event.Timestamp,
			},
		},
	}

	return receive, nil
}

// ReceiveMessages uses event queues to receive messages. Starts a long running process
func (c *Channel) ReceiveMessages(receiveChan chan messages.Receive) {
	defer close(receiveChan)

	if c.socketclient == nil {
		return
	}

	go func() {
		for evt := range c.socketclient.Events {
			switch evt.Type {
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

				c.socketclient.Ack(*evt.Request)

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

	c.socketclient.Run()
}
