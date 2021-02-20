package bot

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// CLI allows you to chat with botto from the cli. Great for testing
type CLI struct {
	Port int
	URL  string
	HTTP *retryablehttp.Client
}

// NewCLI instantiates a new botto command line interface
func NewCLI(url string, port int) *CLI {
	cli := &CLI{
		Port: port,
		URL:  fmt.Sprintf("%s:%d/channels/rest", url, port),
		HTTP: retryablehttp.NewClient(),
	}

	cli.HTTP.Logger = log.New()

	return cli
}

// Run starts the command line interface
func (c *CLI) Run() {
	reader := bufio.NewReader(os.Stdin)
	for {
		color.Magenta("you: ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err)
			continue
		}

		answers := c.sendAndReceive(&query.Question{
			Sender: "cli",
			Text:   strings.TrimSuffix(cmd, "\n"),
		}, c.URL)

		color.Cyan("bot:")

		for _, answer := range answers {
			fmt.Println(answer.Text)
		}
	}
}

// sendAndReceive send a message to localhost endpoint and receives an answer
func (c *CLI) sendAndReceive(question *query.Question, url string) []query.Answer {
	jsonMess, err := json.Marshal(question)
	if err != nil {
		log.Warn(err)
		return []query.Answer{}
	}

	resp, err := c.HTTP.Post(url, "Content-Type: application/json", bytes.NewBuffer(jsonMess))
	if err != nil {
		log.Warn(err)
		return nil
	}
	defer resp.Body.Close()

	answer := []query.Answer{}
	if err := json.NewDecoder(resp.Body).Decode(&answer); err != nil {
		log.Warn(err)
		return nil
	}

	return answer
}
