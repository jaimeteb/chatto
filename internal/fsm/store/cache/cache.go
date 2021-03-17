package cache

import (
	"time"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/fsm/store/config"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

// CacheStore struct models an FSM sotred in Cache
type CacheStore struct {
	C *cache.Cache
}

func NewCacheStore(cfg *config.StoreConfig) *CacheStore {
	log.Infof("* TTL:    %v", cfg.TTL)
	log.Infof("* Purge:  %v", cfg.Purge)
	return &CacheStore{
		C: cache.New(
			time.Duration(cfg.TTL)*time.Second,
			time.Duration(cfg.Purge)*time.Second,
		),
	}
}

// Exists for CacheStoreFSM
func (s *CacheStore) Exists(user string) (e bool) {
	_, ok := s.C.Get(user)
	return ok
}

// Get method for CacheStoreFSM
func (s *CacheStore) Get(user string) *fsm.FSM {
	v, _ := s.C.Get(user)
	return v.(*fsm.FSM)
}

// Set method for CacheStoreFSM
func (s *CacheStore) Set(user string, m *fsm.FSM) {
	s.C.Set(user, m, 0)
}
