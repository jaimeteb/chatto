package main

import (
	"flag"

	"github.com/jaimeteb/chatto/bot"
	"github.com/jaimeteb/chatto/internal/logger"
)

func main() {
	port := flag.Int("port", 4770, "Specify port to use.")
	path := flag.String("path", ".", "Path to YAML files.")
	debug := flag.Bool("debug", false, "Enable debug logging.")
	flag.Parse()

	logger.SetLogger(*debug)

	server := bot.NewServer(*path, *port)

	server.Run()
}
