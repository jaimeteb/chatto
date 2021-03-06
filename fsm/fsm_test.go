package fsm_test

import (
	"reflect"
	"testing"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

var (
	// Hello
	helloFunctions = []fsm.Transition{
		{
			From:    []string{"any"},
			Into:    "hello",
			Command: "hey_friend",
			Answers: []fsm.Answer{{
				Text: "Hey friend!",
			}},
		},
	}

	// Pokemon test
	pokemonFunctions = []fsm.Transition{
		{
			From:    []string{"initial"},
			Into:    "search_pokemon",
			Command: "search_pokemon",
			Answers: []fsm.Answer{{
				Text: "What is the Pokémon's name or number?",
			}},
		},
		{
			From:    []string{"initial"},
			Into:    "search_pokemon",
			Command: "greet",
			Answers: []fsm.Answer{{
				Text: "What is the Pokémon's name or number?",
			}},
		},
		{
			From:    []string{"search_pokemon"},
			Into:    "initial",
			Command: "any",
			Extension: fsm.Extension{
				Server: "pokemon",
				Name:   "search_pokemon",
			},
			Slot: fsm.Slot{
				Name: "pokemon",
				Mode: "whole_text",
			},
		},
	}

	// On off test
	onOffFunctions = []fsm.Transition{
		{
			From:    []string{"initial"},
			Into:    "on",
			Command: "turn_on",
			Answers: []fsm.Answer{{
				Text: "Turning on.",
			}},
			Slot: fsm.Slot{
				Name: "on",
				Mode: "whole_text",
			},
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
			Slot: fsm.Slot{
				Name:  "off",
				Mode:  "regex",
				Regex: "^turn.*$",
			},
		},
	}
	defaultResponses = fsm.Defaults{
		Unknown: "Can't do that.",
		Unsure:  "???",
		Error:   "Error",
	}
)

func TestFSM_ExecuteCmd(t *testing.T) {
	type fields struct {
		State int
		Slots map[string]string
	}
	type args struct {
		command   string
		txt       string
		fsmDomain *fsm.Domain
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantAnswers   []query.Answer
		wantExtension *fsm.Extension
		wantState     int
		wantSlots     map[string]string
		wantErr       bool
	}{
		{
			name: "invalid command should be unknown",
			fields: fields{
				State: fsm.StateInitial,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "ruhrow",
				txt:       "blah blah blah",
				fsmDomain: fsm.NewDomain(onOffFunctions, defaultResponses),
			},
			wantAnswers:   nil,
			wantExtension: nil,
			wantState:     fsm.StateInitial,
			wantSlots:     map[string]string{},
			wantErr:       true,
		},
		{
			name: "empty command should be unsure",
			fields: fields{
				State: fsm.StateInitial,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "",
				txt:       "blah blah blah",
				fsmDomain: fsm.NewDomain(onOffFunctions, defaultResponses),
			},
			wantAnswers:   nil,
			wantExtension: nil,
			wantState:     fsm.StateInitial,
			wantSlots:     map[string]string{},
			wantErr:       true,
		},
		{
			name: "hello command should run hey_friend from any state",
			fields: fields{
				State: fsm.StateInitial,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "hey_friend",
				txt:       "hey there",
				fsmDomain: fsm.NewDomain(helloFunctions, defaultResponses),
			},
			wantAnswers: []query.Answer{{
				Text: "Hey friend!",
			}},
			wantExtension: nil,
			wantState:     1,
			wantSlots:     map[string]string{},
			wantErr:       false,
		},
		{
			name: "any command should run extension search_pokemon",
			fields: fields{
				State: 1,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "any",
				txt:       "pikachu",
				fsmDomain: fsm.NewDomain(pokemonFunctions, defaultResponses),
			},
			wantAnswers: nil,
			wantExtension: &fsm.Extension{
				Server: "pokemon",
				Name:   "search_pokemon",
			},
			wantState: fsm.StateInitial,
			wantSlots: map[string]string{"pokemon": "pikachu"},
			wantErr:   false,
		},
		{
			name: "turn_on command should turn on",
			fields: fields{
				State: fsm.StateInitial,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "turn_on",
				txt:       "turn it on",
				fsmDomain: fsm.NewDomain(onOffFunctions, defaultResponses),
			},
			wantAnswers: []query.Answer{{
				Text: "Turning on.",
			}},
			wantExtension: nil,
			wantState:     1,
			wantSlots:     map[string]string{"on": "turn it on"},
			wantErr:       false,
		},
		{
			name: "turn_off command should turn off",
			fields: fields{
				State: 1,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "turn_off",
				txt:       "turn it off",
				fsmDomain: fsm.NewDomain(onOffFunctions, defaultResponses),
			},
			wantAnswers: []query.Answer{
				{
					Text: "Turning off.",
				},
				{
					Text: "❌",
				},
			},
			wantExtension: nil,
			wantState:     fsm.StateInitial,
			wantSlots:     map[string]string{"off": "turn it off"},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &fsm.FSM{
				State: tt.fields.State,
				Slots: tt.fields.Slots,
			}
			gotAnswers, gotExtension, err := m.ExecuteCmd(tt.args.command, tt.args.txt, tt.args.fsmDomain)
			if (err != nil) != tt.wantErr {
				t.Errorf("FSM.ExecuteCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAnswers, tt.wantAnswers) {
				t.Errorf("FSM.ExecuteCmd() gotAnswers = %v, want %v", gotAnswers, tt.wantAnswers)
			}
			if !reflect.DeepEqual(gotExtension, tt.wantExtension) {
				t.Errorf("FSM.ExecuteCmd() gotExtension = %v, want %v", gotExtension, tt.wantExtension)
			}
			if !reflect.DeepEqual(m.Slots, tt.wantSlots) {
				t.Errorf("FSM.Slot gotSlot = %v, want %v", m.Slots, tt.wantSlots)
			}
		})
	}
}
