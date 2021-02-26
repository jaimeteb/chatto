package main

import (
	"log"

	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/query"
)

func greetFunc(req *extension.Request) (res *extension.Response) {
	return &extension.Response{
		FSM: req.FSM,
		Answers: []query.Answer{{
			Text:  "Hello Universe",
			Image: "https://i.imgur.com/pPdjh6x.jpg",
		}},
	}
}

var RegisteredCommandFuncs = extension.RegisteredCommandFuncs{
	"any": greetFunc,
}

func main() {
	if err := extension.ServeREST(RegisteredCommandFuncs); err != nil {
		log.Fatalln(err)
	}
}
