package rest

import (
	"encoding/json"
	"net/http"

	"github.com/jaimeteb/chatto/message"
)

// Channel contains a REST client
type Channel struct {
}

// SendMessage for REST
func (c *Channel) SendMessage(msg message.Message, recipient string) error {
	return nil
}

// ReceiveMessage for REST
func (c *Channel) ReceiveMessage(w http.ResponseWriter, r *http.Request) (message.Message, error) {
	decoder := json.NewDecoder(r.Body)
	var mess message.Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return message.Message{}, err
	}

	return mess, nil
}

// ReceiveMessages uses event queues to receive messages. Starts a long running process
func (c *Channel) ReceiveMessages(messageChan chan message.Message) {
	// Not implemented
}
