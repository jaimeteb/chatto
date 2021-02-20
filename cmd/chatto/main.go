package main

import (
	"flag"

	"github.com/jaimeteb/chatto/internal/bot"
	"github.com/jaimeteb/chatto/internal/logger"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger.SetLogger()

	cli := flag.Bool("cli", false, "Run in CLI mode.")
	port := flag.Int("port", 4770, "Specify port to use.")
	path := flag.String("path", ".", "Path to YAML files.")
	flag.Parse()

	if *cli {
		go bot.CLI(port)
	}

	chattoConfig, err := bot.LoadConfig(*path, *port)
	if err != nil {
		log.Fatal(err)
	}

	chatto, err := bot.New(chattoConfig)
	if err != nil {
		log.Fatal(err)
	}

	chatto.Run()
}
