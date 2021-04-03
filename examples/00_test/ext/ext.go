package main

import (
	"log"

	"github.com/jaimeteb/chatto/extensions"
	"github.com/jaimeteb/chatto/query"
)

func greetFunc(req *extensions.ExecuteExtensionRequest) (res *extensions.ExecuteExtensionResponse) {
	return &extensions.ExecuteExtensionResponse{
		FSM: req.FSM,
		Answers: []query.Answer{{
			Text:  "Hello Universe",
			Image: "https://i.imgur.com/pPdjh6x.jpg",
		}},
	}
}

var registeredExtensions = extensions.Registered{
	"any": greetFunc,
}

func main() {
	if err := extensions.ServeREST(registeredExtensions); err != nil {
		log.Fatalln(err)
	}
}
