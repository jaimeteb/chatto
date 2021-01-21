package common

// Message models and incoming/outgoing message
type Message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
	Image  string `json:"image"`
}

// MessageFromMap converts a map of interfaces or strings into a Message
func MessageFromMap(msgMap interface{}) Message {
	msg := Message{}
	switch m := msgMap.(type) {
	case map[interface{}]interface{}:
		msg.Sender, _ = m["sender"].(string)
		msg.Text, _ = m["text"].(string)
		msg.Image, _ = m["image"].(string)
	case map[string]interface{}:
		msg.Sender, _ = m["sender"].(string)
		msg.Text, _ = m["text"].(string)
		msg.Image, _ = m["image"].(string)
	case map[string]string:
		msg.Sender, _ = m["sender"]
		msg.Text, _ = m["text"]
		msg.Image, _ = m["image"]
	}
	return msg
}

// Out creates an outgoing message without empty fields
func (m *Message) Out() map[string]string {
	o := make(map[string]string)
	if m.Sender != "" {
		o["sender"] = m.Sender
	}
	if m.Text != "" {
		o["text"] = m.Text
	}
	if m.Image != "" {
		o["image"] = m.Image
	}
	return o
}
