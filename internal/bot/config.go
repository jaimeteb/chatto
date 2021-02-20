package bot

import (
	"strings"

	"github.com/jaimeteb/chatto/internal/channels"
	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/extension"
	"github.com/jaimeteb/chatto/internal/fsm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config struct models the bot.yml configuration file
type Config struct {
	Name       string           `mapstructure:"bot_name"`
	Extensions extension.Config `mapstructure:"extensions"`
	Store      fsm.StoreConfig  `mapstructure:"store"`
	Port       int              `mapstructure:"port"`
	Path       string
}

// LoadConfig loads bot configuration from bot.yml
func LoadConfig(path string, port int) (*Config, error) {
	config := viper.New()
	config.SetConfigName("bot")
	config.AddConfigPath(path)
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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
		Store:  fsm.NewStore(botConfig.Store),
		Config: botConfig,
	}

	// Load Channels
	channelsConfig, err := channels.LoadConfig(botConfig.Path)
	if err != nil {
		return nil, err
	}
	b.Channels = channels.New(channelsConfig)

	// Load FSM
	fsmConfig, err := fsm.LoadConfig(botConfig.Path)
	if err != nil {
		return nil, err
	}
	b.Domain = fsm.New(fsmConfig)

	// Load Classifier
	classifConfig, err := clf.LoadConfig(botConfig.Path)
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

	log.Infof("My name is '%v'", b.Name)

	return b, nil
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
