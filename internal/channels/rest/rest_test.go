package rest_test

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/internal/channels/rest"
	"github.com/jaimeteb/chatto/query"
)

func TestChannel_ReceiveMessage(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		c       *rest.Channel
		args    args
		want    *messages.Receive
		wantErr bool
	}{
		{
			name: "receive message from rest",
			args: args{
				body: []byte(`{"sender": "jaimeteb", "text": "Hey."}`),
			},
			want: &messages.Receive{
				Question: &query.Question{
					Sender: "jaimeteb",
					Text:   "Hey.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &rest.Channel{}
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
