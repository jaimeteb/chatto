package bot

import (
	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/channels"
	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	Name       string
	Store      fsm.Store
	Domain     *fsm.Domain
	Classifier *clf.Classifier
	Extension  extension.Extension
	Channels   *channels.Channels
	Config     *Config
	Router     *mux.Router
}

// Prediction models a classifier prediction and its original string
type Prediction struct {
	Original    string  `json:"original"`
	Predicted   string  `json:"predicted"`
	Probability float64 `json:"probability"`
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b *Bot) Answer(question *query.Question) ([]query.Answer, error) {
	if !b.Store.Exists(question.Sender) {
		b.Store.Set(
			question.Sender,
			&fsm.FSM{
				State: 0,
				Slots: make(map[string]string),
			},
		)
	}

	cmd, _ := b.Classifier.Predict(question.Text)

	machine := b.Store.Get(question.Sender)

	reply, ext := machine.ExecuteCmd(cmd, question.Text, b.Domain)

	var err error
	if ext != "" && b.Extension != nil {
		reply, err = b.Extension.RunFunc(question, ext, b.Domain, machine)
		if err != nil {
			return nil, err
		}
	}

	b.Store.Set(question.Sender, machine)

	return reply, nil
}
