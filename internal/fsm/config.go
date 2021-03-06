package fsm

import (
	"github.com/jaimeteb/chatto/fsm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config contains the states, commands, functions and
// default messages of the FSM
type Config struct {
	States      []string         `yaml:"states"`
	Commands    []string         `yaml:"commands"`
	Transitions []fsm.Transition `yaml:"transitions"`
	Defaults    fsm.Defaults     `yaml:"defaults"`
}

// LoadConfig loads the FSM configuration from yaml
func LoadConfig(path string) (*Config, error) {
	config := viper.New()
	config.SetConfigName("fsm")
	config.AddConfigPath(path)
	config.SetDefault("defaults.unknown", "Unknown command, try something different.")
	config.SetDefault("defaults.unsure", "Not sure I understood, try something different.")
	config.SetDefault("defaults.error", "There was an error, try again later.")

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	var fsmConfig Config
	if err := config.Unmarshal(&fsmConfig); err != nil {
		return nil, err
	}

	return &fsmConfig, nil
}

// NewDomainFromConfig initializes a FSM Domain from the FSM Config
func NewDomainFromConfig(fsmConfig *Config) *fsm.Domain {
	fsmDomain := fsm.NewDomain(fsmConfig.Commands, fsmConfig.States, fsmConfig.Transitions, fsmConfig.Defaults)

	log.Info("Loaded states:")
	for stateName, stateID := range fsmDomain.StateTable {
		log.Infof("%2d %v", stateID, stateName)
	}

	return fsmDomain
}
