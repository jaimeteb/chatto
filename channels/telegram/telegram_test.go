package telegram_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/mock/gomock"
	"github.com/jaimeteb/chatto/channels/messages"
	"github.com/jaimeteb/chatto/channels/telegram"
	"github.com/jaimeteb/chatto/channels/telegram/mocktelegram"
	"github.com/jaimeteb/chatto/query"
)

func TestChannel_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telegramClient := mocktelegram.NewMockClient(ctrl)

	respValues := url.Values{}
	respValues.Add("chat_id", "123456789")
	respValues.Add("parse_mode", "Markdown")
	respValues.Add("text", "Hey bud *beep* *boop*.")

	type fields struct {
		Client   telegram.Client
		mockCall *gomock.Call
	}
	type args struct {
		response *messages.Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "send message to telegram",
			fields: fields{
				Client:   telegramClient,
				mockCall: telegramClient.EXPECT().Call("SendMessage", respValues, gomock.Any()),
			},
			args: args{response: &messages.Response{
				Answers: []query.Answer{{
					Text: "Hey bud *beep* *boop*.",
				}},
				ReplyOpts: &messages.ReplyOpts{
					Telegram: messages.TelegramReplyOpts{
						Recipient: "123456789",
					},
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &telegram.Channel{
				Client: tt.fields.Client,
			}
			if err := c.SendMessage(tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("Channel.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannel_ReceiveMessage(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *messages.Receive
		wantErr bool
	}{
		{
			name: "receive message from telegram",
			args: args{
				body: []byte(`{"update_id": 123, "message": {"message_id": 456, "text": "Hey.", "from": {"id": 789, "first_name": "jaime", "username": "jaimeteb"}}}`),
			},
			want: &messages.Receive{
				Question: &query.Question{
					Sender: "789",
					Text:   "Hey.",
				},
				ReplyOpts: &messages.ReplyOpts{
					Telegram: messages.TelegramReplyOpts{
						Recipient: "789",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &telegram.Channel{}
			got, err := c.ReceiveMessage(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Channel.ReceiveMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Channel.ReceiveMessage() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}
