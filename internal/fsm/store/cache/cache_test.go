package cache_test

import (
	"testing"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/fsm/store"
	"github.com/jaimeteb/chatto/internal/fsm/store/config"
)

func TestCacheStore(t *testing.T) {
	machines := store.New(&config.StoreConfig{Type: "CACHE"})

	if resp1 := machines.Exists("foo"); resp1 != false {
		t.Errorf("incorrect, got: %v, want: %v.", resp1, "false")
	}

	machines.Set(
		"foo",
		&fsm.FSM{
			State: fsm.StateInitial,
			Slots: make(map[string]string),
		},
	)
	if resp2 := machines.Get("foo"); resp2.State != fsm.StateInitial {
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
