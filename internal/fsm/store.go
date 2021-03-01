package fsm

import (
	"context"
	"fmt"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

// StoreConfig struct models a Store configuration in bot.yml
type StoreConfig struct {
	Type     string `mapstructure:"type"`
	TTL      int    `mapstructure:"ttl"`
	Purge    int    `mapstructure:"purge"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

// Store interface for FSM Store modes
type Store interface {
	Exists(string) bool
	Get(string) *fsm.FSM
	Set(string, *fsm.FSM)
}

// CacheStore struct models an FSM sotred in Cache
type CacheStore struct {
	C *cache.Cache
}

// RedisStore struct models an FSM sotred on Redis
type RedisStore struct {
	R   *redis.Client
	TTL int
}

// Exists for CacheStoreFSM
func (s *CacheStore) Exists(user string) (e bool) {
	_, ok := s.C.Get(user)
	return ok
}

// Exists for RedisStoreFSM
func (s *RedisStore) Exists(user string) (e bool) {
	_, err := s.R.Get(ctx, user+":state").Result()
	if err == redis.Nil || err != nil {
		return false
	}
	return true
}

// Get method for CacheStoreFSM
func (s *CacheStore) Get(user string) *fsm.FSM {
	v, _ := s.C.Get(user)
	return v.(*fsm.FSM)
}

// Get method for RedisStoreFSM
func (s *RedisStore) Get(user string) *fsm.FSM {
	m := &fsm.FSM{}

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
func (s *CacheStore) Set(user string, m *fsm.FSM) {
	s.C.Set(user, m, 0)
}

// Set method for RedisStoreFSM
func (s *RedisStore) Set(user string, m *fsm.FSM) {
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

// NewStore loads a Store according to the configuration
func NewStore(storeConfig *StoreConfig) Store {
	var machines Store

	if storeConfig.TTL == 0 {
		storeConfig.TTL = -1
	}
	if storeConfig.Purge == 0 {
		if storeConfig.TTL != 0 {
			storeConfig.Purge = storeConfig.TTL
		} else {
			storeConfig.Purge = -1
		}
	}
	if storeConfig.Port == "" {
		storeConfig.Port = "6379"
	}

	switch storeConfig.Type {
	case "REDIS":
		RDB := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", storeConfig.Host, storeConfig.Port),
			Password: storeConfig.Password,
			DB:       0,
		})
		if _, err := RDB.Ping(context.Background()).Result(); err != nil {
			machines = &CacheStore{
				C: cache.New(
					time.Duration(storeConfig.TTL)*time.Second,
					time.Duration(storeConfig.Purge)*time.Second,
				),
			}
			log.Warn("Couldn't connect to Redis, using CacheStoreFSM instead")
			log.Infof("* TTL:    %v", storeConfig.TTL)
			log.Infof("* Purge:  %v", storeConfig.Purge)
		} else {
			machines = &RedisStore{R: RDB, TTL: storeConfig.TTL}
			log.Info("Registered RedisStoreFSM")
			log.Infof("* TTL:    %v", storeConfig.TTL)
		}
	default:
		machines = &CacheStore{
			C: cache.New(
				time.Duration(storeConfig.TTL)*time.Second,
				time.Duration(storeConfig.Purge)*time.Second,
			),
		}
		log.Info("Registered CacheStoreFSM")
		log.Infof("* TTL:    %v", storeConfig.TTL)
		log.Infof("* Purge:  %v", storeConfig.Purge)
	}
	return machines
}
