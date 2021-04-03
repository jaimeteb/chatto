package telegram_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/mock/gomock"
	"github.com/jaimeteb/chatto/internal/channels/message"
	"github.com/jaimeteb/chatto/internal/channels/telegram"
	"github.com/jaimeteb/chatto/internal/channels/telegram/mocktelegram"
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
		response *message.Response
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
				mockCall: telegramClient.EXPECT().Call("MessageResponse", respValues, gomock.Any()),
			},
			args: args{response: &message.Response{
				Answers: []query.Answer{{
					Text: "Hey bud *beep* *boop*.",
				}},
				ReplyOpts: &message.ReplyOpts{
					Telegram: message.TelegramReplyOpts{
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
			if err := c.MessageResponse(tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("Channel.MessageResponse() error = %v, wantErr %v", err, tt.wantErr)
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
		want    *message.Request
		wantErr bool
	}{
		{
			name: "receive message from telegram",
			args: args{
				body: []byte(`{"update_id": 123, "message": {"message_id": 456, "text": "Hey.", "from": {"id": 789, "first_name": "jaime", "username": "jaimeteb"}}}`),
			},
			want: &message.Request{
				Question: &query.Question{
					Sender: "789",
					Text:   "Hey.",
				},
				ReplyOpts: &message.ReplyOpts{
					Telegram: message.TelegramReplyOpts{
						Recipient: "789",
					},
				},
				Channel: "telegram",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &telegram.Channel{}
			got, err := c.MessageRequest(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Channel.MessageRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Channel.MessageRequest() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}
