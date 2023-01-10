package config

import "time"

// StoreConfig struct models a Store configuration in bot.yml
type StoreConfig struct {
	Type     string        `mapstructure:"type"`
	TTL      time.Duration `mapstructure:"ttl"`
	Purge    time.Duration `mapstructure:"purge"`
	Host     string        `mapstructure:"host"`
	Port     string        `mapstructure:"port"`
	User     string        `mapstructure:"user"`
	Password string        `mapstructure:"password"`
	Database string        `mapstructure:"database"`
	RDBMS    string        `mapstructure:"rdbms"`
}
