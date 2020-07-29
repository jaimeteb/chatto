package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	States    []string   `yaml:"states"`
	Commands  []string   `yaml:"commands"`
	Functions []Function `yaml:"functions"`
}

type Function struct {
	Tuple      Tuple  `yaml:"tuple"`
	Transition string `yaml:"transition"`
	Message    string `yaml:"message"`
}

type Tuple struct {
	Command string `yaml:"command"`
	State   string `yaml:"state"`
}

type Domain struct {
	StateTable      map[string]int
	CommandList     []string
	TransitionTable map[CmdStateTupple]TransitionFunc
	DefaultMessage  string
}

type CmdStateTupple struct {
	Cmd   string
	State int
}

// type TransitionFunc func(state *int) string
type TransitionFunc Transition

type Transition struct {
	Next    int
	Message string
}

type FSM struct {
	State int
}

func (m *FSM) ExecuteCmd(cmd string, dom *Domain) string {
	// get function from transition table
	tupple := CmdStateTupple{strings.TrimSpace(cmd), m.State}
	trans := dom.TransitionTable[tupple]
	if trans == (TransitionFunc{}) {
		return dom.DefaultMessage
	}
	m.State = trans.Next
	return trans.Message
}

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

func main() {
	config := LoadConfig()

	stateTable := make(map[string]int)
	for i, state := range config.States {
		stateTable[state] = i
	}

	commandList := make([]string, 0)
	for _, command := range config.Commands {
		commandList = append(commandList, command)
	}

	transitionTable := make(map[CmdStateTupple]TransitionFunc)
	var defaultMessage string
	for _, function := range config.Functions {
		if function.Transition != "" {
			tupple := CmdStateTupple{
				Cmd:   function.Tuple.Command,
				State: stateTable[function.Tuple.State],
			}
			transitionTable[tupple] = TransitionFunc{
				stateTable[function.Transition],
				function.Message,
			}
		} else {
			defaultMessage = function.Message
		}
	}

	fmt.Println("Config:")
	fmt.Println(stateTable)
	fmt.Println(commandList)
	fmt.Println(transitionTable)
	fmt.Println(defaultMessage)

	domain := Domain{
		StateTable:      stateTable,
		CommandList:     commandList,
		TransitionTable: transitionTable,
		DefaultMessage:  defaultMessage,
	}

	machine := FSM{State: 0}
	reader := bufio.NewReader(os.Stdin)
	for {
		// read command from stdin
		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}

		x := machine.ExecuteCmd(cmd, &domain)
		fmt.Println(x)
	}
}
