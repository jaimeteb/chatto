package fsm

import (
	"testing"

	"github.com/jaimeteb/chatto/ext"
)

func TestFSM1(t *testing.T) {
	path := "../examples/00_test/"
	domain := Create(&path)

	machine := FSM{State: 0}
	// extensionREST := ext.LoadExtensions(ext.ExtensionsConfig{
	// 	Type: "REST",
	// 	URL:  "http://localhost:8770",
	// })
	// extensionREST2 := ext.LoadExtensions(ext.ExtensionsConfig{
	// 	Type: "REST",
	// 	URL:  "http://localhost:8771",
	// })

	resp1, _ := machine.ExecuteCmd("turn_on", "turn_on", domain)
	if resp1 != "Turning on." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp1, "Turning on.")
	}

	resp2, _ := machine.ExecuteCmd("turn_on", "turn_on", domain)
	if resp2 != "Can't do that." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp2, "Can't do that.")
	}

	resp3, _ := machine.ExecuteCmd("hello_universe", "hello", domain)
	if resp3 != "Hello Universe" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp3, "Hello Universe")
	}

	resp4, _ := machine.ExecuteCmd("hello_universe", "hello", domain)
	if resp4 != "Error" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp4, "Error")
	}

	resp5, _ := machine.ExecuteCmd("", "f o o", domain)
	if resp5 != "???" {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp5, "???")
	}
}

func TestFSM2(t *testing.T) {
	path := "../examples/04_trivia/"
	domain := Create(&path)
	machine := FSM{
		State: 1,
		Slots: make(map[string]string),
	}
	machine.ExecuteCmd("start", "1", domain)
}

func TestCacheStore(t *testing.T) {
	machines := LoadStore(StoreConfig{Type: "CACHE"})

	if resp1 := machines.Exists("foo"); resp1 != false {
		t.Errorf("incorrect, got: %v, want: %v.", resp1, "false")
	}

	machines.Set(
		"foo",
		&FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	)
	if resp2 := machines.Get("foo"); resp2.State != 0 {
		t.Errorf("incorrect, got: %v, want: %v.", resp2, "0")
	}

	newFsm := &FSM{
		State: 1,
		Slots: map[string]string{
			"abc": "xyz",
		},
	}
	machines.Set("foo", newFsm)
	if resp3 := machines.Get("foo"); resp3.State != 1 {
		t.Errorf("incorrect, got: %v, want: %v.", resp3, "1")
	}
}

func TestRedisStore(t *testing.T) {
	machines := LoadStore(StoreConfig{
		Type:     "REDIS",
		Host:     "localhost",
		Password: "pass",
	})

	if resp1 := machines.Exists("foo"); resp1 != false {
		t.Errorf("incorrect, got: %v, want: %v.", resp1, "false")
	}

	machines.Set(
		"foo",
		&FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	)
	if resp2 := machines.Get("foo"); resp2.State != 0 {
		t.Errorf("incorrect, got: %v, want: %v.", resp2, "0")
	}

	newFsm := &FSM{
		State: 1,
		Slots: map[string]string{
			"abc": "xyz",
		},
	}
	machines.Set("foo", newFsm)
	if resp3 := machines.Get("foo"); resp3.State != 1 {
		t.Errorf("incorrect, got: %v, want: %v.", resp3, "1")
	}
}

func TestRedisStoreFail(t *testing.T) {
	machines := LoadStore(StoreConfig{
		Type:     "REDIS",
		Host:     "localhost",
		Password: "foo",
	})
	switch machines.(type) {
	case *CacheStoreFSM:
		break
	default:
		t.Error("incorrect, want: *CacheStoreFSM")
	}
}

func TestFSM3(t *testing.T) {
	extensionRPC := ext.LoadExtensions(ext.ExtensionsConfig{
		Type: "RPC",
		Host: "localhost",
		Port: 6770,
	})
	switch extensionRPC.(type) {
	case *ext.ExtensionRPC:
		break
	default:
		t.Error("incorrect, want: *ExtensionRPC")
	}

	path := "../examples/03_pokemon/"
	domain := Create(&path)
	machine := FSM{
		State: 1,
		Slots: make(map[string]string),
	}
	machine.ExecuteCmd("search_pokemon", "search ditto", domain)

	extensionRPC2 := ext.LoadExtensions(ext.ExtensionsConfig{
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
}
