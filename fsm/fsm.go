package fsm

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config models the yaml configuration
type Config struct {
	States    []string          `yaml:"states"`
	Commands  []string          `yaml:"commands"`
	Functions []Function        `yaml:"functions"`
	Defaults  map[string]string `yaml:"defaults"`
}

// Function models a function in yaml
type Function struct {
	Transition Transition `yaml:"transition"`
	Command    string     `yaml:"command"`
	Slot       Slot       `yaml:"slot"`
	Message    string     `yaml:"message"`
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

// Domain models the final configuration of an FSM
type Domain struct {
	StateTable      map[string]int
	CommandList     []string
	TransitionTable map[CmdStateTuple]TransitionFunc
	SlotTable       map[CmdStateTuple]Slot
	DefaultMessages map[string]string
}

// CmdStateTuple is a tuple of Command and State
type CmdStateTuple struct {
	Cmd   string
	State int
}

// TransitionFunc models a transition function
type TransitionFunc func(m *FSM) string

// FSM models a Finite State Machine
type FSM struct {
	State int                    `json:"state"`
	Slots map[string]interface{} `json:"slots"`
}

// NewTransitionFunc generates a new transition function
func NewTransitionFunc(s int, r string) TransitionFunc {
	return func(m *FSM) string {
		(*m).State = s
		return r
	}
}

// ExecuteCmd executes a command in FSM
func (m *FSM) ExecuteCmd(cmd, txt string, dom Domain, ext Extension) (response string) {
	//// if cmd == "" {
	//// 	return dom.DefaultMessages["unsure"]
	//// }

	var trans TransitionFunc
	var tuple CmdStateTuple

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

	if trans == nil {
		response = dom.DefaultMessages["unknown"]
	} else {
		response = trans(m)
		if strings.HasPrefix(response, "ext_") {
			// extFunc := ext.GetFunc(response)
			// response = extFunc(m, &dom, txt).(string) // response = fmt.Sprintf("%v", extFunc(m))
		}
	}

	return
}

// Load loads configuration from yaml
func Load(path *string) Config {
	config := viper.New()
	config.SetConfigName("fsm")
	config.AddConfigPath(*path)

	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	var botConfig Config
	if err := config.Unmarshal(&botConfig); err != nil {
		panic(err)
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

	log.Println("Loaded states:")
	for state, i := range stateTable {
		log.Printf("%v\t%v\n", i, state)
	}

	return domain
}
