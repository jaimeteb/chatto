package channels

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/jaimeteb/chatto/channels/options"
	"github.com/jaimeteb/chatto/channels/rest"
	"github.com/jaimeteb/chatto/channels/slack"
	"github.com/jaimeteb/chatto/channels/telegram"
	"github.com/jaimeteb/chatto/channels/twilio"
	"github.com/jaimeteb/chatto/message"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config struct combines all available client configurations
type Config struct {
	Telegram telegram.Config `mapstructure:"telegram"`
	Twilio   twilio.Config   `mapstructure:"twilio"`
	Slack    slack.Config    `mapstructure:"slack"`
}

// Channels combines all available channel clients
type Channels struct {
	Telegram *telegram.Channel
	Twilio   *twilio.Channel
	REST     *rest.Channel
	Slack    *slack.Channel
}

// Channel interface implements a channel to send and receive messages on
type Channel interface {
	// SendMessage to the channel
	SendMessage(msg message.Message, sendOpts options.SendOptions) error
	// ReceiveMessage from the channel
	ReceiveMessage(w http.ResponseWriter, r *http.Request) (message.Message, error)
	// ReceiveMessages starts a long running process, receives message events and sends them to the messageChan
	ReceiveMessages(messageChan chan message.Message)
}

// SendMessages through the channel
func SendMessages(msgs interface{}, chnl Channel, sendOpts options.SendOptions) ([]map[string]string, error) {
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
		case message.Message:
			ans = append(ans, m.Out())
			if err := chnl.SendMessage(m, sendOpts); err != nil {
				return nil, err
			}
		case string:
			msg := message.Message{
				Text: m,
			}
			ans = append(ans, msg.Out())
			if err := chnl.SendMessage(msg, sendOpts); err != nil {
				return nil, err
			}
		case map[interface{}]interface{}, map[string]interface{}, map[string]string:
			msg := message.FromMap(m)
			ans = append(ans, msg.Out())
			if err := chnl.SendMessage(msg, sendOpts); err != nil {
				return nil, err
			}
		default:
			err := fmt.Errorf("Message type unsupported: %T", m)
			return nil, err
		}
	}

	return ans, nil
}

// Load registered clients/channels in the chn.yml file
func Load(path *string) *Channels {
	config := viper.New()
	config.SetConfigName("chn")
	config.AddConfigPath(*path)
	config.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)

	chnls := &Channels{}

	if err := config.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warn("File chn.yml not found, skipping channels")
		default:
			log.Warn(err)
		}
		return nil
	}

	var cfg Config
	if err := config.Unmarshal(&cfg); err != nil {
		log.Warn(err)
		return nil
	}

	// TELEGRAM
	if cfg.Telegram != (telegram.Config{}) {
		chnls.Telegram = telegram.NewChannel(cfg.Telegram)
	}

	// TWILIO
	if cfg.Twilio != (twilio.Config{}) {
		chnls.Twilio = twilio.NewChannel(cfg.Twilio)
	}

	// SLACK
	if cfg.Slack != (slack.Config{}) {
		chnls.Slack = slack.NewChannel(cfg.Slack)
	}

	return chnls
}
