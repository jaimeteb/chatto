package channels

import (
	"net/http"
	"strings"

	"github.com/jaimeteb/chatto/channels/messages"
	"github.com/jaimeteb/chatto/channels/rest"
	"github.com/jaimeteb/chatto/channels/slack"
	"github.com/jaimeteb/chatto/channels/telegram"
	"github.com/jaimeteb/chatto/channels/twilio"
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
	Telegram Channel
	Twilio   Channel
	REST     Channel
	Slack    Channel
}

// Channel interface implements a channel to send and receive messages on
type Channel interface {
	// ReceiveMessage from the channel
	ReceiveMessage(w http.ResponseWriter, r *http.Request) (*messages.Receive, error)
	// ReceiveMessages from the channel. Starts a long running process, receives questions and sends them to the receiveChan
	ReceiveMessages(receiveChan chan messages.Receive)
	// SendMessage to the channel
	SendMessage(response *messages.Response) error
}

// Load registered clients/channels in the chn.yml file
func Load(path *string) *Channels {
	config := viper.New()
	config.SetConfigName("chn")
	config.AddConfigPath(*path)
	config.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)

	chnls := Channels{}

	// REST
	chnls.REST = &rest.Channel{}

	if err := config.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warn("File chn.yml not found, using only REST channel")
		default:
			log.Warn(err)
		}
		return &chnls
	}

	var cfg Config
	if err := config.Unmarshal(&cfg); err != nil {
		log.Warn(err)
		return &chnls
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

	return &chnls
}
