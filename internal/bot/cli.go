package bot

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// SendAndReceive send a message to localhost endpoint and receives an answer
func SendAndReceive(question *query.Question, url string) []query.Answer {
	jsonMess, err := json.Marshal(question)
	if err != nil {
		log.Warn(err)
		return []query.Answer{}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMess))
	if err != nil {
		log.Warn(err)
		return []query.Answer{}
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Warn(err)
		return []query.Answer{}
	}
	defer resp.Body.Close()

	ans := []query.Answer{}
	if err := json.NewDecoder(resp.Body).Decode(&ans); err != nil {
		log.Warn(err)
		return []query.Answer{}
	}

	return ans
}

// CLI runs a bot in a command line interface
func CLI(port *int) {
	time.Sleep(time.Second * 10)

	localEndpoint := fmt.Sprintf("http://localhost:%v/endpoints/rest", *port)

	reader := bufio.NewReader(os.Stdin)
	for {
		color.Magenta("you: ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err)
			continue
		}

		respMess := SendAndReceive(&query.Question{
			Sender: "cli",
			Text:   strings.TrimSuffix(cmd, "\n"),
		}, localEndpoint)

		color.Cyan("bot:")

		for _, msg := range respMess {
			fmt.Println(msg.Text)
		}
	}
}
