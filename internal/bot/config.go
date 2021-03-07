package bot

import (
	"strings"
	"time"

	"github.com/jaimeteb/chatto/internal/channels"
	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/extension"
	"github.com/jaimeteb/chatto/internal/fsm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ConversationConfig for the bot
type ConversationConfig struct {
	ReplyUnsure  bool `mapstructure:"reply_unsure"`
	ReplyUnknown bool `mapstructure:"reply_unknown"`
	ReplyError   bool `mapstructure:"reply_error"`
}

// Conversation settings for new and existing conversations
type Conversation struct {
	New      ConversationConfig `mapstructure:"new"`
	Existing ConversationConfig `mapstructure:"existing"`
}

// Auth is authorization for the bot API
type Auth struct {
	Token string `mapstructure:"token"`
}

// Config struct models the bot.yml configuration file
type Config struct {
	Name         string           `mapstructure:"bot_name"`
	Extensions   extension.Config `mapstructure:"extensions"`
	Store        fsm.StoreConfig  `mapstructure:"store"`
	Port         int              `mapstructure:"port"`
	Path         string
	Conversation Conversation `mapstructure:"conversation"`
	Auth         Auth         `mapstructure:"auth"`
}

// ShouldReplyUnsure depending on the conversational settings lets
// the bot know if it should reply with Unsure to the channel
func (c *Config) ShouldReplyUnsure(isExistingConversation bool) bool {
	if isExistingConversation {
		return c.Conversation.Existing.ReplyUnsure
	}

	return c.Conversation.New.ReplyUnsure
}

// ShouldReplyUnknown depending on the conversational settings lets
// the bot know if it should reply with Unknown to the channel
func (c *Config) ShouldReplyUnknown(isExistingConversation bool) bool {
	if isExistingConversation {
		return c.Conversation.Existing.ReplyUnknown
	}

	return c.Conversation.New.ReplyUnknown
}

// ShouldReplyError depending on the conversational settings lets
// the bot know if it should reply with Error to the channel
func (c *Config) ShouldReplyError(isExistingConversation bool) bool {
	if isExistingConversation {
		return c.Conversation.Existing.ReplyError
	}

	return c.Conversation.New.ReplyError
}

// LoadConfig loads bot configuration from bot.yml
func LoadConfig(path string, port int) (*Config, error) {
	config := viper.New()
	config.SetConfigName("bot")
	config.AddConfigPath(path)
	config.AutomaticEnv()
	config.SetEnvPrefix("CHATTO_BOT")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.SetDefault("conversation.new.reply_unsure", true)
	config.SetDefault("conversation.new.reply_unknown", true)
	config.SetDefault("conversation.new.reply_error", true)
	config.SetDefault("conversation.existing.reply_unsure", true)
	config.SetDefault("conversation.existing.reply_unknown", true)
	config.SetDefault("conversation.existing.reply_error", true)

	if err := config.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warn("File bot.yml not found, using default values")
		default:
			return nil, err
		}
	}

	var bc Config
	err := config.Unmarshal(&bc)
	if err != nil {
		return nil, err
	}

	bc.Path = path
	bc.Port = port

	return &bc, nil
}

// loadName loads the bot name from the configuration file
func loadName(name string) string {
	if name != "" {
		return name
	}

	return "botto"
}

// New initializes and returns a new Bot
func New(botConfig *Config) (*Bot, error) {
	b := &Bot{
		Name:   loadName(botConfig.Name),
		Store:  fsm.NewStore(&botConfig.Store),
		Config: botConfig,
	}

	// Load Channels
	channelsConfig, err := channels.LoadConfig(botConfig.Path)
	if err != nil {
		return nil, err
	}
	b.Channels = channels.New(channelsConfig)

	// Load FSM Domain
	fsmReloadChan := make(chan fsm.Config)
	fsmConfig, err := fsm.LoadConfig(botConfig.Path, fsmReloadChan)
	if err != nil {
		return nil, err
	}
	b.Domain = fsm.NewDomainFromConfig(fsmConfig)

	// Load Classifier
	classifReloadChan := make(chan clf.Config)
	classifConfig, err := clf.LoadConfig(botConfig.Path, classifReloadChan)
	if err != nil {
		return nil, err
	}
	b.Classifier = clf.New(classifConfig)

	// Load Extensions
	ext, err := extension.New(botConfig.Extensions)
	if err != nil {
		return nil, err
	}
	b.Extension = ext

	// Register HTTP handlers
	b.RegisterRoutes()

	// Reload the FSM Domain and CLF Classifier if the configs change
	receiveAndReload(b, fsmReloadChan, classifReloadChan)

	log.Infof("My name is '%v'", b.Name)

	return b, nil
}

func receiveAndReload(b *Bot, fsmReloadChan chan fsm.Config, classifReloadChan chan clf.Config) {
	go func() {
		for {
			select {
			case fsmConfig := <-fsmReloadChan:
				b.Domain = fsm.NewDomainFromConfig(&fsmConfig)
			case classifConfig := <-classifReloadChan:
				b.Classifier = clf.New(&classifConfig)
			default:
				time.Sleep(5 * time.Second)
			}
		}
	}()
}

// LOGO for Chatto
const LOGO = `
                           *******                          
                  *************************                 
             *********                *********             
          *******                           ******.         
        *****                                  ******       
      *****                                       *****     
    *****                                           *****   
   ****                                              .****  
  ****         ********,             *********         **** 
 ****       .******.******         ******.******       .****
 ****       ****       ****       ****       ****       ****
****                                                    ****
****                                                     ***
****                                                    ****
 ****                                                   ****
 ****                  ****       ****                 ****.
  ****                 *****     ****                  **** 
   ****.                 ***********                 *****  
    *****                                           ****    
      *****                                       *****     
        ******                                 ******       
          .******                          .******          
              *********               *********             
                  .***********************.                 
`
