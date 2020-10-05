package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jaimeteb/chatto/fsm"
)

func greetFunc(req *fsm.Request) (res *fsm.Response) {
	return &fsm.Response{
		FSM: req.FSM,
		Res: "Hello Universe",
	}
}

func goodbyeFunc(req *fsm.Request) (res *fsm.Response) {
	return &fsm.Response{
		FSM: req.FSM,
		Res: "Goodbye Universe",
	}
}

func sayNameAgeFunc(req *fsm.Request) (res *fsm.Response) {
	m := req.FSM

	name := m.Slots["name"]
	age := m.Slots["age"]

	var message string
	if _, err := strconv.Atoi(age); err != nil {
		message = fmt.Sprintf("Your name is %v", name)
	} else {
		message = fmt.Sprintf("Your name is %v and you're %v years old", name, age)
	}

	return &fsm.Response{
		FSM: req.FSM,
		Res: message,
	}
}

var myExtMap = fsm.ExtensionMap{
	"greet":        greetFunc,
	"goodbye":      goodbyeFunc,
	"ext_name_age": sayNameAgeFunc,
}

func main() {
	if err := fsm.ServeExtension(myExtMap); err != nil {
		log.Fatalln(err)
	}
}
