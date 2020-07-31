package fsm

import (
	"github.com/spf13/viper"
)

// Load loads configuration from yaml
func Load() Config {
	config := viper.New()
	config.SetConfigName("fsm")
	config.AddConfigPath(".")
	config.AddConfigPath("config")

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
func Create() Domain {
	config := Load()
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
}
