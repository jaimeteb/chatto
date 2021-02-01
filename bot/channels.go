package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/ajg/form"
	cmn "github.com/jaimeteb/chatto/common"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	"github.com/kevinburke/twilio-go"
	"github.com/kimrgrey/go-telegram"
	"github.com/spf13/viper"
)

// ClientsConfig struct combines all available clients configuration
type ClientsConfig struct {
	Telegram TelegramConfig `mapstructure:"telegram"`
	Twilio   TwilioConfig   `mapstructure:"twilio"`
	Slack    SlackConfig    `mapstructure:"slack"`
}

// TelegramConfig models Telegram configuration
type TelegramConfig struct {
	BotKey string `mapstructure:"bot_key"`
}

// TwilioConfig models Twilio configuration
type TwilioConfig struct {
	AccountSid string `mapstructure:"account_sid"`
	AuthToken  string `mapstructure:"auth_token"`
	Number     string `mapstructure:"number"`
}

// SlackConfig contains the Slack token
type SlackConfig struct {
	Token    string `mapstructure:"token"`
	AppToken string `mapstructure:"app_token"`
}

// Clients struct combines all available clients
type Clients struct {
	Telegram TelegramClient
	Twilio   TwilioClient
	REST     RESTClient
	Slack    SlackClient
}

// TwilioClient contains a Twilio client as well as the Twilio number
type TwilioClient struct {
	Client *twilio.Client
	Number string
}

// TelegramClient contains a Telegram client
type TelegramClient struct {
	Client *telegram.Client
}

// RESTClient contains a REST client
type RESTClient struct {
}

// SlackClient contains a Slack Client
type SlackClient struct {
	Client     *slack.Client
	Socketmode *socketmode.Client
}

// Client interface implements a SendMessage method that sends message through an API client
type Client interface {
	SendMessage(msg cmn.Message, recipient string) error
	RecieveMessage(w http.ResponseWriter, r *http.Request) (cmn.Message, error)
}

// SendMessage for Twilio
func (t *TwilioClient) SendMessage(msg cmn.Message, recipient string) error {
	var imageURL []*url.URL

	if msg.Image != "" {
		u, _ := url.Parse(msg.Image)
		imageURL = append(imageURL, u)
	}
	ret, err := t.Client.Messages.SendMessage(t.Number, recipient, msg.Text, imageURL)
	log.Debug(ret, err)
	return err
}

// RecieveMessage for Twilio
func (t *TwilioClient) RecieveMessage(w http.ResponseWriter, r *http.Request) (cmn.Message, error) {
	decoder := form.NewDecoder(r.Body)
	var twilioMessage TwilioMessageIn
	if err := decoder.Decode(&twilioMessage); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return cmn.Message{}, err
	}

	log.Debug(twilioMessage)
	sender := twilioMessage.From
	text := twilioMessage.Body
	mess := cmn.Message{
		Sender: sender,
		Text:   text,
	}

	return mess, nil
}

// SendMessage for Telegram
func (t *TelegramClient) SendMessage(msg cmn.Message, recipient string) error {
	respValues := url.Values{}
	respValues.Add("chat_id", recipient)
	respValues.Add("parse_mode", "Markdown")

	var method string
	if msg.Image != "" {
		respValues.Add("photo", msg.Image)
		respValues.Add("caption", msg.Text)
		method = "SendPhoto"
	} else {
		respValues.Add("text", msg.Text)
		method = "SendMessage"
	}

	apiResp := new(interface{})
	t.Client.Call(method, respValues, apiResp)
	log.Debug(*apiResp)

	return nil
}

// RecieveMessage for Telegram
func (t *TelegramClient) RecieveMessage(w http.ResponseWriter, r *http.Request) (cmn.Message, error) {
	decoder := json.NewDecoder(r.Body)
	var telegramMess TelegramMessageIn

	err := decoder.Decode(&telegramMess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return cmn.Message{}, err
	}

	log.Debug(telegramMess)
	sender := strconv.Itoa(telegramMess.Message.From.ID)
	mess := cmn.Message{
		Sender: sender,
		Text:   telegramMess.Message.Text,
	}

	return mess, nil
}

// SendMessage for REST
func (c *RESTClient) SendMessage(msg cmn.Message, recipient string) error {
	return nil
}

// RecieveMessage for REST
func (c *RESTClient) RecieveMessage(w http.ResponseWriter, r *http.Request) (cmn.Message, error) {
	decoder := json.NewDecoder(r.Body)
	var mess cmn.Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return cmn.Message{}, err
	}

	return mess, nil
}

// SendMessage for Slack
func (s *SlackClient) SendMessage(msg cmn.Message, recipient string) error {
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

	ret, _, err := s.Client.PostMessage(recipient, slackMsgOptions...)
	log.Debugf("%v - %v\n", ret, err)

	return nil
}

// RecieveMessage for Slack
func (s *SlackClient) RecieveMessage(w http.ResponseWriter, r *http.Request) (cmn.Message, error) {
	log.Debug(r.Body)

	var event SlackMessage

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return cmn.Message{}, err
	}

	if event.Type == "url_verification" {
		js, err := json.Marshal(map[string]string{"challenge": event.Challenge})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return cmn.Message{}, nil
	}

	if event.Event.BotID != "" {
		return cmn.Message{}, nil
	}

	log.Debug(event.Type)
	log.Debugf("%+v\n", event.Event)

	msg := cmn.Message{
		Sender: event.Event.Channel,
		Text:   event.Event.Text,
	}

	return msg, nil
}

// Messages from the fsm.
func Messages(msgs interface{}) ([]cmn.Message, []map[string]string, error) {
	// Create slice of messages
	msgsArr := make([]interface{}, 0)
	if rt := reflect.TypeOf(msgs); rt.Kind() == reflect.Slice {
		msgsArr = msgs.([]interface{})
	} else {
		msgsArr = append(msgsArr, msgs)
	}

	answer := make([]map[string]string, 0, len(msgsArr))
	messages := make([]cmn.Message, 0, len(msgsArr))

	for _, msgElem := range msgsArr {
		switch m := msgElem.(type) {
		case cmn.Message:
			answer = append(answer, m.Out())
			messages = append(messages, m)
		case string:
			msg := cmn.Message{Text: m}
			answer = append(answer, msg.Out())
			messages = append(messages, msg)
		case map[interface{}]interface{}, map[string]interface{}, map[string]string:
			msg := cmn.MessageFromMap(m)
			answer = append(answer, msg.Out())
			messages = append(messages, msg)
		default:
			err := fmt.Errorf("Message type unsupported: %T", m)
			return nil, nil, err
		}
	}

	return messages, answer, nil
}

// SendMessages sends messages through the clients
func SendMessages(msgs interface{}, client Client, recipient string, w http.ResponseWriter) error {
	ans := make([]map[string]string, 0)

	// Create slice of messages
	msgsArr := make([]interface{}, 0)
	if rt := reflect.TypeOf(msgs); rt.Kind() == reflect.Slice {
		msgsArr = msgs.([]interface{})
	} else {
		msgsArr = append(msgsArr, msgs)
	}

	for _, msgElem := range msgsArr {
		switch m := msgElem.(type) {
		case cmn.Message:
			ans = append(ans, m.Out())
			if err := client.SendMessage(m, recipient); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}
		case string:
			msg := cmn.Message{
				Text: m,
			}
			ans = append(ans, msg.Out())
			if err := client.SendMessage(msg, recipient); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}
		case map[interface{}]interface{}, map[string]interface{}, map[string]string:
			msg := cmn.MessageFromMap(m)
			ans = append(ans, msg.Out())
			if err := client.SendMessage(msg, recipient); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return err
			}
		default:
			err := fmt.Errorf("Message type unsupported: %T", m)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}
	}

	js, err := json.Marshal(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	return nil
}

// LoadClients loads registered clients/channels in the chn.yml file
func LoadClients(path *string) Clients {
	config := viper.New()
	config.SetConfigName("chn")
	config.AddConfigPath(*path)
	config.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)

	var cts Clients

	if err := config.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warn("File chn.yml not found, skipping channels")
		default:
			log.Warn(err)
		}
		return cts
	}

	var end ClientsConfig
	if err := config.Unmarshal(&end); err != nil {
		log.Warn(err)
		return cts
	}

	// TELEGRAM
	if end.Telegram != (TelegramConfig{}) {
		telegramClient := telegram.NewClient(end.Telegram.BotKey)
		cts.Telegram = TelegramClient{telegramClient}
		log.Infof("Added Telegram client: %v\n", telegramClient.GetMe().ID)
	}

	// TWILIO
	if end.Twilio != (TwilioConfig{}) {
		twilioClient := twilio.NewClient(end.Twilio.AccountSid, end.Twilio.AuthToken, nil)
		cts.Twilio = TwilioClient{twilioClient, end.Twilio.Number}
		log.Infof("Added Twilio client: %v\n", twilioClient.AccountSid)
	}

	// SLACK
	if end.Slack != (SlackConfig{}) {
		var slackOpts []slack.Option

		if end.Slack.AppToken != "" {
			slackOpts = append(slackOpts, slack.OptionAppLevelToken(end.Slack.AppToken))
		}

		slackClient := slack.New(end.Slack.Token, slackOpts...)
		cts.Slack = SlackClient{Client: slackClient}

		if end.Slack.AppToken != "" {
			cts.Slack.Socketmode = socketmode.New(slackClient)
		}

		log.Infof("Added Slack client: %v...\n", end.Slack.Token[:10])
	}

	return cts
}
