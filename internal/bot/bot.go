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
func (b *Bot) SubmitMessageRequest(messageRequest *message.Request) error {
	messageResponse := message.Response{
		ReplyOpts: messageRequest.ReplyOpts,
		Channel:   messageRequest.Channel,
	}

	isExistingConversation := b.Store.Exists(messageRequest.Conversation())

	if !isExistingConversation {
		b.Store.Set(messageRequest.Conversation(), fsm.NewFSM())
	}

	cmd, _ := b.Classifier.Model.Predict(messageRequest.Question.Text, b.Classifier.Pipeline)

	machine := b.Store.Get(messageRequest.Conversation())

	previousState := machine.State

	// Set existing conversation to false if in the initial state
	// because initial state means this is a new conversation
	if machine.State == fsm.StateInitial {
		isExistingConversation = false
	}

	answers, extension, err := machine.ExecuteCmd(cmd, messageRequest.Question.Text, b.Domain)
	if err != nil {
		switch e := err.(type) {
		case *fsm.ErrUnsureCommand:
			if b.Config.ShouldReplyUnsure(isExistingConversation) {
				messageResponse.Answers = []query.Answer{{Text: e.Error()}}
				b.MessageResponseQueue <- messageResponse
				return nil
			}

			return nil
		case *fsm.ErrUnknownCommand:
			if b.Config.ShouldReplyUnknown(isExistingConversation) {
				messageResponse.Answers = []query.Answer{{Text: e.Error()}}
				b.MessageResponseQueue <- messageResponse
				return nil
			}

			return nil
		default:
			return err
		}
	}

	b.Store.Set(messageRequest.Conversation(), machine)

	log.Debugf("FSM | State transitioned from '%d' -> '%d'", previousState, machine.State)

	switch {
	case len(answers) > 0:
		messageResponse.Answers = answers
		b.MessageResponseQueue <- messageResponse
	case extension != nil:
		if _, ok := b.Extensions[extension.Server]; !ok {
			return &ErrUnknownExtension{Extension: extension.Server}
		}

		err = b.Extensions[extension.Server].Execute(extension.Name, *messageRequest, b.Domain, machine)
		if err != nil {
			return err
		}
	default:
		log.Info("no action for message request: %+v", messageRequest)
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
