package ext

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/message"
)

func TestRESTExt(t *testing.T) {
	extensionREST1 := LoadExtensions(ExtensionsConfig{
		Type: "REST",
		URL:  "http://localhost:8770",
	})
	extensionREST2 := LoadExtensions(ExtensionsConfig{
		Type: "REST",
		URL:  "http://localhost:8771",
	})

	resp1 := extensionREST1.RunExtFunc("", "ext_any", "hello", fsm.Domain{}, &fsm.FSM{})
	if resp1.(map[string]interface{})["text"] != "Hello Universe" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp1, "Hello Universe")
	}

	resp2 := extensionREST2.RunExtFunc("", "ext_any", "hello", fsm.Domain{DefaultMessages: fsm.Defaults{Error: "Error"}}, &fsm.FSM{})
	if resp2.(string) != "Error" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp2, "Error")
	}
}

func TestRPCExt(t *testing.T) {
	extensionRPC1 := LoadExtensions(ExtensionsConfig{
		Type: "RPC",
		Host: "localhost",
		Port: 6770,
	})
	switch e := extensionRPC1.(type) {
	case *ExtensionRPC:
		break
	default:
		t.Errorf("incorrect, got %T, want: *ExtensionRPC", e)
	}

	extensionRPC2 := LoadExtensions(ExtensionsConfig{
		Type: "RPC",
		Host: "localhost",
		Port: 6771,
	})
	switch extensionRPC2.(type) {
	case nil:
		break
	default:
		t.Error("incorrect, want: nil")
	}

	testDom := fsm.Domain{
		DefaultMessages: fsm.Defaults{Error: "Error"},
	}
	testFSM := fsm.FSM{
		Slots: map[string]string{
			"pokemon": "pikachu",
		},
	}
	resp1 := extensionRPC1.RunExtFunc("", "ext_search_pokemon", "pikachu", testDom, &testFSM)
	if resp1.(string) == "Error" {
		t.Errorf("resp is incorrect, got: %v", resp1)
	}

	resp2 := extensionRPC1.RunExtFunc("", "ext_any", "hello", fsm.Domain{DefaultMessages: fsm.Defaults{Error: "Error"}}, &fsm.FSM{})
	if resp2.(string) != "Error" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp2, "Error")
	}
}

func TestRESTExtServer(t *testing.T) {
	greetFunc := func(req *Request) (res *Response) {
		return &Response{
			FSM: req.FSM,
			Res: message.Message{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			},
		}
	}

	myExtMap := ExtensionMap{
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
			Res: message.Message{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			},
		}
	}

	myExtMap := ExtensionMap{
		"ext_any": greetFunc,
	}

	listener := ListenerRPC{myExtMap}

	listener.GetAllFuncs(new(Request), new(GetAllFuncsResponse))

	req := Request{
		Req: "ext_any",
		FSM: &fsm.FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	}
	listener.GetFunc(&req, new(Response))
}
