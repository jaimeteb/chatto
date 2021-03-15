// nolint:bodyclose
package extensions_test

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
	"github.com/jaimeteb/chatto/extensions"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

var (
	commandsPath = "/extensions"
	commandPath  = "/extension"
	versionPath  = "/version"
	greetFunc    = func(req *extensions.ExecuteExtensionRequest) (res *extensions.ExecuteExtensionResponse) {
		return &extensions.ExecuteExtensionResponse{
			FSM: req.FSM,
			Answers: []query.Answer{{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			}},
		}
	}
	RegisteredExtensions = extensions.RegisteredExtensions{
		"any": greetFunc,
	}
)

func setBearerToken(r *http.Request, t string) *http.Request {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
	return r
}

func TestExtension_ListenerREST_GetBuildVersion(t *testing.T) {
	type args struct {
		l *extensions.ListenerREST
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *version.BuildResponse
		wantErr *extensions.ErrorResponse
	}{
		{
			name: "get build version",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, "my-test-token"),
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
				l: extensions.NewListenerREST(RegisteredExtensions, "my-test-token"),
				r: httptest.NewRequest("POST", versionPath, nil),
			},
			wantErr: &extensions.ErrorResponse{
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
				gotErr := &extensions.ErrorResponse{}
				_ = json.Unmarshal(body, gotErr)
				if !reflect.DeepEqual(gotErr, tt.wantErr) {
					t.Errorf("Extensions.ListenerREST.GetBuildVersion() error = %v, wantErr %v", spew.Sprint(gotErr), spew.Sprint(tt.wantErr))
					return
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extensions.ListenerREST.GetBuildVersion() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}

func TestExtension_ListenerREST_GetAllExtensions(t *testing.T) {
	type args struct {
		l *extensions.ListenerREST
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr *extensions.ErrorResponse
	}{
		{
			name: "get all funcs",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, ""),
				r: httptest.NewRequest("GET", commandsPath, nil),
			},
			want: []string{"any"},
		},
		{
			name: "get all funcs without auth",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, "my-test-token"),
				r: httptest.NewRequest("GET", commandsPath, nil),
			},
			wantErr: &extensions.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "missing or incorrect authorization token",
			},
		},
		{
			name: "get all funcs with auth",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, "my-test-token-abc"),
				r: setBearerToken(httptest.NewRequest("GET", commandsPath, nil), "my-test-token-abc"),
			},
			want: []string{"any"},
		},
		{
			name: "get all funcs with invalid http method",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, "my-test-token-123"),
				r: setBearerToken(httptest.NewRequest("POST", commandsPath, nil), "my-test-token-123"),
			},
			wantErr: &extensions.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "got method 'POST', expected 'GET'",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.args.l.GetAllExtensions(w, tt.args.r)
			var got []string
			body, _ := ioutil.ReadAll(w.Result().Body)
			_ = json.Unmarshal(body, &got)

			if w.Code != http.StatusOK {
				gotErr := &extensions.ErrorResponse{}
				_ = json.Unmarshal(body, gotErr)
				if !reflect.DeepEqual(gotErr, tt.wantErr) {
					t.Errorf("Extensions.ListenerREST.GetAllExtensions() error = %v, wantErr %v", spew.Sprint(gotErr), spew.Sprint(tt.wantErr))
					return
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extensions.ListenerREST.GetAllExtensions() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}

func TestExtension_ListenerREST_ExecuteExtension(t *testing.T) {
	type args struct {
		l *extensions.ListenerREST
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *extensions.ExecuteExtensionResponse
		wantErr *extensions.ErrorResponse
	}{
		{
			name: "execute command func",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, ""),
				r: httptest.NewRequest("POST", commandPath, bytes.NewBuffer([]byte(`{"extension": "any", "fsm": {"state": 0, "slots": {}}}`))),
			},
			want: &extensions.ExecuteExtensionResponse{
				FSM: &fsm.FSM{State: fsm.StateInitial, Slots: fsm.NewFSM().Slots},
				Answers: []query.Answer{{
					Text:  "Hello Universe",
					Image: "https://i.imgur.com/pPdjh6x.jpg",
				}},
			},
		},
		{
			name: "execute command func with auth",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, "my-test-token"),
				r: setBearerToken(httptest.NewRequest("POST", commandPath, bytes.NewBuffer([]byte(`{"extension": "any", "fsm": {"state": 0, "slots": {}}}`))), "my-test-token"),
			},
			want: &extensions.ExecuteExtensionResponse{
				FSM: &fsm.FSM{State: fsm.StateInitial, Slots: fsm.NewFSM().Slots},
				Answers: []query.Answer{{
					Text:  "Hello Universe",
					Image: "https://i.imgur.com/pPdjh6x.jpg",
				}},
			},
		},
		{
			name: "execute command func without auth",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, "my-test-token"),
				r: httptest.NewRequest("POST", commandPath, bytes.NewBuffer([]byte(`{"extension": "any", "fsm": {"state": 0, "slots": {}}}`))),
			},
			want: &extensions.ExecuteExtensionResponse{},
			wantErr: &extensions.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "missing or incorrect authorization token",
			},
		},
		{
			name: "execute invalid command func",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, "my-test-token-123"),
				r: setBearerToken(httptest.NewRequest("POST", commandPath, bytes.NewBuffer([]byte(`{"extension": "i_dont_exist", "fsm": {"state": 0, "slots": {}}}`))), "my-test-token-123"),
			},
			want: &extensions.ExecuteExtensionResponse{},
			wantErr: &extensions.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "extension 'i_dont_exist' not found",
			},
		},
		{
			name: "execute command func with invalid http method",
			args: args{
				l: extensions.NewListenerREST(RegisteredExtensions, "my-test-token-abc"),
				r: setBearerToken(httptest.NewRequest("GET", commandPath, bytes.NewBuffer([]byte(`{"extension": "any", "fsm": {"state": 0, "slots": {}}}`))), "my-test-token-abc"),
			},
			wantErr: &extensions.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "got method 'GET', expected 'POST'",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			tt.args.l.ExecuteExtension(w, tt.args.r)
			got := &extensions.ExecuteExtensionResponse{}
			body, _ := ioutil.ReadAll(w.Result().Body)
			_ = json.Unmarshal(body, got)

			if w.Code != http.StatusOK {
				gotErr := &extensions.ErrorResponse{}
				_ = json.Unmarshal(body, gotErr)
				if !reflect.DeepEqual(gotErr, tt.wantErr) {
					t.Errorf("Extensions.ListenerREST.ExecuteExtension() error = %v, wantErr %v", spew.Sprint(gotErr), spew.Sprint(tt.wantErr))
					return
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extensions.ListenerREST.ExecuteExtension() = %v, want %v", spew.Sprint(got), spew.Sprint(tt.want))
			}
		})
	}
}

func TestExtensionRPCServer(t *testing.T) {
	listener := extensions.ListenerRPC{RegisteredExtensions: RegisteredExtensions}

	err := listener.GetAllExtensions(nil, new(extensions.GetAllExtensionsResponse))
	if err != nil {
		t.Fatal(err)
	}

	req := extensions.ExecuteExtensionRequest{
		Extension: "any",
		FSM: &fsm.FSM{
			State: fsm.StateInitial,
			Slots: make(map[string]string),
		},
	}

	err = listener.ExecuteExtension(&req, new(extensions.ExecuteExtensionResponse))
	if err != nil {
		t.Fatal(err)
	}
}
