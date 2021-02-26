package extension_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jaimeteb/chatto/version"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

var (
	commandsPath = "/ext/commands"
	commandPath  = "/ext/command"
	versionPath  = "/ext/version"
	greetFunc    = func(req *extension.ExecuteCommandFuncRequest) (res *extension.ExecuteCommandFuncResponse) {
		return &extension.ExecuteCommandFuncResponse{
			FSM: req.FSM,
			Answers: []query.Answer{{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			}},
		}
	}
	registeredCommandFuncs = extension.RegisteredCommandFuncs{
		"any": greetFunc,
	}
)

func setBearerToken(r *http.Request, t string) *http.Request {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
	return r
}

func TestExtension_ListenerREST_GetBuildVersion(t *testing.T) {
	type args struct {
		l *extension.ListenerREST
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *version.BuildResponse
		wantErr *extension.ErrorResponse
	}{
		{
			name: "get build version",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: httptest.NewRequest("GET", versionPath, nil),
			},
			want: &version.BuildResponse{
				Version: "v0.0.0",
				Commit:  "0000000000000000000000000000000000000000",
				BuiltAt: "0001-01-01 00:00:00 +0000 UTC",
				BuiltBy: "dev",
			},
		},
		{
			name: "get build version with invalid http method",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: httptest.NewRequest("POST", versionPath, nil),
			},
			wantErr: &extension.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "got method 'POST', expected 'GET'",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.args.l.GetBuildVersion(w, tt.args.r)
			body, _ := ioutil.ReadAll(w.Result().Body)
			got := &version.BuildResponse{}
			_ = json.Unmarshal(body, &got)

			if w.Code != http.StatusOK {
				gotErr := &extension.ErrorResponse{}
				_ = json.Unmarshal(body, gotErr)
				if !reflect.DeepEqual(gotErr, tt.wantErr) {
					t.Errorf("Extension.ListenerREST.GetBuildVersion() error = %v, wantErr %v", spew.Sprint(gotErr), spew.Sprint(tt.wantErr))
					return
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extension.ListenerREST.GetBuildVersion() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}

func TestExtension_ListenerREST_GetAllCommandFuncs(t *testing.T) {
	type args struct {
		l *extension.ListenerREST
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr *extension.ErrorResponse
	}{
		{
			name: "get all funcs",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, ""),
				r: httptest.NewRequest("GET", commandsPath, nil),
			},
			want: []string{"any"},
		},
		{
			name: "get all funcs without auth",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: httptest.NewRequest("GET", commandsPath, nil),
			},
			wantErr: &extension.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "missing or incorrect authorization token",
			},
		},
		{
			name: "get all funcs with auth",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: setBearerToken(httptest.NewRequest("GET", commandsPath, nil), "my-test-token"),
			},
			want: []string{"any"},
		},
		{
			name: "get all funcs with invalid http method",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: setBearerToken(httptest.NewRequest("POST", commandsPath, nil), "my-test-token"),
			},
			wantErr: &extension.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "got method 'POST', expected 'GET'",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.args.l.GetAllCommandFuncs(w, tt.args.r)
			var got []string
			body, _ := ioutil.ReadAll(w.Result().Body)
			_ = json.Unmarshal(body, &got)

			if w.Code != http.StatusOK {
				gotErr := &extension.ErrorResponse{}
				_ = json.Unmarshal(body, gotErr)
				if !reflect.DeepEqual(gotErr, tt.wantErr) {
					t.Errorf("Extension.ListenerREST.GetAllCommandFuncs() error = %v, wantErr %v", spew.Sprint(gotErr), spew.Sprint(tt.wantErr))
					return
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extension.ListenerREST.GetAllCommandFuncs() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}

func TestExtension_ListenerREST_ExecuteCommandFunc(t *testing.T) {
	type args struct {
		l *extension.ListenerREST
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *extension.ExecuteCommandFuncResponse
		wantErr *extension.ErrorResponse
	}{
		{
			name: "execute command func",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, ""),
				r: httptest.NewRequest("POST", commandPath, bytes.NewBuffer([]byte(`{"command": "any", "fsm": {"state": 0, "slots": {}}}`))),
			},
			want: &extension.ExecuteCommandFuncResponse{
				FSM: &fsm.FSM{State: 0, Slots: fsm.NewFSM().Slots},
				Answers: []query.Answer{{
					Text:  "Hello Universe",
					Image: "https://i.imgur.com/pPdjh6x.jpg",
				}},
			},
		},
		{
			name: "execute command func with auth",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: setBearerToken(httptest.NewRequest("POST", commandPath, bytes.NewBuffer([]byte(`{"command": "any", "fsm": {"state": 0, "slots": {}}}`))), "my-test-token"),
			},
			want: &extension.ExecuteCommandFuncResponse{
				FSM: &fsm.FSM{State: 0, Slots: fsm.NewFSM().Slots},
				Answers: []query.Answer{{
					Text:  "Hello Universe",
					Image: "https://i.imgur.com/pPdjh6x.jpg",
				}},
			},
		},
		{
			name: "execute command func without auth",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: httptest.NewRequest("POST", commandPath, bytes.NewBuffer([]byte(`{"command": "any", "fsm": {"state": 0, "slots": {}}}`))),
			},
			want: &extension.ExecuteCommandFuncResponse{},
			wantErr: &extension.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "missing or incorrect authorization token",
			},
		},
		{
			name: "execute invalid command func",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: setBearerToken(httptest.NewRequest("POST", commandPath, bytes.NewBuffer([]byte(`{"command": "i_dont_exist", "fsm": {"state": 0, "slots": {}}}`))), "my-test-token"),
			},
			want: &extension.ExecuteCommandFuncResponse{},
			wantErr: &extension.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "extension command 'i_dont_exist' not found",
			},
		},
		{
			name: "execute command func with invalid http method",
			args: args{
				l: extension.NewListenerREST(registeredCommandFuncs, "my-test-token"),
				r: setBearerToken(httptest.NewRequest("GET", commandPath, bytes.NewBuffer([]byte(`{"command": "any", "fsm": {"state": 0, "slots": {}}}`))), "my-test-token"),
			},
			wantErr: &extension.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "got method 'GET', expected 'POST'",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.args.l.ExecuteCommandFunc(w, tt.args.r)
			got := &extension.ExecuteCommandFuncResponse{}
			body, _ := ioutil.ReadAll(w.Result().Body)
			_ = json.Unmarshal(body, got)

			if w.Code != http.StatusOK {
				gotErr := &extension.ErrorResponse{}
				_ = json.Unmarshal(body, gotErr)
				if !reflect.DeepEqual(gotErr, tt.wantErr) {
					t.Errorf("Extension.ListenerREST.ExecuteCommandFunc() error = %v, wantErr %v", spew.Sprint(gotErr), spew.Sprint(tt.wantErr))
					return
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extension.ListenerREST.ExecuteCommandFunc() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}

func TestExtensionRPCServer(t *testing.T) {
	listener := extension.ListenerRPC{RegisteredCommandFuncs: registeredCommandFuncs}

	err := listener.GetAllCommandFuncs(nil, new(extension.GetAllCommandFuncsResponse))
	if err != nil {
		t.Fatal(err)
	}

	req := extension.ExecuteCommandFuncRequest{
		Command: "any",
		FSM: &fsm.FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	}

	err = listener.ExecuteCommandFunc(&req, new(extension.ExecuteCommandFuncResponse))
	if err != nil {
		t.Fatal(err)
	}
}
