package rest_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
				Channel: "rest",
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

func TestChannel_ValidateCallback(t *testing.T) {
	setBearerToken := func(r *http.Request, t string) *http.Request {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
		return r
	}

	type args struct {
		c *rest.Channel
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "unauthorized call to secured rest channel",
			args: args{
				c: rest.New(rest.Config{CallbackToken: "my-test-token"}),
				r: httptest.NewRequest("POST", "/channels/rest", nil),
			},
			want: false,
		},
		{
			name: "call to secured rest channel with wrong token",
			args: args{
				c: rest.New(rest.Config{CallbackToken: "my-test-token"}),
				r: setBearerToken(httptest.NewRequest("POST", "/channels/rest", nil), "my-wrong-token"),
			},
			want: false,
		},
		{
			name: "authorized call to secured rest channel",
			args: args{
				c: rest.New(rest.Config{CallbackToken: "my-test-token"}),
				r: setBearerToken(httptest.NewRequest("POST", "/channels/rest", nil), "my-test-token"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.c.ValidateCallback(tt.args.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Channel.ReceiveMessage() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}
