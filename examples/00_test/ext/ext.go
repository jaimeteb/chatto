package main

import (
	"log"

	cmn "github.com/jaimeteb/chatto/common"
	"github.com/jaimeteb/chatto/ext"
)

func greetFunc(req *ext.Request) (res *ext.Response) {
	return &ext.Response{
		FSM: req.FSM,
		Res: cmn.Message{
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
