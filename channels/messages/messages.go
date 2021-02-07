package messages

import "github.com/jaimeteb/chatto/query"

// Receive question from channel with reply options
type Receive struct {
	Question  *query.Question
	ReplyOpts *ReplyOpts
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
