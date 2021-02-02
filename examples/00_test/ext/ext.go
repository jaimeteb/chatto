package main

import (
	"log"

	"github.com/jaimeteb/chatto/ext"
	"github.com/jaimeteb/chatto/message"
)

func greetFunc(req *ext.Request) (res *ext.Response) {
	return &ext.Response{
		FSM: req.FSM,
		Res: message.Message{
			Text:  "Hello Universe",
			Image: "https://i.imgur.com/pPdjh6x.jpg",
		},
	}
}

var myExtMap = ext.ExtensionMap{
	"ext_any": greetFunc,
}

func main() {
	if err := ext.ServeExtensionREST(myExtMap); err != nil {
		log.Fatalln(err)
	}
}
