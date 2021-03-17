package config

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
