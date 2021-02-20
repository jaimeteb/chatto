package fsm_test

import (
	"reflect"
	"testing"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

var domainOnOff = fsm.NewDomain(
	[]string{"turn_on", "turn_off"},
	[]string{"on", "off"},
	[]fsm.Function{
		{
			Transition: fsm.Transition{
				From: "off",
				Into: "on",
			},
			Command: "turn_on",
			Message: []fsm.Message{{
				Text: "Turning on.",
			}},
		},
		{
			Transition: fsm.Transition{
				From: "on",
				Into: "off",
			},
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
	fsm.Defaults{
		Unknown: "Can't do that.",
		Unsure:  "???",
		Error:   "Error",
	},
)

func TestFSM_ExecuteCmd(t *testing.T) {
	type fields struct {
		State int
		Slots map[string]string
	}
	type args struct {
		cmd       string
		txt       string
		fsmDomain *fsm.Domain
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantAnswers   []query.Answer
		wantExtension string
	}{
		{
			fields: fields{
				State: 0,
				Slots: make(map[string]string),
			},
			args: args{
				cmd:       "turn_on",
				txt:       "blarg",
				fsmDomain: domainOnOff,
			},
			wantAnswers: []query.Answer{{
				Text: "Can't do that.",
			}},
			wantExtension: "",
		},
		{
			fields: fields{
				State: 0,
				Slots: make(map[string]string),
			},
			args: args{
				cmd:       "turn_on",
				txt:       "on",
				fsmDomain: domainOnOff,
			},
			wantAnswers: []query.Answer{{
				Text: "Turning on.",
			}},
			wantExtension: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &fsm.FSM{
				State: tt.fields.State,
				Slots: tt.fields.Slots,
			}
			gotAnswers, gotExtension := m.ExecuteCmd(tt.args.cmd, tt.args.txt, tt.args.fsmDomain)
			if !reflect.DeepEqual(gotAnswers, tt.wantAnswers) {
				t.Errorf("FSM.ExecuteCmd() gotAnswers = %v, want %v", gotAnswers, tt.wantAnswers)
			}
			if gotExtension != tt.wantExtension {
				t.Errorf("FSM.ExecuteCmd() gotExtension = %v, want %v", gotExtension, tt.wantExtension)
			}
		})
	}
}
