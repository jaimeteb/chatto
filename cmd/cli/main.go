package main

import (
	"flag"

	"github.com/jaimeteb/chatto/internal/bot"
	"github.com/jaimeteb/chatto/internal/logger"
)

func main() {
	url := flag.String("url", "http://localhost", "Specify url to use.")
	port := flag.Int("port", 4770, "Specify port to use.")
	debug := flag.Bool("debug", false, "Enable debug logging.")
	flag.Parse()

	logger.SetLogger(*debug)

	cli := bot.NewCLI(*url, *port)

	cli.Run()
}
