package main

import (
	"fmt"

	"github.com/jaimeteb/chatto/fsm"
)

func greetFunc(m *fsm.FSM) interface{} {
	return "Hello Universe"
}

func goodbyeFunc(m *fsm.FSM) interface{} {
	return "Goodbye Universe"
}

func sayNameAgeFunc(m *fsm.FSM) interface{} {
	return fmt.Sprintf("Your name is %v and you're %v years old", m.Slots["name"], m.Slots["age"])
}

// Ext is exported
var Ext = fsm.FuncMap{
	"greet":        greetFunc,
	"goodbye":      goodbyeFunc,
	"ext_name_age": sayNameAgeFunc,
}

func main() {}
