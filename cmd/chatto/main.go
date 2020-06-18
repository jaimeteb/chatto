package main

import (
	"github.com/jaimeteb/chatto/core"
	"github.com/jaimeteb/chatto/server"
)

func main() {
	bot := core.LoadYAML()
	core.LoadConv()

	server.ServeBot(&bot)
}
