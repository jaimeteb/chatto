package pkg

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// LoadYAML function
func LoadYAML() Bot {
	config := viper.New()
	config.AddConfigPath("./config")
	config.AddConfigPath(".")
	config.SetConfigName("bot")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	var bot Bot
	err = config.Unmarshal(&bot)
	if err != nil {
		log.Fatalf(err.Error())
	}

	bot.History.MaxHist = bot.MaxHist
	return bot
}

// LoadConv loads conversations.md as Conversations
func LoadConv() (convs []Conversation) {
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
			newConv := Conversation{Name: trimmed}
			convs = append(convs, newConv)
		case strings.HasPrefix(trimmed, "* "):
			trimmed = strings.TrimPrefix(trimmed, "* ")
			newMess := Message{
				Sender: "usr",
				Text:   trimmed,
			}
			convs[len(convs)-1].Path = append(convs[len(convs)-1].Path, newMess)
		case strings.HasPrefix(trimmed, "- "):
			trimmed = strings.TrimPrefix(trimmed, "- ")
			newMess := Message{
				Sender: "bot",
				Text:   trimmed,
			}
			convs[len(convs)-1].Path = append(convs[len(convs)-1].Path, newMess)
		}
	}

	// fmt.Println(convs)
	return
}
