package extension_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

func TestExtensionRESTServer(t *testing.T) {
	greetFunc := func(req *extension.Request) (res *extension.Response) {
		return &extension.Response{
			FSM: req.FSM,
			Answers: []query.Answer{{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			}},
		}
	}

	registeredFuncs := extension.RegisteredFuncs{
		"any": greetFunc,
	}

	listener := extension.ListenerREST{RegisteredFuncs: registeredFuncs}

	req1, err := http.NewRequest("GET", "/ext/get_all_funcs", nil)
	if err != nil {
		t.Fatal(err)
	}

	w1 := httptest.NewRecorder()
	listener.GetAllFuncs(w1, req1)

	jsonStr2 := []byte(`{"extension": "any", "fsm": {"state": 0, "slots": {}}}`)
	req2, err := http.NewRequest("POST", "/ext/get_func", bytes.NewBuffer(jsonStr2))
	if err != nil {
		t.Fatal(err)
	}

	w2 := httptest.NewRecorder()
	listener.GetFunc(w2, req2)
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

	registeredFuncs := extension.RegisteredFuncs{
		"any": greetFunc,
	}

	listener := extension.ListenerRPC{registeredFuncs}

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
