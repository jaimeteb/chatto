package messages

import (
	"strings"

	"github.com/jaimeteb/chatto/query"
)

// Receive question from channel with reply options
type Receive struct {
	Question  *query.Question `json:"question"`
	ReplyOpts *ReplyOpts
}

// Conversation returns a string of the unique conversations
// the bot is having. The definition of a "conversation" is
// different depending on the channel used. For example Slack
// conversations happen in Slack threads, which is different
// than a conversation in Twilio between a Sender and Recipient.
// In the Slack example this allows us to have multiple different
// conversations with the bot in different Slack channels and
// threads
func (r *Receive) Conversation() string {
	if r.ReplyOpts == nil {
		return r.Question.Sender
	}

	if r.ReplyOpts.Slack != (SlackReplyOpts{}) {
		return strings.Join([]string{
			r.ReplyOpts.Slack.Channel,
			r.ReplyOpts.Slack.TS,
		}, "/")
	} else if r.ReplyOpts.Telegram != (TelegramReplyOpts{}) {
		return r.Question.Sender
	} else if r.ReplyOpts.Twilio != (TwilioReplyOpts{}) {
		return r.Question.Sender
	}

	return r.Question.Sender
}

// Response with answers to channel with reply options
type Response struct {
	Answers   []query.Answer
	ReplyOpts *ReplyOpts
}

// ReplyOpts allow you to configure how the reply is sent
type ReplyOpts struct {
	Telegram TelegramReplyOpts
	Twilio   TwilioReplyOpts
	Slack    SlackReplyOpts
}

// TelegramReplyOpts are options used to reply with Telegram
type TelegramReplyOpts struct {
	Recipient string
}

// TwilioReplyOpts are options used to reply with Twilio
type TwilioReplyOpts struct {
	Recipient string
}

// SlackReplyOpts are options used to reply with Slack
type SlackReplyOpts struct {
	Channel string
	TS      string
}
