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

// SendAndReceive send a message to localhost endpoint and receives an answer
func SendAndReceive(mess *Message, url string) *[]Message {
	jsonMess, _ := json.Marshal(mess)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMess))
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
func CLI(port *int) {
	time.Sleep(time.Second * 10)

	localEndpoint := fmt.Sprintf("http://localhost:%v/endpoints/rest", *port)

	reader := bufio.NewReader(os.Stdin)
	for {
		log.Info("you:\t| ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err)
		}

		respMess := SendAndReceive(&Message{"cli", strings.TrimSuffix(cmd, "\n")}, localEndpoint)
		for _, msg := range *respMess {
			log.Infof("botto:\t| %v\n", msg.Text)
		}
	}
}
