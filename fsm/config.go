package fsm

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config contains the states, commands, functions and
// default messages of the FSM
type Config struct {
	States    []string   `yaml:"states"`
	Commands  []string   `yaml:"commands"`
	Functions []Function `yaml:"functions"`
	Defaults  Defaults   `yaml:"defaults"`
}

// Function lists the transitions available for the FSM
type Function struct {
	Transition Transition `yaml:"transition"`
	Command    string     `yaml:"command"`
	Slot       Slot       `yaml:"slot"`
	Extension  string     `yaml:"extension"`
	Message    []Message  `yaml:"message"`
}

// Message that is sent when a transition is executed
type Message struct {
	Text  string `yaml:"text"`
	Image string `yaml:"image"`
}

// Transition describes the states of the transition
// (from one state into another) if the functions command
// is executed
type Transition struct {
	From string `yaml:"from"`
	Into string `yaml:"into"`
}

// Slot is used to save information from the user's input
type Slot struct {
	Name  string `yaml:"name"`
	Mode  string `yaml:"mode"`
	Regex string `yaml:"regex"`
}

// Defaults set the messages that will be returned when
// Unknown, Unsure or Error events happen during FSM execution
type Defaults struct {
	Unknown string `yaml:"unknown" json:"unknown"`
	Unsure  string `yaml:"unsure" json:"unsure"`
	Error   string `yaml:"error" json:"error"`
}

// LoadConfig loads the FSM configuration from yaml
func LoadConfig(path string) (*Config, error) {
	config := viper.New()
	config.SetConfigName("fsm")
	config.AddConfigPath(path)

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	var fsmConfig Config
	if err := config.Unmarshal(&fsmConfig); err != nil {
		return nil, err
	}

	return &fsmConfig, nil
}

// New initializes the FSM
func New(fsmConfig *Config) *Domain {
	fsmDomain := &Domain{}

	stateTable := make(map[string]int)
	for i, state := range fsmConfig.States {
		stateTable[state] = i
	}
	stateTable["any"] = -1 // Add state "any"

	transitionTable := make(map[CmdStateTuple]TransitionFunc, len(fsmConfig.Functions))

	slotTable := make(map[CmdStateTuple]Slot, len(fsmConfig.Functions))

	for n := range fsmConfig.Functions {
		tuple := CmdStateTuple{
			Cmd:   fsmConfig.Functions[n].Command,
			State: stateTable[fsmConfig.Functions[n].Transition.From],
		}

		transitionTable[tuple] = NewTransitionFunc(
			stateTable[fsmConfig.Functions[n].Transition.Into],
			fsmConfig.Functions[n].Extension,
			fsmConfig.Functions[n].Message,
		)

		if fsmConfig.Functions[n].Slot != (Slot{}) {
			slotTable[tuple] = fsmConfig.Functions[n].Slot
		}
	}

	fsmDomain.StateTable = stateTable
	fsmDomain.CommandList = fsmConfig.Commands
	fsmDomain.TransitionTable = transitionTable
	fsmDomain.DefaultMessages = fsmConfig.Defaults
	fsmDomain.SlotTable = slotTable

	log.Info("Loaded states:")
	for state, i := range stateTable {
		log.Infof("%v - %v", i, state)
	}

	return fsmDomain
}
