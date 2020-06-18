package core

import (
	"log"

	"github.com/jaimeteb/chatto/models"
	"github.com/spf13/viper"
)

// LoadYAML function
func LoadYAML() models.Bot {
	config := viper.New()
	config.AddConfigPath("./config")
	config.AddConfigPath(".")
	config.SetConfigName("bot")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	var bot models.Bot
	err = config.Unmarshal(&bot)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return bot
}

// LoadConv loads conversations.md and parses Conversations
func LoadConv() {

}
