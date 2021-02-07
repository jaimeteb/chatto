package fsm

import (
	"testing"
)

func TestFSM1(t *testing.T) {
	path := "../examples/00_test/"
	domain := Create(&path)
	machine := FSM{State: 0}

	resp1, _ := machine.ExecuteCmd("turn_on", "turn_on", domain)
	if len(resp1) != 1 && resp1[0].Text != "Turning on." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp1, "Turning on.")
	}

	resp2, _ := machine.ExecuteCmd("turn_on", "turn_on", domain)
	if len(resp2) != 1 && resp2[0].Text != "Can't do that." {
		t.Errorf("resp is incorrect, got: %v, want: %v.", resp2, "Can't do that.")
	}

	resp5, _ := machine.ExecuteCmd("", "f o o", domain)
	if len(resp5) != 1 && resp5[0].Text != "???" {
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
