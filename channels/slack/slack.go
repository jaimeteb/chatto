package slack

import (
	"encoding/json"
	"net/http"

	"github.com/jaimeteb/chatto/channels/options"
	"github.com/jaimeteb/chatto/message"
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

	log.Infof("Added Slack client: %v...\n", config.Token[:10])

	return client
}

// SendMessage for Slack
func (c *Channel) SendMessage(msg message.Message, sendOpts options.SendOptions) error {
	slackMsgOptions := []slack.MsgOption{}

	if msg.Image != "" {
		var imageText *slack.TextBlockObject
		if msg.Text != "" {
			imageText = slack.NewTextBlockObject("plain_text", msg.Text, false, false)
		} else {
			imageText = nil
		}
		image := slack.MsgOptionBlocks(slack.NewImageBlock(msg.Image, "image", "1", imageText))
		slackMsgOptions = append(slackMsgOptions, image)
	} else if msg.Text != "" {
		text := slack.MsgOptionText(msg.Text, false)
		slackMsgOptions = append(slackMsgOptions, text)
	}

	if sendOpts.TS != "" {
		slackMsgOptions = append(slackMsgOptions, slack.MsgOptionTS(sendOpts.TS))
	}

	ret, _, err := c.client.PostMessage(sendOpts.Recipient, slackMsgOptions...)
	log.Debugf("%v - %v\n", ret, err)

	return nil
}

// ReceiveMessage for Slack
func (c *Channel) ReceiveMessage(w http.ResponseWriter, r *http.Request) (message.Message, error) {
	log.Debug(r.Body)

	var slackMsg Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&slackMsg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return message.Message{}, err
	}

	if slackMsg.Type == "url_verification" {
		js, err := json.Marshal(map[string]string{"challenge": slackMsg.Challenge})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return message.Message{}, nil
	}

	if slackMsg.Event.BotID != "" {
		return message.Message{}, nil
	}

	log.Debug(slackMsg.Type)
	log.Debugf("%+v\n", slackMsg.Event)

	msg := message.Message{
		Sender: slackMsg.Event.Channel,
		Text:   slackMsg.Event.Text,
	}

	return msg, nil
}

// ReceiveMessages uses event queues to receive messages. Starts a long running process
func (c *Channel) ReceiveMessages(messageChan chan message.Message) {
	defer close(messageChan)

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
					log.Warnf("Ignored %+v\n", evt)
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

						messageChan <- message.Message{Sender: ev.Channel, Text: ev.Text}
					case *slackevents.AppMentionEvent:
						if ev.BotID != "" {
							// Do not interact with bots.
							continue
						}

						messageChan <- message.Message{Sender: ev.Channel, Text: ev.Text}
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
