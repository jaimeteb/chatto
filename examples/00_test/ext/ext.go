package main

import (
	"log"

	"github.com/jaimeteb/chatto/fsm"
)

func greetFunc(req *fsm.Request) (res *fsm.Response) {
	return &fsm.Response{
		FSM: req.FSM,
		Res: "Hello Universe",
	}
}

var myExtMap = fsm.ExtensionMap{
	"ext_any": greetFunc,
}

func main() {
	if err := fsm.ServeExtensionREST(myExtMap); err != nil {
		log.Fatalln(err)
	}
}
