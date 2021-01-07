package fsm

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

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

// Config models the yaml configuration
type Config struct {
	States    []string   `yaml:"states"`
	Commands  []string   `yaml:"commands"`
	Functions []Function `yaml:"functions"`
	Defaults  Defaults   `yaml:"defaults"`
}

// Function models a function in yaml
type Function struct {
	Transition Transition  `yaml:"transition"`
	Command    string      `yaml:"command"`
	Slot       Slot        `yaml:"slot"`
	Message    interface{} `yaml:"message"`
}

// Transition models a state transition
type Transition struct {
	From string `yaml:"from"`
	Into string `yaml:"into"`
}

// Slot models a slot configuration
type Slot struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

// Defaults models the domain's default messages
type Defaults struct {
	Unknown string `yaml:"unknown"`
	Unsure  string `yaml:"unsure"`
	Error   string `yaml:"error"`
}

// Domain models the final configuration of an FSM
type Domain struct {
	StateTable      map[string]int
	CommandList     []string
	TransitionTable map[CmdStateTuple]TransitionFunc
	SlotTable       map[CmdStateTuple]Slot
	DefaultMessages Defaults
}

// DomainNoFuncs models the final configuration of an FSM without functions
// to be used in extensions
type DomainNoFuncs struct {
	StateTable      map[string]int
	CommandList     []string
	DefaultMessages Defaults
}

// CmdStateTuple is a tuple of Command and State
type CmdStateTuple struct {
	Cmd   string
	State int
}

// TransitionFunc models a transition function
type TransitionFunc func(m *FSM) interface{}

// FSM models a Finite State Machine
type FSM struct {
	State int               `json:"state"`
	Slots map[string]string `json:"slots"`
}

// NoFuncs returns a Domain without TransitionFunc items in order
// to serialize it for extensions
func (d *Domain) NoFuncs() *DomainNoFuncs {
	return &DomainNoFuncs{
		StateTable:      d.StateTable,
		CommandList:     d.CommandList,
		DefaultMessages: d.DefaultMessages,
	}
}

// NewTransitionFunc generates a new transition function
func NewTransitionFunc(s int, r interface{}) TransitionFunc {
	return func(m *FSM) interface{} {
		(*m).State = s
		return r
	}
}

// ExecuteCmd executes a command in FSM
func (m *FSM) ExecuteCmd(cmd, txt string, dom Domain, ext Extension) (response interface{}) {
	var trans TransitionFunc
	var tuple CmdStateTuple

	previousState := m.State

	tupleFromAny := CmdStateTuple{cmd, -1}
	tupleNormal := CmdStateTuple{cmd, m.State}
	tupleCmdAny := CmdStateTuple{"any", m.State}

	if dom.TransitionTable[tupleFromAny] == nil {
		if dom.TransitionTable[tupleCmdAny] == nil {
			trans = dom.TransitionTable[tupleNormal] // There is no transition "From Any" with cmd, nor "Cmd Any"
			tuple = tupleNormal
		} else {
			trans = dom.TransitionTable[tupleCmdAny] // There is a transition "Cmd Any"
			tuple = tupleCmdAny
		}
	} else {
		trans = dom.TransitionTable[tupleFromAny] // There is a transition "From Any" with cmd
		tuple = tupleFromAny
	}

	slot := dom.SlotTable[tuple]
	if slot.Name != "" {
		switch slot.Mode {
		case "whole_text":
			m.Slots[slot.Name] = txt
		}
	}
	// log.Println(m.Slots)

	if cmd == "" {
		response = dom.DefaultMessages.Unsure // Threshold not met
	} else if trans == nil {
		response = dom.DefaultMessages.Unknown // Unknown transition
	} else {
		response = trans(m)
		switch r := response.(type) {
		case string:
			if strings.HasPrefix(r, "ext_") && ext != nil {
				response = ext.RunExtFunc(r, txt, dom, m)
			}
		}
	}

	log.Debugf("FSM | transitioned %v -> %v\n", previousState, m.State)
	return
}

// Load loads configuration from yaml
func Load(path *string) Config {
	config := viper.New()
	config.SetConfigName("fsm")
	config.AddConfigPath(*path)

	if err := config.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	var botConfig Config
	if err := config.Unmarshal(&botConfig); err != nil {
		log.Panic(err)
	}

	return botConfig
}

// Create loads a domain struct from loaded configuration
func Create(path *string) Domain {
	config := Load(path)
	var domain Domain

	stateTable := make(map[string]int)
	for i, state := range config.States {
		stateTable[state] = i
	}
	stateTable["any"] = -1 // Add state "any"

	transitionTable := make(map[CmdStateTuple]TransitionFunc)
	slotTable := make(map[CmdStateTuple]Slot)
	for _, function := range config.Functions {
		tuple := CmdStateTuple{
			Cmd:   function.Command,
			State: stateTable[function.Transition.From],
		}
		transitionTable[tuple] = NewTransitionFunc(
			stateTable[function.Transition.Into],
			function.Message,
		)
		if function.Slot != (Slot{}) {
			slotTable[tuple] = function.Slot
		}
	}

	domain.StateTable = stateTable
	domain.CommandList = config.Commands
	domain.TransitionTable = transitionTable
	domain.DefaultMessages = config.Defaults
	domain.SlotTable = slotTable

	log.Info("Loaded states:")
	for state, i := range stateTable {
		log.Infof("%v\t%v\n", i, state)
	}

	return domain
}
