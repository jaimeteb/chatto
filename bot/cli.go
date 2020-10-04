package bot

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const localEndpoint = "http://localhost:4770/endpoints/rest"

// SendAndReceive send a message to localhost endpoint and receives an answer
func SendAndReceive(mess *Message) *Message {
	jsonMess, _ := json.Marshal(mess)

	req, err := http.NewRequest("POST", localEndpoint, bytes.NewBuffer(jsonMess))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return &Message{}
	}
	defer resp.Body.Close()

	ans := &Message{}
	if err := json.NewDecoder(resp.Body).Decode(ans); err != nil {
		log.Println(err.Error())
		return &Message{}
	}
	return ans
}

// CLI runs a bot in a command line interface
func CLI() {
	time.Sleep(time.Second * 10)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("you:\t| ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print("botto:\t| ")

		// resp := bot.Answer(Message{"cli", strings.TrimSuffix(cmd, "\n")})
		respMess := SendAndReceive(&Message{"cli", strings.TrimSuffix(cmd, "\n")})
		fmt.Println(respMess.Text)
	}
}
