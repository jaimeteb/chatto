package bot

import (
	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/channels"
	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/extension"
	fsmint "github.com/jaimeteb/chatto/internal/fsm"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	Name       string
	Store      fsmint.Store
	Domain     *fsm.Domain
	Classifier *clf.Classifier
	Extension  extension.Extension
	Channels   *channels.Channels
	Config     *Config
	Router     *mux.Router
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b *Bot) Answer(question *query.Question) ([]query.Answer, error) {
	if !b.Store.Exists(question.Sender) {
		b.Store.Set(question.Sender, fsm.NewFSM())
	}

	cmd, _ := b.Classifier.Predict(question.Text)

	machine := b.Store.Get(question.Sender)

	previousState := machine.State

	reply, ext := machine.ExecuteCmd(cmd, question.Text, b.Domain)

	log.Debugf("FSM | State transitioned from '%d' -> '%d'", previousState, machine.State)

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
