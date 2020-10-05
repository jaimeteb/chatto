package fsm

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RDB is the Redis client to be used in the chatbot
var RDB = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%v:6379", os.Getenv("REDIS_HOST")),
	Password: os.Getenv("REDIS_PASS"),
	DB:       0,
})

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
		log.Println(err)
	}
	i, err := strconv.Atoi(state)
	if err != nil {
		log.Println(err)
	}
	m.State = i

	slots, err := s.R.HGetAll(ctx, user+":slots").Result()
	if err != nil {
		log.Println(err)
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
		log.Println("Error setting state:", err)
	}
	if len(m.Slots) > 0 {
		if err := s.R.HSet(ctx, user+":slots", m.Slots, 0).Err(); err != nil {
			log.Println("Error setting slots:", err)
		}
	}
}
