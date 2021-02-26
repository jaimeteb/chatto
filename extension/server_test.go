package extension_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

func setBearerToken(r *http.Request, t string) *http.Request {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
	return r
}

func TestExtension_ListenerREST_GetAllFuncs(t *testing.T) {
	greetFunc := func(req *extension.Request) (res *extension.Response) {
		return &extension.Response{
			FSM: req.FSM,
			Answers: []query.Answer{{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			}},
		}
	}

	registeredCommandFuncs := extension.RegisteredCommandFuncs{
		"any": greetFunc,
	}

	type args struct {
		l *extension.ListenerREST
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "get all funcs",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, ""),
				r: httptest.NewRequest("GET", "/ext/get_all_funcs", nil),
			},
			want: []string{"any"},
		},
		{
			name: "get all funcs without auth",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: httptest.NewRequest("GET", "/ext/get_all_funcs", nil),
			},
			wantErr: true,
		},
		{
			name: "get all funcs with auth",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: setBearerToken(httptest.NewRequest("GET", "/ext/get_all_funcs", nil), "my-test-token"),
			},
			want: []string{"any"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.args.l.GetAllFuncs(w, tt.args.r)
			var got []string
			err := json.NewDecoder(w.Result().Body).Decode(&got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extension.ListenerREST.GetAllFuncs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extension.ListenerREST.GetAllFuncs() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}

	// req1, err := http.NewRequest("GET", "/ext/get_all_funcs", nil)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// w1 := httptest.NewRecorder()
	// listener.GetAllFuncs(w1, req1)

	// jsonStr2 := []byte(`{"extension": "any", "fsm": {"state": 0, "slots": {}}}`)
	// req2, err := http.NewRequest("POST", "/ext/get_func", bytes.NewBuffer(jsonStr2))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// w2 := httptest.NewRecorder()
	// listener.GetFunc(w2, req2)
}

func TestExtension_ListenerREST_GetFunc(t *testing.T) {
	greetFunc := func(req *extension.Request) (res *extension.Response) {
		return &extension.Response{
			FSM: req.FSM,
			Answers: []query.Answer{{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			}},
		}
	}

	registeredCommandFuncs := extension.RegisteredCommandFuncs{
		"any": greetFunc,
	}

	type args struct {
		l *extension.ListenerREST
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *extension.Response
		wantErr bool
	}{
		{
			name: "get func",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, ""),
				r: httptest.NewRequest("POST", "/ext/get_func", bytes.NewBuffer([]byte(`{"extension": "any", "fsm": {"state": 0, "slots": {}}}`))),
			},
			want: &extension.Response{
				FSM: &fsm.FSM{State: 0, Slots: fsm.NewFSM().Slots},
				Answers: []query.Answer{{
					Text:  "Hello Universe",
					Image: "https://i.imgur.com/pPdjh6x.jpg",
				}},
			},
		},
		{
			name: "get func with auth",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: setBearerToken(httptest.NewRequest("POST", "/ext/get_func", bytes.NewBuffer([]byte(`{"extension": "any", "fsm": {"state": 0, "slots": {}}}`))), "my-test-token"),
			},
			want: &extension.Response{
				FSM: &fsm.FSM{State: 0, Slots: fsm.NewFSM().Slots},
				Answers: []query.Answer{{
					Text:  "Hello Universe",
					Image: "https://i.imgur.com/pPdjh6x.jpg",
				}},
			},
		},
		{
			name: "get func without auth",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: httptest.NewRequest("POST", "/ext/get_func", nil),
			},
			want:    &extension.Response{nil, nil},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.args.l.GetFunc(w, tt.args.r)
			got := new(extension.Response)
			err := json.NewDecoder(w.Result().Body).Decode(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extension.ListenerREST.GetFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extension.ListenerREST.GetFunc() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}

func TestExtensionRPCServer(t *testing.T) {
	greetFunc := func(req *extension.Request) (res *extension.Response) {
		return &extension.Response{
			FSM: req.FSM,
			Answers: []query.Answer{{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			}},
		}
	}

	RegisteredCommandFuncs := extension.RegisteredCommandFuncs{
		"any": greetFunc,
	}

	listener := extension.ListenerRPC{RegisteredCommandFuncs}

	err := listener.GetAllFuncs(new(extension.Request), new(extension.GetAllFuncsResponse))
	if err != nil {
		t.Fatal(err)
	}

	req := extension.Request{
		Extension: "any",
		FSM: &fsm.FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	}

	err = listener.GetFunc(&req, new(extension.Response))
	if err != nil {
		t.Fatal(err)
	}
}
