package main

import (
	"log"

	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/query"
)

func greetFunc(req *extension.ExecuteExtensionRequest) (res *extension.ExecuteExtensionResponse) {
	return &extension.ExecuteExtensionResponse{
		FSM: req.FSM,
		Answers: []query.Answer{{
			Text:  "Hello Universe",
			Image: "https://i.imgur.com/pPdjh6x.jpg",
		}},
	}
}

var registeredExtensions = extension.RegisteredExtensions{
	"any": greetFunc,
}

func main() {
	if err := extension.ServeREST(registeredExtensions); err != nil {
		log.Fatalln(err)
	}
}
