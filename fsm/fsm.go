package fsm

import (
	"regexp"
	"strings"

	"github.com/jaimeteb/chatto/query"
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
	stateTable := make(StateTable, len(states)+1)

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
	return &d.BaseDomain
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

// NewFSM instantiates a new FSM
func NewFSM() *FSM {
	return &FSM{State: 0, Slots: make(map[string]string)}
}

// ExecuteCmd executes a command in the FSM
// TODO: Document what the ExecuteCmd does
func (m *FSM) ExecuteCmd(command, matchedText string, fsmDomain *Domain) (answers []query.Answer, extension string) {
	cmdStateTuple, transitionFunc := m.SelectStateTransition(command, fsmDomain)

	// Save information from the user's input into the slot
	m.SaveToSlot(matchedText, fsmDomain.SlotTable[cmdStateTuple])

	// Transition FSM state and get answers or extension to execute
	return m.TransitionState(command, transitionFunc, fsmDomain.DefaultMessages)
}

// SelectStateTransition based on the command provided
func (m *FSM) SelectStateTransition(command string, fsmDomain *Domain) (CmdStateTuple, TransitionFunc) {
	// fromAnyState means we can transition from any state
	fromAnyState := CmdStateTuple{command, -1}

	// cmdAnyState transition between any two states
	cmdAnyState := CmdStateTuple{"any", m.State}

	// normalState transition from one existing state to another
	normalState := CmdStateTuple{command, m.State}

	// Special state any can go from any state into another
	if fsmDomain.TransitionTable[fromAnyState] != nil {
		return fromAnyState, fsmDomain.TransitionTable[fromAnyState]
	}

	// Special command any is used to transition between two states,
	// regardless of the command predicted. Useful for taking in any
	// user input for searches
	if fsmDomain.TransitionTable[cmdAnyState] != nil {
		return cmdAnyState, fsmDomain.TransitionTable[cmdAnyState]
	}

	return normalState, fsmDomain.TransitionTable[normalState]
}

// SaveToSlot saves information from the user's input/question
func (m *FSM) SaveToSlot(matchedText string, slot Slot) {
	slotName := strings.TrimSpace(slot.Name)
	slotRegex := strings.TrimSpace(slot.Regex)

	if slotName != "" {
		switch strings.TrimSpace(slot.Mode) {
		case "regex":
			if r, err := regexp.Compile(slotRegex); err == nil {
				match := r.FindAllString(matchedText, 1)
				if len(match) > 0 {
					m.Slots[slotName] = match[0]
				}
			}
		default:
			// Use whole_text by default
			m.Slots[slotName] = matchedText
		}
	}
}

// TransitionState FSM state and return the query answers or extension to execute.
func (m *FSM) TransitionState(command string, transitionFunc TransitionFunc, defaults Defaults) (answers []query.Answer, extension string) {
	// Function command was not found in the matchedText
	if strings.TrimSpace(command) == "" {
		return []query.Answer{{Text: defaults.Unsure}}, ""
	}

	// Function command was found but state transition is unknown or not valid
	if transitionFunc == nil {
		return []query.Answer{{Text: defaults.Unknown}}, ""
	}

	// Execute transition
	extension, messages := transitionFunc(m)

	// Tell the bot to execute an extension to get the answer
	if strings.TrimSpace(extension) != "" {
		return nil, extension
	}

	for n := range messages {
		answers = append(answers, query.Answer{Text: messages[n].Text, Image: messages[n].Image})
	}

	return answers, ""
}
