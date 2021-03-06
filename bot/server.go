package bot

import (
	"net/http"

	"github.com/jaimeteb/chatto/internal/bot"
	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/query"
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

// Answer passes a querty.Question into the bot to produce answers
func (s *Server) Answer(question *query.Question) (answers []query.Answer, err error) {
	message := messages.Receive{Question: question}
	answers, err = s.bot.Answer(&message)
	return
}

// RESTHandler passes an incoming http.Request to the REST channel
func (s *Server) RESTHandler(w http.ResponseWriter, r *http.Request) {
	s.bot.ChannelHandler(w, r, s.bot.Channels.REST)
}

// TelegramHandler passes an incoming http.Request to the Telegram channel
func (s *Server) TelegramHandler(w http.ResponseWriter, r *http.Request) {
	s.bot.ChannelHandler(w, r, s.bot.Channels.Telegram)
}

// TwilioHandler passes an incoming http.Request to the Twilio channel
func (s *Server) TwilioHandler(w http.ResponseWriter, r *http.Request) {
	s.bot.ChannelHandler(w, r, s.bot.Channels.Twilio)
}

// SlackHandler passes an incoming http.Request to the Slack channel
func (s *Server) SlackHandler(w http.ResponseWriter, r *http.Request) {
	s.bot.ChannelHandler(w, r, s.bot.Channels.Slack)
}

// Run botto bot server
func (s *Server) Run() {
	s.bot.Run()
}
