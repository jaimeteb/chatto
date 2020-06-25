package main

import (
	"github.com/jaimeteb/chatto/pkg"
	"github.com/jaimeteb/chatto/server"
)

func main() {
	bot := pkg.LoadYAML()

	convs := pkg.LoadConv()

	chain := pkg.NewChain(bot.StateSize)
	chain.Build(convs)

	// gen := chain.Generate(5)
	// fmt.Println(gen)

	server.ServeBot(&bot)
}
