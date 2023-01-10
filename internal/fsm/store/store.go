package store

import (
	"strings"
	"time"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/fsm/store/cache"
	"github.com/jaimeteb/chatto/internal/fsm/store/config"
	"github.com/jaimeteb/chatto/internal/fsm/store/redis"
	"github.com/jaimeteb/chatto/internal/fsm/store/sql"
	log "github.com/sirupsen/logrus"
)

var defaultDuration = -1 * time.Second

// Store interface for FSM Store modes
type Store interface {
	Exists(string) bool
	Get(string) *fsm.FSM
	Set(string, *fsm.FSM)
}

// New loads a Store according to the configuration
func New(cfg *config.StoreConfig) Store {
	var machines Store

	if cfg.Purge == defaultDuration && cfg.TTL != defaultDuration {
		cfg.Purge = cfg.TTL
	}

	switch strings.ToLower(cfg.Type) {
	case "redis":
		redisStore, err := redis.NewStore(cfg)
		if err != nil {
			log.Errorf("Error: %v", err)
			log.Warn("Couldn't connect to Redis, using CacheStoreFSM instead")
			machines = cache.NewStore(cfg)
		} else {
			log.Info("Connected to RedisStoreFSM")
			machines = redisStore
		}
	case "sql":
		sqlStore, err := sql.NewStore(cfg)
		if err != nil {
			log.Errorf("Error: %v", err)
			log.Warn("Couldn't connect to SQL database, using CacheStoreFSM instead")
			machines = cache.NewStore(cfg)
		} else {
			log.Info("Connected to SQLStoreFSM")
			machines = sqlStore
		}
	default:
		log.Info("Connected to CacheStoreFSM")
		machines = cache.NewStore(cfg)
	}
	return machines
}
