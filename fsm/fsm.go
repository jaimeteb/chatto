package fsm

import (
	"regexp"
	"strings"

	"github.com/jaimeteb/chatto/query"
)

const (
	// StateInitial is the first state the FSM enters upon
	// initialization of a new conversation and when ending
	// an existing conversation
	StateInitial = 0
	// StateAny allows transitioning from any state
	StateAny = -1
)

// Transition lists the transitions available for the FSM
// Describes the states of the transition
// (from one state into another) if the functions command
// is executed
type Transition struct {
	From      []string `yaml:"from"`
	Into      string   `yaml:"into"`
	Command   string   `yaml:"command"`
	Slot      Slot     `yaml:"slot"`
	Extension string   `yaml:"extension"`
	Answers   []Answer `yaml:"answers"`
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

// Answer that is sent when a transition is executed
type Answer struct {
	Text  string `yaml:"text"`
	Image string `yaml:"image"`
}

// StateTable contains a mapping of state names to state ids
// TODO: Document how the StateTable works
type StateTable map[string]int

// NewStateTable initializes a new StateTable
func NewStateTable(transitions []Transition) StateTable {
	stateTableDefaultSize := 2

	stateTable := make(StateTable, len(transitions)+stateTableDefaultSize)

	stateTable["any"] = -1    // Add state "any" ID
	stateTable["initial"] = 0 // Add state "initial" ID

	// Starting state ID
	stateID := 1

	for n := range transitions {
		state := strings.TrimSpace(transitions[n].Into)

		// Do not add duplicate states
		if _, ok := stateTable[state]; ok {
			continue
		}

		// Set state name to id mapping
		stateTable[state] = stateID

		// Increment state ID
		stateID++
	}

	return stateTable
}

// TransitionTable contains the mapping of state tuples to transition functions
// TODO: Document how the TransitionTable works
type TransitionTable map[CmdStateTuple]TransitionFunc

// NewTransitionTable initializes a new TransitionTable
func NewTransitionTable(transitions []Transition, stateTable StateTable) TransitionTable {
	transitionTable := make(TransitionTable, len(transitions))

	for n := range transitions {
		transition := transitions[n]

		for _, from := range transition.From {
			cmdStateTuple := CmdStateTuple{
				Cmd:   transition.Command,
				State: stateTable[from],
			}

			transitionTable[cmdStateTuple] = NewTransitionFunc(
				stateTable[transition.Into],
				transition.Extension,
				transition.Answers,
			)
		}
	}

	return transitionTable
}

// SlotTable contains the mapping of state tuples to slots
// TODO: Document how the SlotTable works
type SlotTable map[CmdStateTuple]Slot

// NewSlotTable initializes a new SlotTable
func NewSlotTable(transitions []Transition, stateTable StateTable) SlotTable {
	slotTable := make(SlotTable, len(transitions))

	for n := range transitions {
		transition := transitions[n]

		for _, from := range transition.From {
			cmdStateTuple := CmdStateTuple{
				Cmd:   transition.Command,
				State: stateTable[from],
			}

			if transition.Slot != (Slot{}) {
				slotTable[cmdStateTuple] = transition.Slot
			}
		}
	}

	return slotTable
}

// BaseDomain contains the data required for a minimally functioning FSM
type BaseDomain struct {
	StateTable      StateTable `json:"state_table"`
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
func NewDomain(transitions []Transition, defaults Defaults) *Domain {
	fsmDomain := &Domain{}
	fsmDomain.DefaultMessages = defaults
	fsmDomain.StateTable = NewStateTable(transitions)
	fsmDomain.TransitionTable = NewTransitionTable(transitions, fsmDomain.StateTable)
	fsmDomain.SlotTable = NewSlotTable(transitions, fsmDomain.StateTable)

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
type TransitionFunc func(m *FSM) (extension string, answers []Answer)

// NewTransitionFunc generates a new transition function
// that will transition the FSM into the specified state
// and return the extension and the states defined messages
func NewTransitionFunc(state int, extension string, answers []Answer) TransitionFunc {
	return func(m *FSM) (string, []Answer) {
		m.State = state
		return extension, answers
	}
}

// FSM models a Finite State Machine
type FSM struct {
	State int               `json:"state"`
	Slots map[string]string `json:"slots"`
}

// NewFSM instantiates a new FSM
func NewFSM() *FSM {
	return &FSM{State: StateInitial, Slots: make(map[string]string)}
}

// ExecuteCmd executes a state transition in the FSM based on
// the function command provided and if configured will save
// the classified text to a slot
func (m *FSM) ExecuteCmd(command, classifiedText string, fsmDomain *Domain) (answers []query.Answer, extension string, err error) {
	// Function command was not found by the classifier
	if strings.TrimSpace(command) == "" {
		if fsmDomain.DefaultMessages.Unsure == "" {
			return nil, "", nil
		}

		return nil, "", &ErrUnsureCommand{Msg: fsmDomain.DefaultMessages.Unsure}
	}

	cmdStateTuple, transitionFunc := m.SelectStateTransition(command, fsmDomain)

	// Save information from the user's input into the slot
	m.SaveToSlot(classifiedText, fsmDomain.SlotTable[cmdStateTuple])

	// Transition FSM state and get answers or extension to execute
	return m.TransitionState(transitionFunc, fsmDomain.DefaultMessages)
}

// SelectStateTransition based on the command provided
func (m *FSM) SelectStateTransition(command string, fsmDomain *Domain) (CmdStateTuple, TransitionFunc) {
	// fromAnyState means we can transition from any state
	fromAnyState := CmdStateTuple{command, StateAny}

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
func (m *FSM) SaveToSlot(classifiedText string, slot Slot) {
	slotName := strings.TrimSpace(slot.Name)
	slotRegex := strings.TrimSpace(slot.Regex)

	if slotName != "" {
		switch strings.TrimSpace(slot.Mode) {
		case "regex":
			if r, err := regexp.Compile(slotRegex); err == nil {
				match := r.FindAllString(classifiedText, 1)
				if len(match) > 0 {
					m.Slots[slotName] = match[0]
				}
			}
		default:
			// Use whole_text by default
			m.Slots[slotName] = classifiedText
		}
	}
}

// TransitionState FSM state and return the query answers or extension to execute.
func (m *FSM) TransitionState(transitionFunc TransitionFunc, defaults Defaults) (answers []query.Answer, extension string, err error) {
	// Function command was found by the classifier but state transition is unknown or not valid
	if transitionFunc == nil {
		if defaults.Unknown == "" {
			return nil, "", nil
		}

		return nil, "", &ErrUnknownCommand{Msg: defaults.Unknown}
	}

	// Execute transition
	extension, messages := transitionFunc(m)

	// Tell the bot to execute an extension to get the answer
	if strings.TrimSpace(extension) != "" {
		return nil, extension, nil
	}

	for n := range messages {
		answers = append(answers, query.Answer{Text: messages[n].Text, Image: messages[n].Image})
	}

	return answers, "", nil
}

// ErrUnsureCommand is returned by the FSM when no function
// command was found by the classifier
type ErrUnsureCommand struct {
	Msg string
}

// Error returns the ErrUnsureCommand error message
func (e *ErrUnsureCommand) Error() string {
	return e.Msg
}

// ErrUnknownCommand is returned by the FSM when a function
// command was found by the classifier but the state
// transition is unknown or not valid
type ErrUnknownCommand struct {
	Msg string
}

// Error returns the ErrUnknownCommand error message
func (e *ErrUnknownCommand) Error() string {
	return e.Msg
}
