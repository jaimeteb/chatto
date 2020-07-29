package pkg

import (
	"strings"
)

// State the FSM state for turnstile
type State uint32

const (
	// Initial is the initial state
	Initial State = iota
	// AskMood asks for mood
	AskMood
	// SayGood replies good
	SayGood
	// SayBad replies bad
	SayBad
	// SayBye is the final state
	SayBye
)

const (
	// CmdGreet command greet
	CmdGreet = "greet"
	// CmdGood command good
	CmdGood = "good"
	// CmdBad command bad
	CmdBad = "bad"
	// CmdYes command yes
	CmdYes = "yes"
	// CmdNo command no
	CmdNo = "no"
)

// FSM the finite state machine
type FSM struct {
	State State
}

// ExecuteCmd execute command
func (p *FSM) ExecuteCmd(cmd string) string {
	// get function from transition table
	tupple := CmdStateTupple{strings.TrimSpace(cmd), p.State}
	if f := StateTransitionTable[tupple]; f == nil {
		return "unknown command, try again please"
	} else {
		return f(&p.State)
	}
}

// CmdStateTupple tupple for state-command combination
type CmdStateTupple struct {
	Cmd   string
	State State
}

// TransitionFunc transition function
type TransitionFunc func(state *State) string

// StateTransitionTable trsition table
var StateTransitionTable = map[CmdStateTupple]TransitionFunc{
	{CmdGreet, Initial}: func(state *State) string {
		*state = AskMood
		return "Hello! How are you?"
	},
	{CmdGood, AskMood}: func(state *State) string {
		*state = Initial
		return "Great! :)"
	},
	{CmdBad, AskMood}: func(state *State) string {
		*state = SayBad
		return "Oh don't be sad :(\nDid that help?"
	},
	{CmdYes, SayBad}: func(state *State) string {
		*state = Initial
		return "I'm glad! :)"
	},
	{CmdNo, SayBad}: func(state *State) string {
		*state = Initial
		return "Oh I'm sorry :("
	},
}

// func prompt(s State) {
// 	m := map[State]string{
// 		Initial:   "Initial",
// 		AskMood:   "AskMood",
// 		SayGood:   "SayGood",
// 		SayBad:    "SayBad",
// 		SayBye:    "SayBye",
// 	}
// 	fmt.Printf("current state is [%s], please input command [greet|good|bad|yes|no]\n", m[s])
// }

// func main() {
// 	machine := &FSM{State: Initial}
// 	// prompt(machine.State)
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		// read command from stdin
// 		cmd, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Fatalln(err)
// 		}

// 		machine.ExecuteCmd(cmd)
// 	}
// }
