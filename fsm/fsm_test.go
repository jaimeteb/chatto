package fsm_test

import (
	"reflect"
	"testing"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

var (
	// Hello
	helloCommands  = []string{"hey_friend"}
	helloStates    = []string{"hello"}
	helloFunctions = []fsm.Function{
		{
			Transition: fsm.Transition{
				From: "any",
				Into: "hello",
			},
			Command: "hey_friend",
			Message: []fsm.Message{{
				Text: "Hey friend!",
			}},
		},
	}

	// Pokemon test
	pokemonCommands  = []string{"initial", "search_pokemon"}
	pokemonStates    = []string{"greet", "search_pokemon"}
	pokemonFunctions = []fsm.Function{
		{
			Transition: fsm.Transition{
				From: "initial",
				Into: "search_pokemon",
			},
			Command: "search_pokemon",
			Message: []fsm.Message{{
				Text: "What is the Pokémon's name or number?",
			}},
		},
		{
			Transition: fsm.Transition{
				From: "initial",
				Into: "search_pokemon",
			},
			Command: "greet",
			Message: []fsm.Message{{
				Text: "What is the Pokémon's name or number?",
			}},
		},
		{
			Transition: fsm.Transition{
				From: "search_pokemon",
				Into: "initial",
			},
			Command:   "any",
			Extension: "search_pokemon",
			Slot: fsm.Slot{
				Name: "pokemon",
				Mode: "whole_text",
			},
		},
	}

	// On off test
	onOffCommands  = []string{"turn_on", "turn_off"}
	onOffStates    = []string{"off", "on"}
	onOffFunctions = []fsm.Function{
		{
			Transition: fsm.Transition{
				From: "off",
				Into: "on",
			},
			Command: "turn_on",
			Message: []fsm.Message{{
				Text: "Turning on.",
			}},
			Slot: fsm.Slot{
				Name: "on",
				Mode: "whole_text",
			},
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
		wantExtension string
		wantState     int
		wantSlots     map[string]string
		wantErr       bool
	}{
		{
			name: "invalid command should be unknown",
			fields: fields{
				State: 0,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "ruhrow",
				txt:       "blah blah blah",
				fsmDomain: fsm.NewDomain(onOffCommands, onOffStates, onOffFunctions, defaultResponses),
			},
			wantAnswers:   nil,
			wantExtension: "",
			wantState:     0,
			wantSlots:     map[string]string{},
			wantErr:       true,
		},
		{
			name: "empty command should be unsure",
			fields: fields{
				State: 0,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "",
				txt:       "blah blah blah",
				fsmDomain: fsm.NewDomain(onOffCommands, onOffStates, onOffFunctions, defaultResponses),
			},
			wantAnswers:   nil,
			wantExtension: "",
			wantState:     0,
			wantSlots:     map[string]string{},
			wantErr:       true,
		},
		{
			name: "hello command should run hey_friend from any state",
			fields: fields{
				State: 0,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "hey_friend",
				txt:       "hey there",
				fsmDomain: fsm.NewDomain(helloCommands, helloStates, helloFunctions, defaultResponses),
			},
			wantAnswers: []query.Answer{{
				Text: "Hey friend!",
			}},
			wantExtension: "",
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
				fsmDomain: fsm.NewDomain(pokemonCommands, pokemonStates, pokemonFunctions, defaultResponses),
			},
			wantAnswers:   nil,
			wantExtension: "search_pokemon",
			wantState:     0,
			wantSlots:     map[string]string{"pokemon": "pikachu"},
			wantErr:       false,
		},
		{
			name: "turn_on command should turn on",
			fields: fields{
				State: 0,
				Slots: make(map[string]string),
			},
			args: args{
				command:   "turn_on",
				txt:       "turn it on",
				fsmDomain: fsm.NewDomain(onOffCommands, onOffStates, onOffFunctions, defaultResponses),
			},
			wantAnswers: []query.Answer{{
				Text: "Turning on.",
			}},
			wantExtension: "",
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
				fsmDomain: fsm.NewDomain(onOffCommands, onOffStates, onOffFunctions, defaultResponses),
			},
			wantAnswers: []query.Answer{
				{
					Text: "Turning off.",
				},
				{
					Text: "❌",
				},
			},
			wantExtension: "",
			wantState:     0,
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
			if gotExtension != tt.wantExtension {
				t.Errorf("FSM.ExecuteCmd() gotExtension = %v, want %v", gotExtension, tt.wantExtension)
			}
			if !reflect.DeepEqual(m.Slots, tt.wantSlots) {
				t.Errorf("FSM.Slot gotSlot = %v, want %v", m.Slots, tt.wantSlots)
			}
		})
	}
}
