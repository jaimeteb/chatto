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

	log "github.com/sirupsen/logrus"
)

var localEndpoint = fmt.Sprintf("http://localhost:%v/endpoints/rest", chattoPort)

// SendAndReceive send a message to localhost endpoint and receives an answer
func SendAndReceive(mess *Message) *[]Message {
	jsonMess, _ := json.Marshal(mess)

	req, err := http.NewRequest("POST", localEndpoint, bytes.NewBuffer(jsonMess))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Warn(err.Error())
		return &[]Message{}
	}
	defer resp.Body.Close()

	ans := &[]Message{}
	if err := json.NewDecoder(resp.Body).Decode(ans); err != nil {
		log.Warn(err.Error())
		return &[]Message{}
	}
	return ans
}

// CLI runs a bot in a command line interface
func CLI() {
	time.Sleep(time.Second * 10)

	reader := bufio.NewReader(os.Stdin)
	for {
		log.Info("you:\t| ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err)
		}

		// resp := bot.Answer(Message{"cli", strings.TrimSuffix(cmd, "\n")})
		respMess := SendAndReceive(&Message{"cli", strings.TrimSuffix(cmd, "\n")})
		for _, msg := range *respMess {
			log.Info("botto:\t| %v\n", msg.Text)
		}
	}
}
