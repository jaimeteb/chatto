package main

import (
	"flag"

	"github.com/jaimeteb/chatto/bot"
)

func main() {
	cli := flag.Bool("cli", false, "Run in CLI mode.")
	path := flag.String("path", ".", "Path to YAML files.")
	flag.Parse()
	if *cli {
		bot.CLI(path)
	} else {
		bot.ServeBot(path)
	}
}
