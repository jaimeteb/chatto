package main

import (
	"log"

	cmn "github.com/jaimeteb/chatto/common"
	"github.com/jaimeteb/chatto/fsm"
)

func greetFunc(req *fsm.Request) (res *fsm.Response) {
	return &fsm.Response{
		FSM: req.FSM,
		Res: cmn.Message{
			Text:  "Hello Universe",
			Image: "https://i.imgur.com/pPdjh6x.jpg",
		},
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
