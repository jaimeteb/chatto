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
				States:   []string{"off", "on"},
				Commands: []string{"turn_on", "turn_off"},
				Transitions: []fsm.Transition{
					{
						// Transition: fsm.Transition{
						From: []string{"off"},
						Into: "on",
						// },
						Command: "turn_on",
						Message: []fsm.Message{{
							Text: "Turning on.",
						}},
					},
					{
						// Transition: fsm.Transition{
						From: []string{"on"},
						Into: "off",
						// },
						Command: "turn_off",
						Message: []fsm.Message{
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
