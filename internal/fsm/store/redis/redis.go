package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/fsm/store/config"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

// Store struct models an FSM sotred on Redis
type Store struct {
	R   Client
	TTL time.Duration
}

type Client interface {
	Get(context.Context, string) *redis.StringCmd
	HGetAll(context.Context, string) *redis.StringStringMapCmd
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
	HSet(context.Context, string, ...interface{}) *redis.IntCmd
	Expire(context.Context, string, time.Duration) *redis.BoolCmd
}

func NewStore(cfg *config.StoreConfig) (*Store, error) {
	if cfg.Port == "" {
		cfg.Port = "6379"
	}
	var TLSConfig *tls.Config
	if cfg.TLS {
		TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	RDB := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:  cfg.Password,
		DB:        0,
		TLSConfig: TLSConfig,
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	log.Infof("* TTL:    %v", cfg.TTL)
	return &Store{R: RDB, TTL: cfg.TTL}, nil
}

// Exists for Store
func (s *Store) Exists(user string) (e bool) {
	_, err := s.R.Get(ctx, user+":state").Result()
	if err == redis.Nil || err != nil {
		return false
	}
	return true
}

// Get method for Store
func (s *Store) Get(user string) *fsm.FSM {
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

// Set method for Store
func (s *Store) Set(user string, m *fsm.FSM) {
	if err := s.R.Set(ctx, user+":state", m.State, s.TTL).Err(); err != nil {
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
		if err := s.R.Expire(ctx, user+":slots", s.TTL).Err(); err != nil {
			log.Error("Error expiring slots:", err)
		}
	}
}
