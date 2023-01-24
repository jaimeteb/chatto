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

// Client allows you to chat with chatto through the REST channel
// Great for testing or integrating with your own tools
type Client struct {
	URL   string
	token string
	http  *retryablehttp.Client
}

// NewClient instantiates a new chatto client
func NewClient(url string, port int, token string) *Client {
	client := &Client{
		URL:   fmt.Sprintf("%s:%d/channels/rest", url, port),
		token: token,
		http:  retryablehttp.NewClient(),
	}
	client.http.Logger = nil

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
		if err != nil {
			log.Error(err)
			continue
		}

		color.Cyan("bot:")

		for _, answer := range answers {
			fmt.Println(answer.Text)
		}
	}
}

// Submit a question to chatto bot and get an answer
func (c *Client) Submit(question *query.Question) ([]query.Answer, error) {
	questionJSON, err := json.Marshal(question)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest("POST", c.URL, bytes.NewBuffer(questionJSON))
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	answers := []query.Answer{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&answers); decodeErr != nil {
		return nil, decodeErr
	}

	return answers, nil
}
