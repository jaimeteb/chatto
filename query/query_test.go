package query_test

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/query"
)

func TestQuery_NewMessageFromMap(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    query.Answer
		wantErr bool
	}{
		{
			name: "answer from map[interface{}]interface{}",
			args: args{
				i: map[interface{}]interface{}{"text": "text", "image": "image"},
			},
			want: query.Answer{Text: "text", Image: "image"},
		},
		{
			name: "answer from map[string]interface{}",
			args: args{
				i: map[string]interface{}{"text": "text", "image": "image"},
			},
			want: query.Answer{Text: "text", Image: "image"},
		},
		{
			name: "answer from map[string]string",
			args: args{
				i: map[string]string{"text": "text", "image": "image"},
			},
			want: query.Answer{Text: "text", Image: "image"},
		},
		{
			name: "answer from invalid interface",
			args: args{
				i: map[string]int{"text": 0, "image": 0},
			},
			want: query.Answer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := query.NewMessageFromMap(tt.args.i)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.NewMessageFromMap() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}

func TestQuery_Answers(t *testing.T) {
	tests := []struct {
		name    string
		args    interface{}
		want    []query.Answer
		wantErr bool
	}{
		{
			name: "answers from SubmitMessageRequest",
			args: query.Answer{Text: "text", Image: "image"},
			want: []query.Answer{{Text: "text", Image: "image"}},
		},
		{
			name: "answers from string",
			args: "text",
			want: []query.Answer{{Text: "text"}},
		},
		{
			name: "answers from map[string]string",
			args: map[string]string{"text": "text", "image": "image"},
			want: []query.Answer{{Text: "text", Image: "image"}},
		},
		{
			name: "answers from invalid interface",
			args: 0,
			want: []query.Answer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := query.Answers(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.Answers() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}
