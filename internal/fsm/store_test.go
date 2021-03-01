package fsm_test

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/jaimeteb/chatto/fsm"
	fsmint "github.com/jaimeteb/chatto/internal/fsm"
)

var redisServer *miniredis.Miniredis = miniredis.NewMiniRedis()

func startRedisServer(pw string) (host, port string) {
	redisServer.RequireAuth(pw)
	redisServer.Start()
	return redisServer.Host(), redisServer.Port()
}

func closeRedisServer() {
	redisServer.Close()
}

func TestCacheStore(t *testing.T) {
	machines := fsmint.NewStore(fsmint.StoreConfig{Type: "CACHE"})

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
	redisHost, redisPort := startRedisServer("pass")
	defer closeRedisServer()

	fmt.Println(redisServer.Addr())

	machines := fsmint.NewStore(fsmint.StoreConfig{
		Type:     "REDIS",
		Host:     redisHost,
		Port:     redisPort,
		Password: "pass",
	})

	switch machines.(type) {
	case *fsmint.RedisStore:
		break
	default:
		t.Error("incorrect, want: *RedisStore")
	}

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
	machines := fsmint.NewStore(fsmint.StoreConfig{
		Type:     "REDIS",
		Host:     "localhost",
		Password: "foo",
	})
	switch machines.(type) {
	case *fsmint.CacheStore:
		break
	default:
		t.Error("incorrect, want: *CacheStoreFSM")
	}
}
