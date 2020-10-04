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
	SetState(string, int)
	GetState(string) int
	Exists(string) bool

	// SetSlot(string, string, interface{})
	// GetSlot(string, string) interface{}
}

// CacheStoreFSM struct models an FSM sotred in Cache
type CacheStoreFSM map[string]*FSM

// RedisStoreFSM struct models an FSM sotred on Redis
type RedisStoreFSM struct {
	R *redis.Client
}

// SetState for CacheStoreFSM
func (s *CacheStoreFSM) SetState(user string, i int) {
	// if s.Exists(user) {
	// 	(*s)[user].State = i
	// } else {
	// 	(*s)[user] = &FSM{
	// 		State: i,
	// 		Slots: make(map[string]interface{}),
	// 	}
	// }
	(*s)[user].State = i
}

// GetState for CacheStoreFSM
func (s *CacheStoreFSM) GetState(user string) (i int) {
	// if s.Exists(user) {
	// 	return (*s)[user].State
	// }
	// return -1
	return (*s)[user].State
}

// Exists for CacheStoreFSM
func (s *CacheStoreFSM) Exists(user string) (e bool) {
	_, ok := (*s)[user]
	return ok
}

// SetState for RedisStoreFSM
func (s *RedisStoreFSM) SetState(user string, i int) {
	if err := s.R.Set(ctx, user+":state", i, 0).Err(); err != nil {
		log.Println(err)
		// return err
	}
}

// GetState for RedisStoreFSM
func (s *RedisStoreFSM) GetState(user string) (i int) {
	val, err := s.R.Get(ctx, user+":state").Result()
	if err != nil {
		log.Println(err)
		return -1
	}
	i, err = strconv.Atoi(val)
	if err != nil {
		log.Println(err)
		return -1
	}
	return i
}

// Exists for RedisStoreFSM
func (s *RedisStoreFSM) Exists(user string) (e bool) {
	_, err := s.R.Get(ctx, user+":state").Result()
	if err == redis.Nil || err != nil {
		return false
	}
	return true
}
