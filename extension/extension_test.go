package extension

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

func TestRESTExt(t *testing.T) {
	extensionREST1, err := LoadExtensions(Config{
		Type: "REST",
		URL:  "http://localhost:8770",
	})
	if err != nil {
		t.Errorf("unable to load extensionREST1: %s", err)
	}

	extensionREST2, err := LoadExtensions(Config{
		Type: "REST",
		URL:  "http://localhost:8771",
	})
	if err != nil {
		t.Errorf("unable to load extensionREST2: %s", err)
	}

	resp1, err := extensionREST1.RunExtFunc(&query.Question{Text: "hello"}, "ext_any", &fsm.DB{}, &fsm.FSM{})
	if err != nil {
		t.Errorf("unable to run extensionREST1: %s", err)
	}

	if len(resp1) == 1 && resp1[0].Text != "Hello Universe" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp1, "Hello Universe")
	}

	fsmDB := &fsm.DB{}
	fsmDB.DefaultMessages = fsm.Defaults{Error: "Error"}

	resp2, err := extensionREST2.RunExtFunc(&query.Question{Text: "hello"}, "ext_any", fsmDB, &fsm.FSM{})
	if err != nil {
		t.Errorf("unable to run extensionREST2: %s", err)
	}

	if len(resp2) == 1 && resp2[0].Text != "Error" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp2, "Error")
	}
}

func TestRPCExt(t *testing.T) {
	extensionRPC1, err := LoadExtensions(Config{
		Type: "RPC",
		Host: "localhost",
		Port: 6770,
	})
	if err != nil {
		t.Errorf("unable to load extensionRPC1: %s", err)
	}

	switch e := extensionRPC1.(type) {
	case *RPC:
		break
	default:
		t.Errorf("incorrect, got %T, want: *ExtensionRPC", e)
	}

	extensionRPC2, err := LoadExtensions(Config{
		Type: "RPC",
		Host: "localhost",
		Port: 6771,
	})
	if err != nil {
		t.Errorf("unable to load extensionRPC2: %s", err)
	}

	switch extensionRPC2.(type) {
	case nil:
		break
	default:
		t.Error("incorrect, want: nil")
	}

	testDB := &fsm.DB{}
	testDB.DefaultMessages = fsm.Defaults{Error: "Error"}

	testFSM := fsm.FSM{
		Slots: map[string]string{
			"pokemon": "pikachu",
		},
	}

	resp1, err := extensionRPC1.RunExtFunc(&query.Question{Text: "pikachu"}, "ext_search_pokemon", testDB, &testFSM)
	if err != nil {
		t.Errorf("unable to run extensionRPC1: %s", err)
	}

	if len(resp1) == 1 && resp1[0].Text == "Error" {
		t.Errorf("resp is incorrect, got: %v", resp1)
	}

	resp2, err := extensionRPC2.RunExtFunc(&query.Question{Text: "hello"}, "ext_any", testDB, &fsm.FSM{})
	if err != nil {
		t.Errorf("unable to run extensionRPC2: %s", err)
	}

	if len(resp2) == 1 && resp2[0].Text != "Error" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp2, "Error")
	}
}

func TestRESTExtServer(t *testing.T) {
	greetFunc := func(req *Request) (res *Response) {
		return &Response{
			FSM: req.FSM,
			Answers: []query.Answer{{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			}},
		}
	}

	myExtMap := RegisteredFuncs{
		"ext_any": greetFunc,
	}

	listener := ListenerREST{myExtMap}

	req1, _ := http.NewRequest("GET", "/ext/get_all_funcs", nil)
	w1 := httptest.NewRecorder()
	listener.GetAllFuncs(w1, req1)

	jsonStr2 := []byte(`{"req": "ext_any", "fsm": {"state": 0, "slots": {}}}`)
	req2, _ := http.NewRequest("POST", "/ext/get_func", bytes.NewBuffer(jsonStr2))
	w2 := httptest.NewRecorder()
	listener.GetFunc(w2, req2)
}

func TestRPCExtServer(t *testing.T) {
	greetFunc := func(req *Request) (res *Response) {
		return &Response{
			FSM: req.FSM,
			Answers: []query.Answer{{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			}},
		}
	}

	myExtMap := RegisteredFuncs{
		"ext_any": greetFunc,
	}

	listener := ListenerRPC{myExtMap}

	err := listener.GetAllFuncs(new(Request), new(GetAllFuncsResponse))
	if err != nil {
		t.Fatal(err)
	}

	req := Request{
		Extension: "ext_any",
		FSM: &fsm.FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	}

	err = listener.GetFunc(&req, new(Response))
	if err != nil {
		t.Fatal(err)
	}
}
