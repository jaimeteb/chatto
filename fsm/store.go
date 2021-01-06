package fsm

import (
	"context"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	redis "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// StoreConfig struct models a Store configuration in bot.yml
type StoreConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
}

// StoreFSM interface for FSM Store modes
type StoreFSM interface {
	Exists(string) bool
	Get(string) *FSM
	Set(string, *FSM)
}

// CacheStoreFSM struct models an FSM sotred in Cache
type CacheStoreFSM map[string]*FSM

// RedisStoreFSM struct models an FSM sotred on Redis
type RedisStoreFSM struct {
	R *redis.Client
}

// Exists for CacheStoreFSM
func (s *CacheStoreFSM) Exists(user string) (e bool) {
	_, ok := (*s)[user]
	return ok
}

// Exists for RedisStoreFSM
func (s *RedisStoreFSM) Exists(user string) (e bool) {
	_, err := s.R.Get(ctx, user+":state").Result()
	if err == redis.Nil || err != nil {
		return false
	}
	return true
}

// Get method for CacheStoreFSM
func (s *CacheStoreFSM) Get(user string) *FSM {
	return (*s)[user]
}

// Get method for RedisStoreFSM
func (s *RedisStoreFSM) Get(user string) *FSM {
	m := &FSM{}

	state, err := s.R.Get(ctx, user+":state").Result()
	if err != nil {
		log.Error(err)
	}
	i, err := strconv.Atoi(state)
	if err != nil {
		log.Error(err)
	}
	m.State = i

	slots, err := s.R.HGetAll(ctx, user+":slots").Result()
	if err != nil {
		log.Error(err)
	}
	m.Slots = slots

	return m
}

// Set method for CacheStoreFSM
func (s *CacheStoreFSM) Set(user string, m *FSM) {
	(*s)[user] = m
}

// Set method for RedisStoreFSM
func (s *RedisStoreFSM) Set(user string, m *FSM) {
	if err := s.R.Set(ctx, user+":state", m.State, 0).Err(); err != nil {
		log.Error("Error setting state:", err)
	}
	if len(m.Slots) > 0 {
		kvs := make([]string, 0)
		for k, v := range m.Slots {
			kvs = append(kvs, k, v)
		}

		if err := s.R.HSet(ctx, user+":slots", kvs).Err(); err != nil {
			log.Error("Error setting slots:", err)
		}
	}
}

// LoadStore loads a Store according to the configuration
func LoadStore(sc StoreConfig) StoreFSM {
	var machines StoreFSM
	switch sc.Type {
	case "REDIS":
		RDB := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:6379", sc.Host),
			Password: sc.Password,
			DB:       0,
		})
		if _, err := RDB.Ping(context.Background()).Result(); err != nil {
			machines = &CacheStoreFSM{}
			log.Warn("Couldn't connect to Redis, using CacheStoreFSM instead")
		} else {
			machines = &RedisStoreFSM{R: RDB}
			log.Info("Registered RedisStoreFSM")
		}
	default:
		machines = &CacheStoreFSM{}
		log.Info("Registered CacheStoreFSM")
	}
	return machines
}
