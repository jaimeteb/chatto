package cmd

import (
	"fmt"

	"github.com/jaimeteb/chatto/bot"
	"github.com/jaimeteb/chatto/internal/logger"
	"github.com/jaimeteb/chatto/version"
	"github.com/spf13/cobra"
)

var (
	chattoCliURL   string
	chattoCliToken string
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Chatto CLI.",
	Long:  `With the Chatto CLI you can talk to your bot directly from your terminal.`,
	Run:   cli,
}

func init() {
	rootCmd.AddCommand(cliCmd)

	cliCmd.Flags().StringVar(&chattoCliURL, "url", "http://localhost", "Specify REST channel url to connect")
	cliCmd.Flags().IntVar(&chattoPort, "port", 4770, "Specify REST channel port to connect")
	cliCmd.Flags().StringVar(&chattoCliToken, "token", "", "Specify REST channel auth token to use")
}

func cli(cmd *cobra.Command, args []string) {
	if chattoVersion {
		fmt.Println(version.BuildStr())
		return
	}
	logger.SetLogger(debug)

	client := bot.NewClient(chattoCliURL, chattoPort, chattoCliToken)
	client.CLI()
}
