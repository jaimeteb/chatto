package extension_test

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
	"github.com/jaimeteb/chatto/testutils"
)

func TestExtensionREST(t *testing.T) {
	extensionPort := testutils.GetFreePort(t)

	testutils.RunGoExtension(t, testutils.Examples00TestPath, extensionPort)

	extensionREST, err := extension.New(extension.Config{
		Type: "REST",
		URL:  fmt.Sprintf("http://localhost:%s", extensionPort),
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := extensionREST.RunFunc(&query.Question{Text: "hello"}, "any", &fsm.Domain{}, &fsm.FSM{})
	if err != nil {
		t.Fatal(err)
	}

	want := "Hello Universe"

	if len(resp) == 1 && resp[0].Text != want {
		t.Errorf("extension.RunFunc() = %v, want %v.", resp[0].Text, want)
	}
}

func TestExtensionRESTError(t *testing.T) {
	extensionPort := testutils.GetFreePort(t)

	extensionREST, err := extension.New(extension.Config{
		Type: "REST",
		URL:  fmt.Sprintf("http://localhost:%s", extensionPort),
	})

	if err == nil {
		t.Errorf("extension.New() = %v, want %v.", nil, net.OpError{})
	}

	if extensionREST != nil {
		t.Errorf("extension.New() = %v, want %v.", spew.Sprint(extensionREST), nil)
	}
}

func TestExtensionRPCPokemon(t *testing.T) {
	extensionPort := testutils.GetFreePort(t)

	testutils.RunGoExtension(t, testutils.Examples03PokemonPath, extensionPort)

	extPort, err := strconv.Atoi(extensionPort)
	if err != nil {
		t.Fatal(err)
	}

	extensionRPC, err := extension.New(extension.Config{
		Type: "RPC",
		Host: "localhost",
		Port: extPort,
	})
	if err != nil {
		t.Fatal(err)
	}

	switch e := extensionRPC.(type) {
	case *extension.RPC:
		break
	default:
		t.Fatalf("incorrect, got %T, want: *ExtensionRPC", e)
	}

	fsmDomain := &fsm.Domain{}
	fsmDomain.DefaultMessages = fsm.Defaults{Error: "Error"}

	testFSM := fsm.FSM{
		Slots: map[string]string{
			"pokemon": "pikachu",
		},
	}

	resp, err := extensionRPC.RunFunc(&query.Question{Text: "pikachu"}, "search_pokemon", fsmDomain, &testFSM)
	if err != nil {
		t.Fatal(err)
	}

	want := `Name: pikachu 
ID: 25 
Height: 4 
Weight: 60`

	if len(resp) == 1 && resp[0].Text != want {
		t.Errorf("extension.RunFunc() = %v, want %v.", resp[0].Text, want)
	}
}

func TestExtensionRPCError(t *testing.T) {
	extensionPort := testutils.GetFreePort(t)

	extPort, err := strconv.Atoi(extensionPort)
	if err != nil {
		t.Fatal(err)
	}

	extensionRPC, err := extension.New(extension.Config{
		Type: "RPC",
		Host: "localhost",
		Port: extPort,
	})

	if err == nil {
		t.Errorf("extension.New() = %v, want %v.", nil, net.OpError{})
	}

	if extensionRPC != nil {
		t.Errorf("extension.New() = %v, want %v.", spew.Sprint(extensionRPC), nil)
	}
}

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

	myExtMap := extension.RegisteredFuncs{
		"any": greetFunc,
	}

	listener := extension.ListenerREST{myExtMap}

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

	myExtMap := extension.RegisteredFuncs{
		"any": greetFunc,
	}

	listener := extension.ListenerRPC{myExtMap}

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
