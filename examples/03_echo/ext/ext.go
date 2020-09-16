package main

import (
	"fmt"
	"strconv"

	"github.com/jaimeteb/chatto/fsm"
)

func greetFunc(m *fsm.FSM, dom *fsm.Domain, txt string) interface{} {
	return "Hello Universe"
}

func goodbyeFunc(m *fsm.FSM, dom *fsm.Domain, txt string) interface{} {
	return "Goodbye Universe"
}

func sayNameAgeFunc(m *fsm.FSM, dom *fsm.Domain, txt string) interface{} {
	name := m.Slots["name"].(string)
	age := m.Slots["age"].(string)

	if _, err := strconv.Atoi(age); err != nil {
		return fmt.Sprintf("Your name is %v", name)
	}

	return fmt.Sprintf("Your name is %v and you're %v years old", name, age)
}

// Ext is exported
var Ext = fsm.FuncMap{
	"greet":        greetFunc,
	"goodbye":      goodbyeFunc,
	"ext_name_age": sayNameAgeFunc,
}

func main() {}
