package bot

import (
	"github.com/jaimeteb/chatto/internal/bot"
	log "github.com/sirupsen/logrus"
)

// Server runs botto bot
type Server struct {
	bot *bot.Bot
}

// NewServer for running botto bot
func NewServer(path string, port int) *Server {
	botConfig, err := bot.LoadConfig(path, port)
	if err != nil {
		log.Fatal(err)
	}

	b, err := bot.New(botConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{bot: b}
}

// Run botto bot server
func (s *Server) Run() {
	s.bot.Run()
}
