package pipeline_test

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/internal/clf/pipeline"
)

func TestPipeline(t *testing.T) {
	type args struct {
		text string
		pipe *pipeline.Config
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "no remove_symbols no lower",
			args: args{
				text: "I don't know...",
				pipe: &pipeline.Config{false, false, 0},
			},
			want: []string{"I", "don't", "know..."},
		},
		{
			name: "remove_symbols lower",
			args: args{
				text: "I don't know...",
				pipe: &pipeline.Config{true, true, 0},
			},
			want: []string{"i", "don", "t", "know"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pipeline.Pipeline(tt.args.text, tt.args.pipe)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pipeline.Pipeline() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}
