package store_test

import (
	"reflect"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/jaimeteb/chatto/internal/fsm/store"
	"github.com/jaimeteb/chatto/internal/fsm/store/cache"
	"github.com/jaimeteb/chatto/internal/fsm/store/config"
	"github.com/jaimeteb/chatto/internal/fsm/store/redis"
	"github.com/jaimeteb/chatto/internal/fsm/store/sql"
	"github.com/jaimeteb/chatto/internal/testutils"
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

func TestNew(t *testing.T) {
	redisHost, redisPort := startRedisServer("pass")
	defer closeRedisServer()

	type args struct {
		cfg *config.StoreConfig
	}
	tests := []struct {
		name string
		args args
		want reflect.Type
	}{
		{
			name: "cache 1",
			args: args{
				cfg: &config.StoreConfig{},
			},
			want: reflect.TypeOf(&cache.Store{}),
		},
		{
			name: "cache 2",
			args: args{
				cfg: &config.StoreConfig{
					TTL: 10,
				},
			},
			want: reflect.TypeOf(&cache.Store{}),
		},
		{
			name: "redis success",
			args: args{
				cfg: &config.StoreConfig{
					Type:     "redis",
					Host:     redisHost,
					Port:     redisPort,
					Password: "pass",
				},
			},
			want: reflect.TypeOf(&redis.Store{}),
		},
		{
			name: "redis fail",
			args: args{
				cfg: &config.StoreConfig{
					Type:     "redis",
					Host:     redisHost,
					Port:     redisPort,
					Password: "passss",
				},
			},
			want: reflect.TypeOf(&cache.Store{}),
		},
		{
			name: "sql success",
			args: args{
				cfg: &config.StoreConfig{
					Type:     "sql",
					RDBMS:    "sqlite",
					Database: "test.db",
				},
			},
			want: reflect.TypeOf(&sql.Store{}),
		},
		{
			name: "sql fail mysql",
			args: args{
				cfg: &config.StoreConfig{
					Type:  "sql",
					RDBMS: "mysql",
				},
			},
			want: reflect.TypeOf(&cache.Store{}),
		},
		{
			name: "sql fail postgresql",
			args: args{
				cfg: &config.StoreConfig{
					Type:  "sql",
					RDBMS: "postgresql",
				},
			},
			want: reflect.TypeOf(&cache.Store{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reflect.TypeOf(store.New(tt.args.cfg))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
	t.Cleanup(func() {
		testutils.RemoveFiles("db")
	})
}
