package bot

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/channels"
	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/extension"
	store "github.com/jaimeteb/chatto/internal/fsm/store"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	Name       string
	Store      store.Store
	Domain     *fsm.Domain
	Classifier *clf.Classifier
	Extensions extension.ServerMap
	Channels   *channels.Channels
	Config     *Config
	Router     *mux.Router
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b *Bot) Answer(receiveMsg *messages.Receive) ([]query.Answer, error) {
	isExistingConversation := b.Store.Exists(receiveMsg.Conversation())

	if !isExistingConversation {
		b.Store.Set(receiveMsg.Conversation(), fsm.NewFSM())
	}

	cmd, _ := b.Classifier.Model.Predict(receiveMsg.Question.Text, b.Classifier.Pipeline)

	machine := b.Store.Get(receiveMsg.Conversation())

	previousState := machine.State

	// Set existing conversation to false if in the initial state
	// because initial state means this is a new conversation
	if machine.State == fsm.StateInitial {
		isExistingConversation = false
	}

	answers, ext, err := machine.ExecuteCmd(cmd, receiveMsg.Question.Text, b.Domain)
	if err != nil {
		switch e := err.(type) {
		case *fsm.ErrUnsureCommand:
			if b.Config.ShouldReplyUnsure(isExistingConversation) {
				return []query.Answer{{Text: e.Error()}}, nil
			}

			return []query.Answer{}, nil
		case *fsm.ErrUnknownCommand:
			if b.Config.ShouldReplyUnknown(isExistingConversation) {
				return []query.Answer{{Text: e.Error()}}, nil
			}

			return []query.Answer{}, nil
		default:
			return nil, err
		}
	}

	log.Debugf("FSM | State transitioned from '%d' -> '%d'", previousState, machine.State)

	if ext != nil {
		if _, ok := b.Extensions[ext.Server]; !ok {
			return nil, &ErrUnknownExtension{Extension: ext.Server}
		}

		answers, err = b.Extensions[ext.Server].ExecuteExtension(receiveMsg.Question, ext.Name, receiveMsg.Channel, cmd, b.Domain, machine)
		if err != nil {
			return []query.Answer{{Text: b.Domain.DefaultMessages.Error}}, nil
		}
	}

	b.Store.Set(receiveMsg.Conversation(), machine)

	return answers, nil
}

// ErrUnknownExtension is returned by the Bot when
// the provided extension name does not exist
type ErrUnknownExtension struct {
	Extension string
}

// Error returns the ErrUnknownExtension error message
func (e *ErrUnknownExtension) Error() string {
	return fmt.Sprintf("cannot answer: extension %s is unknown", e.Extension)
}
