package extensions_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/extensions"
	"github.com/jaimeteb/chatto/internal/testutils"
	"github.com/jaimeteb/chatto/query"
)

func TestExtensionRESTError(t *testing.T) {
	extensionPort := testutils.GetFreePort(t)

	extCfg := make(map[string]extensions.Config)
	extCfg["test"] = extensions.Config{
		Type: "REST",
		URL:  fmt.Sprintf("http://localhost:%s", extensionPort),
	}

	extensions, err := extensions.New(extCfg)
	if err != nil {
		t.Errorf("extension.New() = %v, want %v.", err, nil)
	}

	if len(extensions) > 0 {
		t.Errorf("extension.New() = %v, want %v.", spew.Sprint(extensions), "map[]")
	}
}

func TestExtensionREST(t *testing.T) {
	extensionPort := testutils.GetFreePort(t)

	testutils.RunGoExtension(t, "../"+testutils.Examples00TestPath, extensionPort)

	extCfg := make(map[string]extensions.Config)
	extCfg["test"] = extensions.Config{
		Type: "REST",
		URL:  fmt.Sprintf("http://localhost:%s", extensionPort),
	}

	extensions, err := extensions.New(extCfg)
	if err != nil {
		t.Errorf("extension.New() = %v, want %v.", err, nil)
	}

	resp, err := extensions["test"].Execute(&query.Question{Text: "hello"}, "any", "", &fsm.Domain{}, &fsm.FSM{})
	if err != nil {
		t.Fatal(err)
	}

	want := "Hello Universe"

	if len(resp) == 1 && resp[0].Text != want {
		t.Errorf("extension.Execute() = %v, want %v.", resp[0].Text, want)
	}
}

func TestExtensionRPCPokemon(t *testing.T) {
	extensionPort := testutils.GetFreePort(t)

	testutils.RunGoExtension(t, "../"+testutils.Examples03PokemonPath, extensionPort)

	extPort, err := strconv.Atoi(extensionPort)
	if err != nil {
		t.Fatal(err)
	}

	extCfg := make(map[string]extensions.Config)
	extCfg["pokemon"] = extensions.Config{
		Type: "RPC",
		Host: "localhost",
		Port: extPort,
	}

	extensions, err := extensions.New(extCfg)
	if err != nil {
		t.Errorf("extension.New() = %v, want %v.", err, nil)
	}

	switch e := extensions["pokemon"].(type) {
	case *extensions.RPC:
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

	resp, err := extensions["pokemon"].Execute(&query.Question{Text: "pikachu"}, "search_pokemon", "", fsmDomain, &testFSM)
	if err != nil {
		t.Fatal(err)
	}

	want := `Name: pikachu 
ID: 25.00 
Height: 4.00 
Weight: 60.00`

	if len(resp) == 1 && resp[0].Text != want {
		t.Errorf("extension.Execute() = %v, want %v.", resp[0].Text, want)
	}
}

func TestExtensionRPCError(t *testing.T) {
	extensionPort := testutils.GetFreePort(t)

	extPort, err := strconv.Atoi(extensionPort)
	if err != nil {
		t.Fatal(err)
	}

	extCfg := make(map[string]extensions.Config)
	extCfg["test"] = extensions.Config{
		Type: "RPC",
		Host: "localhost",
		Port: extPort,
	}

	extensions, err := extensions.New(extCfg)
	if err != nil {
		t.Errorf("extension.New() = %v, want %v.", err, nil)
	}

	if len(extensions) > 0 {
		t.Errorf("extension.New() = %v, want %v.", spew.Sprint(extensions), "map[]")
	}
}
