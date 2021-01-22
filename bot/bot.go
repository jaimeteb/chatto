package bot

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jaimeteb/chatto/clf"
	cmn "github.com/jaimeteb/chatto/common"
	"github.com/jaimeteb/chatto/ext"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/spf13/viper"
)

func init() {
	lvl := os.Getenv("LOG_LEVEL")
	switch lvl {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

// TelegramMessageIn models a telegram incoming message
type TelegramMessageIn struct {
	UpdateID int                    `json:"update_id"`
	Message  TelegramMessageInInner `json:"message"`
}

// TelegramMessageInInner models a telegram incoming message inner struct
type TelegramMessageInInner struct {
	MessageID int                        `json:"message_id"`
	From      TelegramMessageInInnerFrom `json:"from"`
	Date      int                        `json:"date"`
	Text      string                     `json:"text"`
}

// TelegramMessageInInnerFrom models a telegram incoming message inner struct
type TelegramMessageInInnerFrom struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

// TwilioMessageIn models an incoming Twilio message
type TwilioMessageIn struct {
	From             string `form:"From"`
	Body             string `form:"Body"`
	To               string `form:"To"`
	MediaURL         string `form:"MediaUrl"`
	MediaContentType string `form:"MediaContentType"`
	MessageSid       string `form:"MessageSid"`
	SmsStatus        string `form:"SmsStatus"`
	AccountSid       string `form:"AccountSid"`
	Sid              string `form:"Sid"`
	SmsSid           string `form:"SmsSid"`
	SmsMessageSid    string `form:"SmsMessageSid"`
	NumMedia         int    `form:"NumMedia"`
	NumSegments      int    `form:"NumSegments"`
	APIVersion       string `form:"ApiVersion"`
}

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	Name       string
	Machines   fsm.StoreFSM
	Domain     fsm.Domain
	Classifier clf.Classifier
	Extension  ext.Extension
	Clients    Clients
}

// Prediction models a classifier prediction and its orignal string
type Prediction struct {
	Original    string  `json:"original"`
	Predicted   string  `json:"predicted"`
	Probability float64 `json:"probability"`
}

// Config struct models the bot.yml configuration file
type Config struct {
	Name       string               `mapstructure:"bot_name"`
	Extensions ext.ExtensionsConfig `mapstructure:"extensions"`
	Store      fsm.StoreConfig      `mapstructure:"store"`
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b Bot) Answer(mess cmn.Message) interface{} {
	if !b.Machines.Exists(mess.Sender) {
		b.Machines.Set(
			mess.Sender,
			&fsm.FSM{
				State: 0,
				Slots: make(map[string]string),
			},
		)
	}

	inputMessage := mess.Text
	cmd, _ := b.Classifier.Predict(inputMessage)

	m := b.Machines.Get(mess.Sender)
	resp, runExt := m.ExecuteCmd(cmd, inputMessage, b.Domain)
	if runExt != "" && b.Extension != nil {
		resp = b.Extension.RunExtFunc(runExt, inputMessage, b.Domain, m)
	}
	b.Machines.Set(mess.Sender, m)

	return resp
}

// LoadBotConfig loads bot configuration from bot.yml
func LoadBotConfig(path *string) Config {
	config := viper.New()
	config.SetConfigName("bot")
	config.AddConfigPath(*path)
	config.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)

	if err := config.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warn("File bot.yml not found, using default values")
		default:
			log.Warn(err)
		}
		return Config{}
	}

	var bc Config
	config.Unmarshal(&bc)

	return bc
}

// LoadName loads the bot name from the configuration file
func LoadName(bcName string) (name string) {
	name = "botto"
	if bcName != "" {
		name = bcName
	}
	log.Infof("My name is '%v'\n", name)
	return
}

// LoadBot loads all configurations and returns a Bot
func LoadBot(path *string) Bot {
	bc := LoadBotConfig(path)

	// Load Name
	name := LoadName(bc.Name)
	// Load Domain
	domain := fsm.Create(path)
	// Load Classifier
	classifier := clf.Create(path)
	// Load Extensions
	extension := ext.LoadExtensions(bc.Extensions)
	// Load clients
	clients := LoadClients(path)
	// Load Store
	machines := fsm.LoadStore(bc.Store)

	return Bot{name, machines, domain, classifier, extension, clients}
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
