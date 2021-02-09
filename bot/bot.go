package bot

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/jaimeteb/chatto/channels"
	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
	"github.com/spf13/viper"
)

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	Name       string
	Machines   fsm.StoreFSM
	DB         *fsm.DB
	Classifier clf.Classifier
	Extension  extension.Extension
	Channels   *channels.Channels
}

// Prediction models a classifier prediction and its original string
type Prediction struct {
	Original    string  `json:"original"`
	Predicted   string  `json:"predicted"`
	Probability float64 `json:"probability"`
}

// Config struct models the bot.yml configuration file
type Config struct {
	Name       string           `mapstructure:"bot_name"`
	Extensions extension.Config `mapstructure:"extensions"`
	Store      fsm.StoreConfig  `mapstructure:"store"`
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b Bot) Answer(question *query.Question) ([]query.Answer, error) {
	if !b.Machines.Exists(question.Sender) {
		b.Machines.Set(
			question.Sender,
			&fsm.FSM{
				State: 0,
				Slots: make(map[string]string),
			},
		)
	}

	cmd, _ := b.Classifier.Predict(question.Text)

	machine := b.Machines.Get(question.Sender)

	reply, runExt := machine.ExecuteCmd(cmd, question.Text, b.DB)

	var err error
	if runExt != "" && b.Extension != nil {
		reply, err = b.Extension.RunExtFunc(question, runExt, b.DB, machine)
		if err != nil {
			return nil, err
		}
	}

	b.Machines.Set(question.Sender, machine)

	return reply, nil
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
	log.Infof("My name is '%v'", name)
	return
}

// LoadBot loads all configurations and returns a Bot
func LoadBot(path *string) (*Bot, error) {
	bc := LoadBotConfig(path)

	// Load Name
	name := LoadName(bc.Name)

	// Load DB
	db := fsm.Create(path)

	// Load Classifier
	classifier := clf.Create(path)

	// Load Extensions
	extension, err := extension.LoadExtensions(bc.Extensions)
	if err != nil {
		return nil, err
	}

	// Load channels
	chnls := channels.Load(path)

	// Load Store
	machines := fsm.LoadStore(bc.Store)

	return &Bot{name, machines, db, classifier, extension, chnls}, nil
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
