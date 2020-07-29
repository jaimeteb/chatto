package pkg

import (
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
	Tuple      Tuple  `yaml:"tuple"`
	Transition string `yaml:"transition"`
	Message    string `yaml:"message"`
}

// Tuple models a command state tuple in yaml
type Tuple struct {
	Command string `yaml:"command"`
	State   string `yaml:"state"`
}

// Domain models the final configuration of an FSM
type Domain struct {
	StateTable      map[string]int
	CommandList     []string
	TransitionTable map[CmdStateTupple]TransitionFunc
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
}

// ExecuteCmd executes a command in FSM
func (m *FSM) ExecuteCmd(cmd string, dom *Domain) string {
	// get function from transition table
	tupple := CmdStateTupple{strings.TrimSpace(cmd), m.State}
	trans := dom.TransitionTable[tupple]
	if trans == (TransitionFunc{}) {
		return dom.DefaultMessages["unknown"]
	}
	m.State = trans.Next
	return trans.Message
}

// LoadConfig loads configuration from yaml
func LoadConfig() Config {
	config := viper.New()
	config.SetConfigName("states")
	config.AddConfigPath(".")

	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}

	var botConfig Config
	if err := config.Unmarshal(&botConfig); err != nil {
		panic(err)
	}

	return botConfig
}

// LoadDomain loads a domain struct from loaded configuration
func LoadDomain() Domain {
	config := LoadConfig()
	var domain Domain

	stateTable := make(map[string]int)
	for i, state := range config.States {
		stateTable[state] = i
	}

	transitionTable := make(map[CmdStateTupple]TransitionFunc)
	for _, function := range config.Functions {
		tupple := CmdStateTupple{
			Cmd:   function.Tuple.Command,
			State: stateTable[function.Tuple.State],
		}
		transitionTable[tupple] = TransitionFunc{
			stateTable[function.Transition],
			function.Message,
		}
	}

	domain.StateTable = stateTable
	domain.CommandList = config.Commands
	domain.TransitionTable = transitionTable
	domain.DefaultMessages = config.Defaults

	return domain

	// machine := FSM{State: 0}
	// reader := bufio.NewReader(os.Stdin)
	// for {
	// 	// read command from stdin
	// 	cmd, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}

	// 	x := machine.ExecuteCmd(cmd, &domain)
	// 	fmt.Println(x)
	// }
}
