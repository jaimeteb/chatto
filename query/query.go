package query

import (
	log "github.com/sirupsen/logrus"
)

// Question for the FSM
type Question struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
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
		msg.Text = m["text"].(string)
		msg.Image = m["image"].(string)
	case map[string]interface{}:
		msg.Text = m["text"].(string)
		msg.Image = m["image"].(string)
	case map[string]string:
		msg.Text = m["text"]
		msg.Image = m["image"]
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
