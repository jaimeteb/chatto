package main

import (
	"flag"
	"fmt"

	"github.com/jaimeteb/chatto/bot"
	"github.com/jaimeteb/chatto/internal/logger"
	"github.com/jaimeteb/chatto/internal/version"
)

func main() {
	url := flag.String("url", "http://localhost", "Specify url to use.")
	port := flag.Int("port", 4770, "Specify port to use.")
	debug := flag.Bool("debug", false, "Enable debug logging.")
	vers := flag.Bool("version", false, "Display version.")
	flag.Parse()

	if *vers {
		fmt.Println(version.Build())
		return
	}

	logger.SetLogger(*debug)

	client := bot.NewClient(*url, *port)

	client.CLI()
}
