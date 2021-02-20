package fsm

import (
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// Transition describes the states of the transition
// (from one state into another) if the functions command
// is executed
type Transition struct {
	From string `yaml:"from"`
	Into string `yaml:"into"`
}

// Function lists the transitions available for the FSM
type Function struct {
	Transition Transition `yaml:"transition"`
	Command    string     `yaml:"command"`
	Slot       Slot       `yaml:"slot"`
	Extension  string     `yaml:"extension"`
	Message    []Message  `yaml:"message"`
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

// Message that is sent when a transition is executed
type Message struct {
	Text  string `yaml:"text"`
	Image string `yaml:"image"`
}

// StateTable contains a mapping of state names to state ids
// TODO: Document how the StateTable works
type StateTable map[string]int

// NewStateTable initializes a new StateTable
func NewStateTable(states []string) StateTable {
	stateTable := make(map[string]int, len(states)+1)

	for id, state := range states {
		stateTable[state] = id
	}

	stateTable["any"] = -1 // Add state "any"

	return stateTable
}

// TransitionTable contains the mapping of state tuples to transition functions
// TODO: Document how the TransitionTable works
type TransitionTable map[CmdStateTuple]TransitionFunc

// NewTransitionTable initializes a new TransitionTable
func NewTransitionTable(functions []Function, stateTable StateTable) TransitionTable {
	transitionTable := make(TransitionTable, len(functions))

	for n := range functions {
		function := functions[n]

		cmdStateTuple := CmdStateTuple{
			Cmd:   function.Command,
			State: stateTable[function.Transition.From],
		}

		transitionTable[cmdStateTuple] = NewTransitionFunc(
			stateTable[function.Transition.Into],
			function.Extension,
			function.Message,
		)
	}

	return transitionTable
}

// SlotTable contains the mapping of state tuples to slots
// TODO: Document how the SlotTable works
type SlotTable map[CmdStateTuple]Slot

// NewSlotTable initializes a new SlotTable
func NewSlotTable(functions []Function, stateTable StateTable) SlotTable {
	slotTable := make(SlotTable, len(functions))

	for n := range functions {
		function := functions[n]

		cmdStateTuple := CmdStateTuple{
			Cmd:   function.Command,
			State: stateTable[function.Transition.From],
		}

		if function.Slot != (Slot{}) {
			slotTable[cmdStateTuple] = function.Slot
		}
	}

	return slotTable
}

// BaseDomain contains the data required for a minimally functioning FSM
type BaseDomain struct {
	StateTable      StateTable `json:"state_table"`
	CommandList     []string   `json:"command_list"`
	DefaultMessages Defaults   `json:"default_messages"`
}

// Domain contains BaseDomain plus the functions required for a fully
// functioning FSM
type Domain struct {
	BaseDomain
	TransitionTable TransitionTable
	SlotTable       SlotTable
}

// NewDomain initializes a new Domain
func NewDomain(commands, states []string, functions []Function, defaults Defaults) *Domain {
	fsmDomain := &Domain{}
	fsmDomain.CommandList = commands
	fsmDomain.DefaultMessages = defaults
	fsmDomain.StateTable = NewStateTable(states)
	fsmDomain.TransitionTable = NewTransitionTable(functions, fsmDomain.StateTable)
	fsmDomain.SlotTable = NewSlotTable(functions, fsmDomain.StateTable)

	return fsmDomain
}

// NoFuncs returns a Domain without TransitionFunc items in order
// to serialize it for extensions
func (d *Domain) NoFuncs() *BaseDomain {
	return &BaseDomain{
		StateTable:      d.StateTable,
		CommandList:     d.CommandList,
		DefaultMessages: d.DefaultMessages,
	}
}

// CmdStateTuple is a tuple of Command and State
// TODO: Document how the CmdStateTuple works
type CmdStateTuple struct {
	Cmd   string
	State int
}

// TransitionFunc performs a state transition for the FSM.
// TODO: Document how the TransitionFunc works
type TransitionFunc func(m *FSM) (extension string, messages []Message)

// NewTransitionFunc generates a new transition function
// that will transition the FSM into the specified state
// and return the extension and the states defined messages
func NewTransitionFunc(state int, extension string, messages []Message) TransitionFunc {
	return func(m *FSM) (string, []Message) {
		m.State = state
		return extension, messages
	}
}

// FSM models a Finite State Machine
type FSM struct {
	State int               `json:"state"`
	Slots map[string]string `json:"slots"`
}

// ExecuteCmd executes a command in the FSM
// TODO: Document what the ExecuteCmd does
func (m *FSM) ExecuteCmd(cmd, question string, fsmDomain *Domain) (answers []query.Answer, extension string) {
	cmdStateTuple, transitionFunc := m.SelectStateTransition(cmd, fsmDomain)

	// Save information from the user's input into the slot
	m.SetSlot(fsmDomain.SlotTable[cmdStateTuple], question)

	// Transition FSM state and get answers or extension to execute
	return m.TransitionState(cmd, transitionFunc, fsmDomain.DefaultMessages)
}

// SelectStateTransition based on the command provided
func (m *FSM) SelectStateTransition(cmd string, fsmDomain *Domain) (CmdStateTuple, TransitionFunc) {
	// fromAnyState means we can transition from any state
	fromAnyState := CmdStateTuple{cmd, -1}

	// normalState means we can transition from one existing state to another
	normalState := CmdStateTuple{cmd, m.State}

	// TODO: Whats the difference between fromAnyState and cmdAny?
	cmdAnyState := CmdStateTuple{"any", m.State}

	log.WithField("type", "fsm").Info(spew.Sprint(m))
	log.WithField("type", "domain").Info(spew.Sprint(fsmDomain))

	if fsmDomain.TransitionTable[fromAnyState] == nil {
		if fsmDomain.TransitionTable[cmdAnyState] == nil {
			// There is no transition "From Any" with cmd, nor "Cmd Any"
			return normalState, fsmDomain.TransitionTable[normalState]
		}

		// There is a transition "Cmd Any"
		return cmdAnyState, fsmDomain.TransitionTable[cmdAnyState]
	}

	// There is a transition "From Any" with cmd
	return fromAnyState, fsmDomain.TransitionTable[fromAnyState]
}

// SetSlot saves information from the user's input/question
func (m *FSM) SetSlot(slot Slot, question string) {
	slotName := strings.TrimSpace(slot.Name)
	slotRegex := strings.TrimSpace(slot.Regex)

	if slotName != "" {
		switch strings.TrimSpace(slot.Mode) {
		case "regex":
			if r, err := regexp.Compile(slotRegex); err == nil {
				match := r.FindAllString(question, 1)
				if len(match) > 0 {
					m.Slots[slotName] = match[0]
				}
			}
		default:
			// Use whole_text by default
			m.Slots[slotName] = question
		}
	}
}

// TransitionState FSM state and return the query answers or extension to execute.
func (m *FSM) TransitionState(cmd string, transitionFunc TransitionFunc, defaults Defaults) (answers []query.Answer, extension string) {
	// Threshold not met
	if strings.TrimSpace(cmd) == "" {
		return []query.Answer{{Text: defaults.Unsure}}, ""
	}

	// Unknown transition
	if transitionFunc == nil {
		return []query.Answer{{Text: defaults.Unknown}}, ""
	}

	// Execute transition
	extension, messages := transitionFunc(m)

	if strings.TrimSpace(extension) != "" {
		return nil, extension
	}

	for n := range messages {
		answers = append(answers, query.Answer{Text: messages[n].Text, Image: messages[n].Image})
	}

	return answers, ""
}
