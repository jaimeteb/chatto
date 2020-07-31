package main

import (
	"flag"

	"github.com/jaimeteb/chatto/bot"
)

func main() {
	cli := flag.Bool("cli", false, "Run in CLI mode")
	flag.Parse()
	if *cli {
		bot.CLI()
	} else {
		bot.ServeBot()
	}
}
