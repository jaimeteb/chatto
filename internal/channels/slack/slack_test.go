package slack_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/mock/gomock"
	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/internal/channels/slack"
	"github.com/jaimeteb/chatto/internal/channels/slack/mockslack"
	"github.com/jaimeteb/chatto/query"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func TestChannel_SendMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	slackClient := mockslack.NewMockClient(ctrl)

	type fields struct {
		Client          slack.Client
		mockPostMessage *gomock.Call
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
			name: "send message to slack",
			fields: fields{
				Client:          slackClient,
				mockPostMessage: slackClient.EXPECT().PostMessage("test_channel", gomock.Any()).Return("", "", nil),
			},
			args: args{response: &messages.Response{
				Answers: []query.Answer{{
					Text: "Hey bud *beep* *boop*.",
				}},
				ReplyOpts: &messages.ReplyOpts{
					Slack: messages.SlackReplyOpts{
						Channel: "test_channel",
						TS:      "2021010202045",
					},
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &slack.Channel{
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
			name: "receive message from slack",
			args: args{
				body: []byte(`{"type": "message", "event": {"thread_ts": "2021010202045", "text": "hey", "user": "jaimeteb", "channel": "test_channel"}}`),
			},
			want: &messages.Receive{
				Question: &query.Question{
					Sender: "jaimeteb",
					Text:   "hey",
				},
				ReplyOpts: &messages.ReplyOpts{
					Slack: messages.SlackReplyOpts{
						Channel: "test_channel",
						TS:      "2021010202045",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &slack.Channel{}
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

func TestChannel_ReceiveMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	socketClient := mockslack.NewMockSocketClient(ctrl)

	type fields struct {
		SocketClient       slack.SocketClient
		SocketClientEvents chan socketmode.Event
		mockAck            *gomock.Call
		mockRun            *gomock.Call
	}
	type args struct {
		receiveChan chan messages.Receive
		slackEvent  socketmode.Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   messages.Receive
	}{
		{
			name: "test slack socketmode",
			fields: fields{
				SocketClient:       socketClient,
				SocketClientEvents: make(chan socketmode.Event),
				mockAck:            socketClient.EXPECT().Ack(gomock.Any()).Return(),
				mockRun: socketClient.EXPECT().Run().Do(func() {
					time.Sleep(5 * time.Second)
				}),
			},
			args: args{
				receiveChan: make(chan messages.Receive),
				slackEvent: socketmode.Event{
					Type: socketmode.EventTypeEventsAPI,
					Data: slackevents.EventsAPIEvent{
						Type: slackevents.CallbackEvent,
						InnerEvent: slackevents.EventsAPIInnerEvent{
							Data: &slackevents.MessageEvent{
								ThreadTimeStamp: "2021010202045",
								Text:            "Hey.",
								User:            "jaimeteb",
								Channel:         "test_channel",
							},
						},
					},
					Request: &socketmode.Request{},
				},
			},
			want: messages.Receive{
				Question: &query.Question{
					Sender: "jaimeteb",
					Text:   "Hey.",
				},
				ReplyOpts: &messages.ReplyOpts{
					Slack: messages.SlackReplyOpts{
						Channel: "test_channel",
						TS:      "2021010202045",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &slack.Channel{
				SocketClient:       tt.fields.SocketClient,
				SocketClientEvents: tt.fields.SocketClientEvents,
			}

			go c.ReceiveMessages(tt.args.receiveChan)

			tt.fields.SocketClientEvents <- tt.args.slackEvent

			for got := range tt.args.receiveChan {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Channel.ReceiveMessages() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
				}

				break
			}
		})
	}
}
