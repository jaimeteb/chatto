package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jaimeteb/chatto/bot"
	"github.com/jaimeteb/chatto/internal/logger"
	"github.com/jaimeteb/chatto/version"
	"github.com/spf13/cobra"
)

var (
	debug         bool
	chattoVersion bool
	chattoPath    string
	chattoPort    int
)

var rootCmd = &cobra.Command{
	Use:   "chatto",
	Short: "Run your Chatto bot.",
	Long: `Simple chatbot framework written in Go, with configurations in YAML.
Chatto helps you create very simple text-based chatbots using a few configuration files.`,
	Run: chatto,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().BoolVarP(&chattoVersion, "version", "v", false, "Display version")

	rootCmd.Flags().IntVar(&chattoPort, "port", 4770, "Specify port to use")
	rootCmd.Flags().StringVarP(&chattoPath, "path", "p", ".", "Path to YAML files")
}

func chatto(cmd *cobra.Command, args []string) {
	if chattoVersion {
		fmt.Println(version.BuildStr())
		return
	}
	if strings.EqualFold(os.Getenv("CHATTO_BOT_DEBUG"), "true") {
		debug = true
	}
	logger.SetLogger(debug)

	server := bot.NewServer(chattoPath, chattoPort)
	server.Run()
}
