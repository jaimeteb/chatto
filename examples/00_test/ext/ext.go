package main

import (
	"github.com/jaimeteb/chatto/fsm"
)

func greetFunc(m *fsm.FSM, dom *fsm.Domain, txt string) interface{} {
	return "Hello Universe"
}

// Ext is exported
var Ext = fsm.FuncMap{
	"ext_any": greetFunc,
}

func main() {}
