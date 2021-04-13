package cache

import (
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/fsm/store/config"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

// Store struct models an FSM sotred in Cache
type Store struct {
	C *cache.Cache
}

func NewStore(cfg *config.StoreConfig) *Store {
	log.Infof("* TTL:    %v", cfg.TTL)
	log.Infof("* Purge:  %v", cfg.Purge)
	return &Store{
		C: cache.New(
			cfg.TTL,
			cfg.Purge,
		),
	}
}

// Exists for Store
func (s *Store) Exists(user string) (e bool) {
	_, ok := s.C.Get(user)
	return ok
}

// Get method for Store
func (s *Store) Get(user string) *fsm.FSM {
	v, _ := s.C.Get(user)
	return v.(*fsm.FSM)
}

// Set method for Store
func (s *Store) Set(user string, m *fsm.FSM) {
	s.C.Set(user, m, 0)
}
