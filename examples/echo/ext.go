package main

import (
	"fmt"

	"github.com/jaimeteb/chatto/fsm"
)

// FuncMap maps function names to functions
type FuncMap map[string]func(*fsm.FSM) interface{}

func greetFunc(m *fsm.FSM) interface{} {
	return "Hello Universe"
}

func goodbyeFunc(m *fsm.FSM) interface{} {
	return "Goodbye Universe"
}

func sayNameAgeFunc(m *fsm.FSM) interface{} {
	return fmt.Sprintf("Your name is %v and you're %v years old", m.Slots["name"], m.Slots["age"])
}

func noFunc(*fsm.FSM) interface{} {
	return nil
}

// Ext is exported
var Ext = FuncMap{
	"greet":        greetFunc,
	"goodbye":      goodbyeFunc,
	"ext_name_age": sayNameAgeFunc,
}

// GetFunc gets a function from the function map
func (fm FuncMap) GetFunc(action string) func(*fsm.FSM) interface{} {
	if _, ok := fm[action]; ok {
		return fm[action]
	}
	return noFunc
}

// GetAllFuncs retreives all functions in function map
func (fm FuncMap) GetAllFuncs() []string {
	allFuncs := make([]string, 0)
	for funcName := range fm {
		allFuncs = append(allFuncs, funcName)
	}
	return allFuncs
}

func main() {}
