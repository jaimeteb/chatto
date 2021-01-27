package main

import (
	"flag"

	"github.com/jaimeteb/chatto/bot"
	cmn "github.com/jaimeteb/chatto/common"
)

func init() {
	cmn.SetLogger()
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
