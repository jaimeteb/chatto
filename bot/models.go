package bot

// Message models and incoming/outgoing message
type Message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}
