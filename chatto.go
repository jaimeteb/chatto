package main

import (
	"flag"
	"fmt"
	"os"
	"plugin"

	"github.com/jaimeteb/chatto/bot"
)

// Runner interface
type Runner interface {
	Run(string)
}

func main() {
	ext := "./examples/echo/ext.so"
	plug, err := plugin.Open(ext)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	echo, err := plug.Lookup("Echo")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var runner Runner
	runner, ok := echo.(Runner)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	// 4. use the module
	runner.Run("greet")
	runner.Run("goodbye")
	runner.Run("foo bar")

	////////////////////////////////////////////////////////

	cli := flag.Bool("cli", false, "Run in CLI mode.")
	path := flag.String("path", ".", "Path to YAML files.")
	flag.Parse()
	if *cli {
		bot.CLI(path)
	} else {
		bot.ServeBot(path)
	}
}
