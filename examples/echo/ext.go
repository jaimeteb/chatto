package main

import (
	"fmt"

	"github.com/jaimeteb/chatto/fsm"
)

// FuncMap maps function names to functions
type FuncMap map[string]func(*fsm.FSM)

func greetFunc(m *fsm.FSM) {
	fmt.Println("Hello Universe")
}

func goodbyeFunc(m *fsm.FSM) {
	fmt.Println("Goodbye Universe")
}

func sayNameAgeFunc(m *fsm.FSM) {
	fmt.Printf("Your name is %v and you're %v years old\n", m.Slots["name"], m.Slots["name"])
}

// Ext is exported
var Ext = FuncMap{
	"greet":        greetFunc,
	"goodbye":      goodbyeFunc,
	"ext_name_age": sayNameAgeFunc,
}

// GetFunc gets a function from the function map
func (fm FuncMap) GetFunc(action string) func(*fsm.FSM) {
	if _, ok := fm[action]; ok {
		return fm[action]
	}
	return func(*fsm.FSM) {}
}

func main() {}
