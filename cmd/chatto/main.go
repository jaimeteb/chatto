package main

import (
	"flag"

	"github.com/jaimeteb/chatto/bot"
	"github.com/jaimeteb/chatto/logger"
)

func init() {
	logger.SetLogger()
}

func main() {
	cli := flag.Bool("cli", false, "Run in CLI mode.")
	port := flag.Int("port", 4770, "Specify port to use.")
	path := flag.String("path", ".", "Path to YAML files.")
	flag.Parse()

	if *cli {
		go bot.CLI(port)
	}
	bot.ServeBot(path, port)
}
