package sql

//go:generate mockgen -source=sql.go -destination=mocksql/mocksql.go -package=mocksql

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/internal/fsm/store/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var userCol string = "user"

// FSMORM models a Finite State Machine with a gorm.Model
type FSMORM struct {
	gorm.Model
	User  string
	State int
	Slots string
}

func (_ *FSMORM) TableName() string {
	return "fsms"
}

func slotsToJSONString(slots map[string]string) string {
	bytes, err := json.Marshal(slots)
	if err != nil {
		log.Error(err)
		return "{}"
	}
	return string(bytes)
}

func jsonStringToSlots(jsonStr string) map[string]string {
	var slots map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &slots); err != nil {
		log.Error(err)
		return make(map[string]string)
	}
	return slots
}

// Store models a SQL store for FSM
type Store struct {
	DB DBClient
}

type DBClient interface {
	First(interface{}, ...interface{}) *gorm.DB
	Where(interface{}, ...interface{}) *gorm.DB
	Save(interface{}) *gorm.DB
}

func NewStore(cfg *config.StoreConfig) (*Store, error) {
	var db *gorm.DB
	var err error

	if cfg.Database == "" {
		cfg.Database = "chatto"
	}

	switch strings.ToLower(cfg.RDBMS) {
	case "mysql":
		if cfg.Port == "" {
			cfg.Port = "3306"
		}

		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return nil, err
		}
	case "postgresql", "postgres":
		if cfg.Port == "" {
			cfg.Port = "5432"
		}

		userCol = fmt.Sprintf("%s.%s", new(FSMORM).TableName(), userCol)
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.Host,
			cfg.User,
			cfg.Password,
			cfg.Database,
			cfg.Port,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return nil, err
		}
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.Database), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("no RDBMS specified for SQL connection")
	}

	if err := db.AutoMigrate(&FSMORM{}); err != nil {
		log.Error(err)
	}

	sqlStore := &Store{db}
	sqlStore.runPurge(cfg.TTL, cfg.Purge)

	return sqlStore, nil
}

// Exists for Store
func (s *Store) Exists(user string) (e bool) {
	machine := FSMORM{}
	if res := s.DB.First(&machine, fmt.Sprintf("%s = ?", userCol), user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Debug(res.Error)
		return false
	}
	return true
}

// Get method for Store
func (s *Store) Get(user string) *fsm.FSM {
	machine := FSMORM{}
	if res := s.DB.First(&machine, fmt.Sprintf("%s = ?", userCol), user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Debug(res.Error)
		return nil
	}
	return &fsm.FSM{
		State: machine.State,
		Slots: jsonStringToSlots(machine.Slots),
	}
}

// Set method for Store
func (s *Store) Set(user string, m *fsm.FSM) {
	machine := FSMORM{}
	s.DB.First(&machine, "user = ?", user)
	machine.User = user
	machine.State = m.State
	machine.Slots = slotsToJSONString(m.Slots)
	if res := s.DB.Save(&machine); res.Error != nil {
		log.Error(res.Error)
	}
}

func (s *Store) runPurge(ttl, purge int) {
	if ttl > 0 && purge > 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(purge) * time.Second)
			for {
				for range ticker.C {
					expired := time.Now().Add(-time.Duration(ttl) * time.Second)
					s.DB.Where("updated_at < ?", expired).Delete(&FSMORM{})
				}
			}
		}()
	}
}
