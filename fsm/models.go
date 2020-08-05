package fsm

import "log"

// Config models the yaml configuration
type Config struct {
	States    []string          `yaml:"states"`
	Commands  []string          `yaml:"commands"`
	Functions []Function        `yaml:"functions"`
	Defaults  map[string]string `yaml:"defaults"`
}

// Function models a function in yaml
type Function struct {
	Tuple      Tuple  `yaml:"tuple"`
	Transition string `yaml:"transition"`
	Message    string `yaml:"message"`
}

// Tuple models a command state tuple in yaml
type Tuple struct {
	Command string `yaml:"command"`
	State   string `yaml:"state"`
	Slot    string `yaml:"slot"`
}

// Domain models the final configuration of an FSM
type Domain struct {
	StateTable      map[string]int
	CommandList     []string
	TransitionTable map[CmdStateTupple]TransitionFunc
	SlotTable       map[CmdStateTupple]string
	DefaultMessages map[string]string
}

// CmdStateTupple is a tuple of Command and State
type CmdStateTupple struct {
	Cmd   string
	State int
}

// TransitionFunc models a transition function
type TransitionFunc Transition

// Transition models a state transition
type Transition struct {
	Next    int
	Message string
}

// FSM models a Finite State Machine
type FSM struct {
	State int
	Slots map[string]interface{}
}

// type Cmder interface {
// 	Original() string
// 	Predicted()
// }

// ExecuteCmd executes a command in FSM
func (m *FSM) ExecuteCmd(cmd, org string, dom Domain) string {
	// if cmd == "" {
	// 	return dom.DefaultMessages["unsure"]
	// }

	tupple := CmdStateTupple{cmd, m.State}
	trans := dom.TransitionTable[tupple]
	if trans == (TransitionFunc{}) {
		return dom.DefaultMessages["unknown"]
	}

	slot := dom.SlotTable[tupple]
	if slot != "" {
		m.Slots[slot] = org
	}
	log.Println(m.Slots)

	m.State = trans.Next
	return trans.Message
}
