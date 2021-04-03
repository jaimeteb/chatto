package bot

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/channels"
	"github.com/jaimeteb/chatto/internal/channels/message"
	"github.com/jaimeteb/chatto/internal/clf"
	"github.com/jaimeteb/chatto/internal/extensions"
	store "github.com/jaimeteb/chatto/internal/fsm/store"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	Name                 string
	Store                store.Store
	Domain               *fsm.Domain
	Classifier           *clf.Classifier
	Extensions           extensions.ServerMap
	WebSocket            *extensions.WebSocketServer
	Channels             *channels.Channels
	Config               *Config
	Router               *mux.Router
	MessageResponseQueue chan message.Response
}

// SubmitMessageRequest takes a users input submits it to the classifier to predict the
// command that should be executed. Then the predicted command gets passed through the
// fsm.FSM which transitions the conversation state and returns the answer or extension
// to be executed
func (b *Bot) SubmitMessageRequest(msgRequest *message.Request) error {
	msgResponse := message.Response{
		ReplyOpts: msgRequest.ReplyOpts,
		Channel:   msgRequest.Channel,
	}

	isExistingConversation := b.Store.Exists(msgRequest.Conversation())

	if !isExistingConversation {
		b.Store.Set(msgRequest.Conversation(), fsm.NewFSM())
	}

	cmd, _ := b.Classifier.Model.Predict(msgRequest.Question.Text, b.Classifier.Pipeline)

	machine := b.Store.Get(msgRequest.Conversation())

	previousState := machine.State

	// Set existing conversation to false if in the initial state
	// because initial state means this is a new conversation
	if machine.State == fsm.StateInitial {
		isExistingConversation = false
	}

	answers, extension, err := machine.ExecuteCmd(cmd, msgRequest.Question.Text, b.Domain)
	if err != nil {
		switch e := err.(type) {
		case *fsm.ErrUnsureCommand:
			if b.Config.ShouldReplyUnsure(isExistingConversation) {
				msgResponse.Answers = []query.Answer{{Text: e.Error()}}
				b.MessageResponseQueue <- msgResponse
				return nil
			}

			return nil
		case *fsm.ErrUnknownCommand:
			if b.Config.ShouldReplyUnknown(isExistingConversation) {
				msgResponse.Answers = []query.Answer{{Text: e.Error()}}
				b.MessageResponseQueue <- msgResponse
				return nil
			}

			return nil
		default:
			return err
		}
	}

	b.Store.Set(msgRequest.Conversation(), machine)

	log.Debugf("FSM | State transitioned from '%d' -> '%d'", previousState, machine.State)

	switch {
	case len(answers) > 0:
		msgResponse.Answers = answers
		b.MessageResponseQueue <- msgResponse
	case extension != nil:
		if _, ok := b.Extensions[extension.Server]; !ok {
			return &ErrUnknownExtension{Extension: extension.Server}
		}

		err = b.Extensions[extension.Server].Execute(extension.Name, *msgRequest, b.Domain, machine)
		if err != nil {
			return err
		}
	default:
		log.Info("no action for message request: %+v", msgRequest)
	}

	return nil
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
