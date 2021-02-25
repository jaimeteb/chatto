package main

import (
	"log"

	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/query"
)

func greetFunc(req *extension.ExecuteCommandFuncRequest) (res *extension.ExecuteCommandFuncResponse) {
	return &extension.ExecuteCommandFuncResponse{
		FSM: req.FSM,
		Answers: []query.Answer{{
			Text:  "Hello Universe",
			Image: "https://i.imgur.com/pPdjh6x.jpg",
		}},
	}
}

var registeredFuncs = extension.RegisteredCommandFuncs{
	"any": greetFunc,
}

func main() {
	if err := extension.ServeREST(registeredFuncs); err != nil {
		log.Fatalln(err)
	}
}
