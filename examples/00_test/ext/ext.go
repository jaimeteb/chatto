package main

import (
	"log"

	ext "github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/query"
)

func greetFunc(req *ext.Request) (res *ext.Response) {
	return &ext.Response{
		FSM: req.FSM,
		Answers: []query.Answer{
			{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			},
		},
	}
}

var myExtMap = ext.RegisteredFuncs{
	"ext_any": greetFunc,
}

func main() {
	if err := ext.ServeExtensionREST(myExtMap); err != nil {
		log.Fatalln(err)
	}
}
