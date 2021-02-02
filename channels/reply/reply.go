package reply

// Message from a channel. Includes settings for responses
type Message struct {
	Reply     Reply
	Message   string
	Recipient string
	Slack     struct {
		TS string
	}
}

// Reply to a message
type Reply struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
	Image  string `json:"image"`
}
