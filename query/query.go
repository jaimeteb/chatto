package query

import (
	log "github.com/sirupsen/logrus"
)

// Question for the FSM
type Question struct {
	Sender string
	Text   string
}

// Answer from the FSM
type Answer struct {
	Text  string `json:"text"`
	Image string `json:"image"`
}

// NewMessageFromMap converts a map of interfaces or strings into an Answer
func NewMessageFromMap(msgMap interface{}) Answer {
	msg := Answer{}
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

// Answers creates a slice of Answers from
func Answers(answers ...interface{}) []Answer {
	finalAnswers := make([]Answer, 0)

	for _, ans := range answers {
		switch a := ans.(type) {
		case Answer:
			finalAnswers = append(finalAnswers, a)
		case string:
			newAns := Answer{Text: a}
			finalAnswers = append(finalAnswers, newAns)
		case map[interface{}]interface{}, map[string]interface{}, map[string]string:
			newAns := NewMessageFromMap(a)
			finalAnswers = append(finalAnswers, newAns)
		default:
			log.Errorf("Message type unsupported: %T", a)
		}
	}

	return finalAnswers
}
