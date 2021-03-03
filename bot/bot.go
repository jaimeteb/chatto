package bot

import (
	"github.com/jaimeteb/chatto/internal/bot"
	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// Bot contains all the bot data and resources
type Bot struct {
	bot *bot.Bot
}

// Answer passes a querty.Question into the bot to produce answers
func (b *Bot) Answer(question *query.Question) (answers []query.Answer, err error) {
	message := messages.Receive{Question: question}
	answers, err = b.bot.Answer(&message)
	return
}

// New creates a new Bot
func New(path string) *Bot {
	botConfig, err := bot.LoadConfig(path, 0)
	if err != nil {
		log.Fatal(err)
	}

	newBot, err := bot.New(botConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &Bot{newBot}
}
