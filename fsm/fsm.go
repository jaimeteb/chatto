package fsm

import (
	"regexp"
	"strings"

	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// BaseDB contains the data required for a minimally functioning FSM
type BaseDB struct {
	StateTable      map[string]int `json:"state_table"`
	CommandList     []string       `json:"command_list"`
	DefaultMessages Defaults       `json:"default_messages"`
}

// DB contains BaseDB plus the functions required for a fully
// functioning FSM
type DB struct {
	BaseDB
	TransitionTable map[CmdStateTuple]TransitionFunc
	SlotTable       map[CmdStateTuple]Slot
}

// NoFuncs returns a DB without TransitionFunc items in order
// to serialize it for extensions
func (d *DB) NoFuncs() *BaseDB {
	return &BaseDB{
		StateTable:      d.StateTable,
		CommandList:     d.CommandList,
		DefaultMessages: d.DefaultMessages,
	}
}

// CmdStateTuple is a tuple of Command and State
type CmdStateTuple struct {
	Cmd   string
	State int
}

// TransitionFunc models a transition function
type TransitionFunc func(m *FSM) interface{}

// NewTransitionFunc generates a new transition function
func NewTransitionFunc(state int, r interface{}) TransitionFunc {
	return func(m *FSM) interface{} {
		(*m).State = state
		return r
	}
}

// FSM models a Finite State Machine
type FSM struct {
	State int               `json:"state"`
	Slots map[string]string `json:"slots"`
}

// ExecuteCmd executes a command in the FSM
func (m *FSM) ExecuteCmd(cmd, txt string, machineState *DB) (answers []query.Answer, runExt string) {
	var transition TransitionFunc
	var tuple CmdStateTuple

	previousState := m.State

	tupleFromAny := CmdStateTuple{cmd, -1}
	tupleNormal := CmdStateTuple{cmd, m.State}
	tupleCmdAny := CmdStateTuple{"any", m.State}

	if machineState.TransitionTable[tupleFromAny] == nil {
		if machineState.TransitionTable[tupleCmdAny] == nil {
			transition = machineState.TransitionTable[tupleNormal] // There is no transition "From Any" with cmd, nor "Cmd Any"
			tuple = tupleNormal
		} else {
			transition = machineState.TransitionTable[tupleCmdAny] // There is a transition "Cmd Any"
			tuple = tupleCmdAny
		}
	} else {
		transition = machineState.TransitionTable[tupleFromAny] // There is a transition "From Any" with cmd
		tuple = tupleFromAny
	}

	slot := machineState.SlotTable[tuple]
	if slot.Name != "" {
		switch slot.Mode {
		case "whole_text":
			m.Slots[slot.Name] = txt
		case "regex":
			if r, err := regexp.Compile(slot.Regex); err == nil {
				match := r.FindAllString(txt, 1)
				if len(match) > 0 {
					m.Slots[slot.Name] = match[0]
				}
			}
		}
	}
	// log.Debug(m.Slots)

	if cmd == "" {
		answers = append(answers, query.Answer{Text: machineState.DefaultMessages.Unsure}) // Threshold not met
	} else if transition == nil {
		answers = append(answers, query.Answer{Text: machineState.DefaultMessages.Unknown}) // Unknown transition
	} else {
		response := transition(m)
		switch r := response.(type) {
		case string:
			if strings.HasPrefix(r, "ext_") {
				runExt = r
			}
		}
	}

	log.Debugf("FSM | transitioned %v -> %v\n", previousState, m.State)

	return
}
