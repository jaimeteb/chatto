package redis_test

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/fsm/store"
	"github.com/jaimeteb/chatto/internal/fsm/store/cache"
	"github.com/jaimeteb/chatto/internal/fsm/store/config"
	"github.com/jaimeteb/chatto/internal/fsm/store/redis"
)

var redisServer *miniredis.Miniredis = miniredis.NewMiniRedis()

func startRedisServer(pw string) (host, port string) {
	redisServer.RequireAuth(pw)
	if err := redisServer.Start(); err != nil {
		return "localhost", "6379"
	}
	return redisServer.Host(), redisServer.Port()
}

func closeRedisServer() {
	redisServer.Close()
}

func TestRedisStore(t *testing.T) {
	redisHost, redisPort := startRedisServer("pass")
	defer closeRedisServer()

	fmt.Println(redisServer.Addr())

	machines := store.New(&config.StoreConfig{
		Type:     "REDIS",
		Host:     redisHost,
		Port:     redisPort,
		Password: "pass",
	})

	switch machines.(type) {
	case *redis.Store:
		break
	default:
		t.Error("incorrect, want: *redis.Store")
	}

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

func TestRedisStoreFail(t *testing.T) {
	machines := store.New(&config.StoreConfig{
		Type:     "REDIS",
		Host:     "localhost",
		Password: "foo",
	})
	switch machines.(type) {
	case *cache.Store:
		break
	default:
		t.Error("incorrect, want: *cache.Store")
	}
}
