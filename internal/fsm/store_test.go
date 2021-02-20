package fsm_test

import (
	"testing"

	"github.com/jaimeteb/chatto/fsm"
	intfsm "github.com/jaimeteb/chatto/internal/fsm"
)

func TestCacheStore(t *testing.T) {
	machines := intfsm.NewStore(intfsm.StoreConfig{Type: "CACHE"})

	if resp1 := machines.Exists("foo"); resp1 != false {
		t.Errorf("incorrect, got: %v, want: %v.", resp1, "false")
	}

	machines.Set(
		"foo",
		&fsm.FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	)
	if resp2 := machines.Get("foo"); resp2.State != 0 {
		t.Errorf("incorrect, got: %v, want: %v.", resp2, "0")
	}

	newFsm := &fsm.FSM{
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
	machines := intfsm.NewStore(intfsm.StoreConfig{
		Type:     "REDIS",
		Host:     "localhost",
		Password: "pass",
	})

	if resp1 := machines.Exists("foo"); resp1 != false {
		t.Errorf("incorrect, got: %v, want: %v.", resp1, "false")
	}

	machines.Set(
		"foo",
		&fsm.FSM{
			State: 0,
			Slots: make(map[string]string),
		},
	)
	if resp2 := machines.Get("foo"); resp2.State != 0 {
		t.Errorf("incorrect, got: %v, want: %v.", resp2, "0")
	}

	newFsm := &fsm.FSM{
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
	machines := intfsm.NewStore(intfsm.StoreConfig{
		Type:     "REDIS",
		Host:     "localhost",
		Password: "foo",
	})
	switch machines.(type) {
	case *intfsm.CacheStore:
		break
	default:
		t.Error("incorrect, want: *CacheStoreFSM")
	}
}
