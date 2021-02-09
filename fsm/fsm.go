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
type TransitionFunc func(m *FSM) (string, []Message)

// NewTransitionFunc generates a new transition function
func NewTransitionFunc(state int, extension string, message []Message) TransitionFunc {
	return func(m *FSM) (string, []Message) {
		(*m).State = state
		return extension, message
	}
}

// FSM models a Finite State Machine
type FSM struct {
	State int               `json:"state"`
	Slots map[string]string `json:"slots"`
}

// ExecuteCmd executes a command in the FSM
func (m *FSM) ExecuteCmd(cmd, txt string, db *DB) (answers []query.Answer, runExt string) {
	var transition TransitionFunc
	var tuple CmdStateTuple

	previousState := m.State

	tupleFromAny := CmdStateTuple{cmd, -1}
	tupleNormal := CmdStateTuple{cmd, m.State}
	tupleCmdAny := CmdStateTuple{"any", m.State}

	if db.TransitionTable[tupleFromAny] == nil {
		if db.TransitionTable[tupleCmdAny] == nil {
			transition = db.TransitionTable[tupleNormal] // There is no transition "From Any" with cmd, nor "Cmd Any"
			tuple = tupleNormal
		} else {
			transition = db.TransitionTable[tupleCmdAny] // There is a transition "Cmd Any"
			tuple = tupleCmdAny
		}
	} else {
		transition = db.TransitionTable[tupleFromAny] // There is a transition "From Any" with cmd
		tuple = tupleFromAny
	}

	slot := db.SlotTable[tuple]

	// Get slots
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

	// Get answers
	if cmd == "" {
		answers = append(answers, query.Answer{Text: db.DefaultMessages.Unsure}) // Threshold not met
	} else if transition == nil {
		answers = append(answers, query.Answer{Text: db.DefaultMessages.Unknown}) // Unknown transition
	} else {
		transition, message := transition(m)

		if strings.TrimSpace(transition) != "" {
			runExt = transition
		} else {
			for _, msg := range message {
				answers = append(answers, query.Answer{Text: msg.Text, Image: msg.Image})
			}
		}
	}

	log.Debugf("FSM | transitioned %v -> %v", previousState, m.State)

	return
}
