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

// Load loads configuration from yaml
func Load(path *string) Config {
	config := viper.New()
	config.SetConfigName("fsm")
	config.AddConfigPath(*path)

	if err := config.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	var fsmConfig Config
	if err := config.Unmarshal(&fsmConfig); err != nil {
		log.Panic(err)
	}

	return fsmConfig
}

// Create initializes the FSM Domain from Config
func Create(path *string) *DB {
	config := Load(path)

	db := &DB{}

	stateTable := make(map[string]int)
	for i, state := range config.States {
		stateTable[state] = i
	}
	stateTable["any"] = -1 // Add state "any"

	transitionTable := make(map[CmdStateTuple]TransitionFunc, len(config.Functions))

	slotTable := make(map[CmdStateTuple]Slot, len(config.Functions))

	for _, function := range config.Functions {
		tuple := CmdStateTuple{
			Cmd:   function.Command,
			State: stateTable[function.Transition.From],
		}

		transitionTable[tuple] = NewTransitionFunc(
			stateTable[function.Transition.Into],
			function.Extension,
			function.Message,
		)

		if function.Slot != (Slot{}) {
			slotTable[tuple] = function.Slot
		}
	}

	db.StateTable = stateTable
	db.CommandList = config.Commands
	db.TransitionTable = transitionTable
	db.DefaultMessages = config.Defaults
	db.SlotTable = slotTable

	log.Info("Loaded states:")
	for state, i := range stateTable {
		log.Infof("%v - %v", i, state)
	}

	return db
}
