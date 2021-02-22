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

// Client allows you to chat with botto. Great for testing
// or integrating with your own tools
type Client struct {
	URL  string
	HTTP *retryablehttp.Client
}

// NewClient instantiates a new botto client
func NewClient(url string, port int) *Client {
	client := &Client{
		URL:  fmt.Sprintf("%s:%d/channels/rest", url, port),
		HTTP: retryablehttp.NewClient(),
	}

	client.HTTP.Logger = log.New()

	return client
}

// CLI starts a command line interface
func (c *Client) CLI() {
	reader := bufio.NewReader(os.Stdin)
	for {
		color.Magenta("you: ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err)
			continue
		}

		question := &query.Question{Sender: "cli", Text: strings.TrimSuffix(cmd, "\n")}

		answers, err := c.Submit(question)

		color.Cyan("bot:")

		for _, answer := range answers {
			fmt.Println(answer.Text)
		}
	}
}

// Submit a question to botto bot and get an answer
func (c *Client) Submit(question *query.Question) ([]query.Answer, error) {
	questionJSON, err := json.Marshal(question)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Post(c.URL, "Content-Type: application/json", bytes.NewBuffer(questionJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	answers := []query.Answer{}
	if err := json.NewDecoder(resp.Body).Decode(&answers); err != nil {
		return nil, err
	}

	return answers, nil
}
