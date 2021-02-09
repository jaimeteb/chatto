package fsm

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	redis "github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var mutex = &sync.RWMutex{}

// StoreConfig struct models a Store configuration in bot.yml
type StoreConfig struct {
	Type     string `mapstructure:"type"`
	TTL      int    `mapstructure:"ttl"`
	Purge    int    `mapstructure:"purge"`
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
type CacheStoreFSM struct {
	C *cache.Cache
}

// RedisStoreFSM struct models an FSM sotred on Redis
type RedisStoreFSM struct {
	R   *redis.Client
	TTL int
}

// Exists for CacheStoreFSM
func (s *CacheStoreFSM) Exists(user string) (e bool) {
	mutex.Lock()
	_, ok := s.C.Get(user)
	mutex.Unlock()
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
	mutex.Lock()
	v, _ := s.C.Get(user)
	mutex.Unlock()
	return v.(*FSM)
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
	mutex.Lock()
	s.C.Set(user, m, 0)
	mutex.Unlock()
}

// Set method for RedisStoreFSM
func (s *RedisStoreFSM) Set(user string, m *FSM) {
	if err := s.R.Set(ctx, user+":state", m.State, time.Duration(s.TTL)*time.Second).Err(); err != nil {
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
		if err := s.R.Expire(ctx, user+":slots", time.Duration(s.TTL)*time.Second).Err(); err != nil {
			log.Error("Error expiring slots:", err)
		}
	}
}

// LoadStore loads a Store according to the configuration
func LoadStore(sc StoreConfig) StoreFSM {
	var machines StoreFSM

	if sc.TTL == 0 {
		sc.TTL = -1
	}
	if sc.Purge == 0 {
		sc.Purge = -1
	}

	switch sc.Type {
	case "REDIS":
		RDB := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%v:6379", sc.Host),
			Password: sc.Password,
			DB:       0,
		})
		if _, err := RDB.Ping(context.Background()).Result(); err != nil {
			machines = &CacheStoreFSM{
				C: cache.New(
					time.Duration(sc.TTL)*time.Second,
					time.Duration(sc.Purge)*time.Second,
				),
			}
			log.Warn("Couldn't connect to Redis, using CacheStoreFSM instead")
			log.Infof("* TTL:    %v", sc.TTL)
			log.Infof("* Purge:  %v", sc.Purge)
		} else {
			machines = &RedisStoreFSM{R: RDB, TTL: sc.TTL}
			log.Info("Registered RedisStoreFSM")
			log.Infof("* TTL:    %v", sc.TTL)
		}
	default:
		machines = &CacheStoreFSM{
			C: cache.New(
				time.Duration(sc.TTL)*time.Second,
				time.Duration(sc.Purge)*time.Second,
			),
		}
		log.Info("Registered CacheStoreFSM")
		log.Infof("* TTL:    %v", sc.TTL)
		log.Infof("* Purge:  %v", sc.Purge)
	}
	return machines
}
