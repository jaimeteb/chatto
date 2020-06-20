package main

import (
	"github.com/jaimeteb/chatto/core"
)

func main() {
	// bot := core.LoadYAML()
	convs := core.LoadConv()

	chain := core.NewChain(3)
	chain.Build(convs)

	// server.ServeBot(&bot)
}
