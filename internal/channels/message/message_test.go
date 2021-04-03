package message_test

import (
	"testing"

	"github.com/jaimeteb/chatto/internal/channels/message"
	"github.com/jaimeteb/chatto/query"
)

func TestReceive_Conversation(t *testing.T) {
	type fields struct {
		Question  *query.Question
		ReplyOpts *message.ReplyOpts
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "should set the conversation value to the sender when ReplyOpts is nil",
			fields: fields{
				Question: &query.Question{
					Sender: "42",
					Text:   "Testing 123...",
				},
				ReplyOpts: nil,
			},
			want: "42",
		},
		{
			name: "should set the conversation value to the sender when ReplyOpts is not nil but empty",
			fields: fields{
				Question: &query.Question{
					Sender: "42",
					Text:   "Testing 123...",
				},
				ReplyOpts: &message.ReplyOpts{},
			},
			want: "42",
		},
		{
			name: "should set the conversation value to the sender when using twilio",
			fields: fields{
				Question: &query.Question{
					Sender: "42",
					Text:   "Testing 123...",
				},
				ReplyOpts: &message.ReplyOpts{
					Twilio: message.TwilioReplyOpts{
						Recipient: "42",
					},
				},
			},
			want: "42",
		},
		{
			name: "should set the conversation value to the sender when using telegram",
			fields: fields{
				Question: &query.Question{
					Sender: "42",
					Text:   "Testing 123...",
				},
				ReplyOpts: &message.ReplyOpts{
					Telegram: message.TelegramReplyOpts{
						Recipient: "42",
					},
				},
			},
			want: "42",
		},
		{
			name: "should set the conversation value to the channel and ts when using slack",
			fields: fields{
				Question: &query.Question{
					Sender: "42",
					Text:   "Testing 123...",
				},
				ReplyOpts: &message.ReplyOpts{
					Slack: message.SlackReplyOpts{
						Channel: "C01L96YPUH4",
						TS:      "1612126789.000200",
					},
				},
			},
			want: "C01L96YPUH4/1612126789.000200",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &message.Request{
				Question:  tt.fields.Question,
				ReplyOpts: tt.fields.ReplyOpts,
			}
			if got := r.Conversation(); got != tt.want {
				t.Errorf("Conversation() = %v, want %v", got, tt.want)
			}
		})
	}
}
