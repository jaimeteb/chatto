package main

import (
	"fmt"

	"github.com/jaimeteb/chatto/core"
)

func main() {
	bot := core.LoadYAML()
	convs := core.LoadConv()

	chain := core.NewChain(bot.PrefixLen)
	chain.Build(convs)

	gen := chain.Generate(5)
	fmt.Println(gen)

	// server.ServeBot(&bot)
}
