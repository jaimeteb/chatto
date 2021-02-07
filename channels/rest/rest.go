package rest

import (
	"encoding/json"
	"net/http"

	"github.com/jaimeteb/chatto/channels/messages"
)

// Channel contains a REST client
type Channel struct {
}

// SendMessage for REST
func (c *Channel) SendMessage(response *messages.Response) error {
	return nil
}

// ReceiveMessage for REST
func (c *Channel) ReceiveMessage(w http.ResponseWriter, r *http.Request) (*messages.Receive, error) {
	decoder := json.NewDecoder(r.Body)

	var receive messages.Receive

	err := decoder.Decode(&receive)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return &receive, nil
}

// ReceiveMessages uses event queues to receive messages. Starts a long running process
func (c *Channel) ReceiveMessages(receiveChan chan messages.Receive) {
	// Not implemented
}
