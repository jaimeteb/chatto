package twilio_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/mock/gomock"
	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/internal/channels/twilio"
	"github.com/jaimeteb/chatto/internal/channels/twilio/mocktwilio"
	"github.com/jaimeteb/chatto/query"
	twlio "github.com/kevinburke/twilio-go"
)

func TestChannel_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	twilioClient := mocktwilio.NewMockClient(ctrl)

	type fields struct {
		Client          twilio.Client
		Number          string
		mockSendMessage *gomock.Call
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
			name: "send message to twilio",
			fields: fields{
				Client:          twilioClient,
				Number:          "123456789",
				mockSendMessage: twilioClient.EXPECT().SendMessage("123456789", "42", "Hey bud *beep* *boop*.", nil).Return(&twlio.Message{}, nil),
			},
			args: args{response: &messages.Response{
				Answers: []query.Answer{{
					Text: "Hey bud *beep* *boop*.",
				}},
				ReplyOpts: &messages.ReplyOpts{
					Twilio: messages.TwilioReplyOpts{
						Recipient: "42",
					},
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &twilio.Channel{
				Client: tt.fields.Client,
				Number: tt.fields.Number,
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
			name: "receive message from twilio",
			args: args{
				body: []byte(url.Values{"From": {"42"}, "Body": {"Hey."}}.Encode()),
			},
			want: &messages.Receive{
				Question: &query.Question{
					Sender: "42",
					Text:   "Hey.",
				},
				ReplyOpts: &messages.ReplyOpts{
					Twilio: messages.TwilioReplyOpts{
						Recipient: "42",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &twilio.Channel{}
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
