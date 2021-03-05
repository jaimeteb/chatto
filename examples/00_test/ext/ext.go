package main

import (
	"log"

	ext "github.com/jaimeteb/chatto/extension"
)

func greetFunc(req *ext.ExecuteCommandFuncRequest) (res *ext.ExecuteCommandFuncResponse) {
	return req.NewExecuteCommandFuncResponse(
		ext.WithTextAnswer("Hello Universe"),
	)
}

var registeredCommandFuncs = ext.RegisteredCommandFuncs{
	"any": greetFunc,
}

func main() {
	if err := ext.ServeREST(registeredCommandFuncs); err != nil {
		log.Fatalln(err)
	}
}
