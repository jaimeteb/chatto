package channels

//go:generate mockgen -source=channels.go -destination=mockchannels/mockchannels.go -package=mockchannels

import (
	"net/http"
	"strings"

	"github.com/jaimeteb/chatto/internal/channels/message"
	"github.com/jaimeteb/chatto/internal/channels/rest"
	"github.com/jaimeteb/chatto/internal/channels/slack"
	"github.com/jaimeteb/chatto/internal/channels/telegram"
	"github.com/jaimeteb/chatto/internal/channels/twilio"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config struct combines all available client configurations
type Config struct {
	Telegram telegram.Config `mapstructure:"telegram"`
	Twilio   twilio.Config   `mapstructure:"twilio"`
	Slack    slack.Config    `mapstructure:"slack"`
	REST     rest.Config     `mapstructure:"rest"`
}

// Channels combines all available channel clients
type Channels struct {
	Telegram Channel
	Twilio   Channel
	REST     Channel
	Slack    Channel
}

// Get channel by name
func (c *Channels) Get(chnl string) Channel {
	switch chnl {
	case c.Slack.String():
		return c.Slack
	case c.Telegram.String():
		return c.Telegram
	case c.Twilio.String():
		return c.Twilio
	case c.REST.String():
		return c.REST
	default:
		return nil
	}
}

// Channel interface implements a channel to send and receive messages on
type Channel interface {
	// ReceiveMessage from the channel
	ReceiveMessage(body []byte) (*message.Request, error)
	// ReceiveMessages from the channel. Starts a long running process, receives questions and sends them to the receiveChan
	ReceiveMessages(receiveChan chan message.Request)
	// SendMessage to the channel
	SendMessage(response *message.Response) error
	// ValidateCallback validates a callback to the channel
	ValidateCallback(r *http.Request) bool
	// String returns the channel's name
	String() string
}

// LoadConfig loads channels configuration from chn.yml
func LoadConfig(path string) (*Config, error) {
	config := viper.New()
	config.SetConfigName("chn")
	config.AddConfigPath(path)
	config.AutomaticEnv()
	config.SetEnvPrefix("CHATTO_CHN")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := config.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warn("File chn.yml not found, using only REST channel")
		default:
			return nil, err
		}
	}

	var channelsConfig Config
	if err := config.Unmarshal(&channelsConfig); err != nil {
		return nil, err
	}

	return &channelsConfig, nil
}

// New initializes all channels
func New(channelsConfig *Config) *Channels {
	chnls := Channels{}

	// REST
	chnls.REST = rest.New(channelsConfig.REST)

	// TELEGRAM
	if channelsConfig.Telegram != (telegram.Config{}) {
		chnls.Telegram = telegram.New(channelsConfig.Telegram)
	}

	// TWILIO
	if channelsConfig.Twilio != (twilio.Config{}) {
		chnls.Twilio = twilio.New(channelsConfig.Twilio)
	}

	// SLACK
	if channelsConfig.Slack != (slack.Config{}) {
		chnls.Slack = slack.New(channelsConfig.Slack)
	}

	return &chnls
}
