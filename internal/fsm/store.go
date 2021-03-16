package fsm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	redis "github.com/go-redis/redis/v8"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

// FSMORM models a Finite State Machine with a gorm.Model
type FSMORM struct {
	gorm.Model
	User  string
	State int
	Slots string
}

func slotsToJsonString(slots map[string]string) string {
	if bytes, err := json.Marshal(slots); err != nil {
		log.Error(err)
		return "{}"
	} else {
		return string(bytes)
	}
}

func jsonStringToSlots(jsonStr string) map[string]string {
	var slots map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &slots); err != nil {
		log.Error(err)
		return make(map[string]string)
	}
	return slots
}

// StoreConfig struct models a Store configuration in bot.yml
type StoreConfig struct {
	Type     string `mapstructure:"type"`
	TTL      int    `mapstructure:"ttl"`
	Purge    int    `mapstructure:"purge"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	RDBMS    string `mapstructure:"rdbms"`
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

func NewCacheStore(cfg *StoreConfig) *CacheStore {
	log.Infof("* TTL:    %v", cfg.TTL)
	log.Infof("* Purge:  %v", cfg.Purge)
	return &CacheStore{
		C: cache.New(
			time.Duration(cfg.TTL)*time.Second,
			time.Duration(cfg.Purge)*time.Second,
		),
	}
}

// RedisStore struct models an FSM sotred on Redis
type RedisStore struct {
	R   *redis.Client
	TTL int
}

func NewRedisStore(cfg *StoreConfig) (*RedisStore, error) {
	if cfg.Port == "" {
		cfg.Port = "6379"
	}
	RDB := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	log.Infof("* TTL:    %v", cfg.TTL)
	return &RedisStore{R: RDB, TTL: cfg.TTL}, nil
}

// SQLStore models a SQL store for FSM
type SQLStore struct {
	DB *gorm.DB
}

func NewSQLStore(cfg *StoreConfig) (*SQLStore, error) {
	var db *gorm.DB
	var err error

	if cfg.Port == "" {
		cfg.Port = "3306"
	}
	// if cfg.Database == "" {
	// 	cfg.Database = "chatto"
	// }

	switch cfg.RDBMS {
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	case "postgresql":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.Host,
			cfg.User,
			cfg.Password,
			cfg.Database,
			cfg.Port,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.Database), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("No RDBMS specified for SQL connection.")
	}
	db.AutoMigrate(&FSMORM{})
	sqlStore := &SQLStore{db}
	if cfg.TTL > 0 && cfg.Purge > 0 {
		go sqlStore.runPurge(cfg.Purge, cfg.TTL)
	}
	return sqlStore, nil
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

// Exists for SQLStoreFSM
func (s *SQLStore) Exists(user string) (e bool) {
	machine := FSMORM{}
	if res := s.DB.First(&machine, "user = ?", user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

// Get method for SQLStoreFSM
func (s *SQLStore) Get(user string) *fsm.FSM {
	machine := FSMORM{}
	if res := s.DB.First(&machine, "user = ?", user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		return &fsm.FSM{
			State: machine.State,
			Slots: jsonStringToSlots(machine.Slots),
		}
	}
}

// Set method for SQLStoreFSM
func (s *SQLStore) Set(user string, m *fsm.FSM) {
	machine := FSMORM{}
	s.DB.First(&machine, "user = ?", user)
	machine.User = user
	machine.State = m.State
	machine.Slots = slotsToJsonString(m.Slots)
	if res := s.DB.Save(&machine); res.Error != nil {
		log.Error(res.Error)
	}
}

func (s *SQLStore) runPurge(purge, ttl int) {
	ticker := time.NewTicker(time.Duration(purge) * time.Second)
	for {
		select {
		case <-ticker.C:
			expired := time.Now().Add(-time.Duration(ttl) * time.Second)
			s.DB.Where("updated_at < ?", expired).Delete(&FSMORM{})
		}
	}
}

// Exists for RedisStoreFSM
func (s *RedisStore) Exists(user string) (e bool) {
	_, err := s.R.Get(ctx, user+":state").Result()
	if err == redis.Nil || err != nil {
		return false
	}
	return true
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
func NewStore(cfg *StoreConfig) Store {
	var machines Store

	if cfg.TTL == 0 {
		cfg.TTL = -1
	}
	if cfg.Purge == 0 {
		if cfg.TTL != 0 {
			cfg.Purge = cfg.TTL
		} else {
			cfg.Purge = -1
		}
	}

	switch strings.ToLower(cfg.Type) {
	case "redis":
		redisStore, err := NewRedisStore(cfg)
		if err != nil {
			log.Errorf("Error: %v", err)
			log.Warn("Couldn't connect to Redis, using CacheStoreFSM instead")
			machines = NewCacheStore(cfg)
		} else {
			log.Info("Connected to RedisStoreFSM")
			machines = redisStore
		}
	case "sql":
		sqlStore, err := NewSQLStore(cfg)
		if err != nil {
			log.Errorf("Error: %v", err)
			log.Warn("Couldn't connect to SQL database, using CacheStoreFSM instead")
			machines = NewCacheStore(cfg)
		} else {
			log.Info("Connected to SQLStoreFSM")
			machines = sqlStore
		}
	default:
		log.Info("Connected to CacheStoreFSM")
		machines = NewCacheStore(cfg)
	}
	return machines
}
