package fsm_test

import (
	"reflect"
	"testing"

	"github.com/jaimeteb/chatto/fsm"
	fsmint "github.com/jaimeteb/chatto/internal/fsm"
	"github.com/jaimeteb/chatto/internal/testutils"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *fsmint.Config
		wantErr bool
	}{
		{
			name: "test loading a valid path",
			args: args{path: "../" + testutils.Examples05SimplePath},
			want: &fsmint.Config{
				Transitions: []fsm.Transition{
					{
						From:    []string{"initial"},
						Into:    "on",
						Command: "turn_on",
						Answers: []fsm.Answer{{
							Text: "Turning on.",
						}},
					},
					{
						From:    []string{"on"},
						Into:    "initial",
						Command: "turn_off",
						Answers: []fsm.Answer{
							{
								Text: "Turning off.",
							},
							{
								Text: "❌",
							},
						},
					},
				},
				Defaults: fsm.Defaults{
					Unknown: "Can't do that.",
					Unsure:  "???",
					Error:   "Error",
				},
			},
			wantErr: false,
		},
		{
			name:    "test loading a invalid path",
			args:    args{path: "../" + testutils.Examples00InvalidPath},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fsmint.LoadConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
