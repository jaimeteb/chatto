package core

import (
	"io/ioutil"
	"log"
	"strings"

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

// LoadConv loads conversations.md as Conversations
func LoadConv() (convs []models.Conversation) {
	content, err := ioutil.ReadFile("./config/conversations.md")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " ")

		switch {
		case strings.HasPrefix(trimmed, "## "):
			trimmed = strings.TrimPrefix(trimmed, "## ")
			newConv := models.Conversation{Name: trimmed}
			convs = append(convs, newConv)
		case strings.HasPrefix(trimmed, "* "):
			trimmed = strings.TrimPrefix(trimmed, "* ")
			mewMess := models.Message{
				Sender: "usr",
				Text:   trimmed,
			}
			convs[len(convs)-1].Path = append(convs[len(convs)-1].Path, mewMess)
		case strings.HasPrefix(trimmed, "- "):
			trimmed = strings.TrimPrefix(trimmed, "- ")
			mewMess := models.Message{
				Sender: "bot",
				Text:   trimmed,
			}
			convs[len(convs)-1].Path = append(convs[len(convs)-1].Path, mewMess)
		}
	}

	// fmt.Println(convs)
	return
}
