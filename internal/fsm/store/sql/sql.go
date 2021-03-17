package sql

//go:generate mockgen -source=sql.go -destination=mocksql/mocksql.go -package=mocksql

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/fsm/store/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

// SQLStore models a SQL store for FSM
type SQLStore struct {
	DB DBClient
}

type DBClient interface {
	First(interface{}, ...interface{}) *gorm.DB
	Where(interface{}, ...interface{}) *gorm.DB
	Save(interface{}) *gorm.DB
}

func NewSQLStore(cfg *config.StoreConfig) (*SQLStore, error) {
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
	sqlStore.runPurge(cfg.TTL, cfg.Purge)

	return sqlStore, nil
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

func (s *SQLStore) runPurge(ttl, purge int) {
	if ttl > 0 && purge > 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(purge) * time.Second)
			for {
				select {
				case <-ticker.C:
					expired := time.Now().Add(-time.Duration(ttl) * time.Second)
					s.DB.Where("updated_at < ?", expired).Delete(&FSMORM{})
				}
			}
		}()
	}
}
