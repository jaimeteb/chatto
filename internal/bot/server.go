package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/internal/channels"
	"github.com/jaimeteb/chatto/internal/channels/messages"
	"github.com/jaimeteb/chatto/internal/channels/slack"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// ErrValidationFailed happens when a channel cannot validate an incoming callback
var ErrValidationFailed = errors.New("the callback token is invalid")

// Prediction models a classifier prediction and its original string
type Prediction struct {
	Original    string  `json:"original"`
	Predicted   string  `json:"predicted"`
	Probability float32 `json:"probability"`
}

func (b *Bot) restChannelHandler(w http.ResponseWriter, r *http.Request) {
	b.ChannelHandler(w, r, b.Channels.REST)
}

func (b *Bot) telegramChannelHandler(w http.ResponseWriter, r *http.Request) {
	b.ChannelHandler(w, r, b.Channels.Telegram)
}

func (b *Bot) twilioChannelHandler(w http.ResponseWriter, r *http.Request) {
	b.ChannelHandler(w, r, b.Channels.Twilio)
}

func (b *Bot) slackChannelHandler(w http.ResponseWriter, r *http.Request) {
	b.ChannelHandler(w, r, b.Channels.Slack)
}

func (b *Bot) healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ChannelHandler takes an incoming http.Request and passes it to a channel for it to respond
func (b *Bot) ChannelHandler(w http.ResponseWriter, r *http.Request, chnl channels.Channel) {
	if !chnl.ValidateCallback(r) {
		http.Error(w, ErrValidationFailed.Error(), http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	receiveMsg, err := chnl.ReceiveMessage(body)
	if err != nil {
		switch e := err.(type) {
		case slack.ErrURLVerification:
			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(e.Challenge)
			if err != nil {
				log.Error(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if receiveMsg.Question == nil || (*receiveMsg.Question == query.Question{}) {
		return
	}

	answers, err := b.Answer(receiveMsg)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = chnl.SendMessage(&messages.Response{Answers: answers, ReplyOpts: receiveMsg.ReplyOpts})
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeAnswer(w, answers)
}

func (b *Bot) slackChannelEvents() {
	if b.Channels.Slack != nil {
		receiveChan := make(chan messages.Receive)

		go b.Channels.Slack.ReceiveMessages(receiveChan)

		go func() {
			for receiveMsg := range receiveChan {
				r := receiveMsg

				answers, err := b.Answer(&r)
				if err != nil {
					log.Error(err)
					continue
				}

				err = b.Channels.Slack.SendMessage(&messages.Response{Answers: answers, ReplyOpts: receiveMsg.ReplyOpts})
				if err != nil {
					log.Error(err)
					continue
				}
			}
		}()
	}
}

func (b *Bot) detailsHandler(w http.ResponseWriter, r *http.Request) {
	if err := b.authorize(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	if vars == nil {
		log.Errorf("unable to get sender from request uri: %s", r.URL.RawPath)
		http.Error(w, "unable to get sender from request uri", http.StatusInternalServerError)
		return
	}

	if !b.Store.Exists(vars["sender"]) {
		log.Errorf("sender does not exist: %s", vars["sender"])
		http.Error(w, "sender does not exist", http.StatusNotFound)
		return
	}

	senderObj := b.Store.Get(vars["sender"])

	js, err := json.Marshal(senderObj)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (b *Bot) predictHandler(w http.ResponseWriter, r *http.Request) {
	if err := b.authorize(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var question query.Question

	err := decoder.Decode(&question)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inputText := question.Text
	prediction, prob := b.Classifier.Model.Predict(inputText, b.Classifier.Pipeline)
	answer := Prediction{inputText, prediction, prob}

	js, err := json.Marshal(answer)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (b *Bot) authorize(r *http.Request) error {
	if b.Config.Auth.Token != "" {
		reqToken := r.Header.Get("Authorization")
		reqToken = strings.TrimPrefix(reqToken, "Bearer ")

		if b.Config.Auth.Token != reqToken {
			return errors.New("unauthorized")
		}
	}
	return nil
}

// Run starts the bot which is a long running process
func (b *Bot) Run() {
	log.Info(smileyFace)
	log.Info("Bot started...")

	// Start event listeners
	b.slackChannelEvents()

	// Start web server
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", b.Config.Port),
		Handler:           b.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

// RegisterRoutes with the bot router
func (b *Bot) RegisterRoutes() {
	if b.Channels == nil {
		log.Warn("no channels configured, not registering routes")
		return
	}

	r := mux.NewRouter()

	// Channel channels
	if b.Channels.REST != nil {
		r.HandleFunc("/channels/rest", b.restChannelHandler).Methods("POST")
	}

	if b.Channels.Telegram != nil {
		r.HandleFunc("/channels/telegram", b.telegramChannelHandler).Methods("POST")
	}

	if b.Channels.Twilio != nil {
		r.HandleFunc("/channels/twilio", b.twilioChannelHandler).Methods("POST")
	}

	if b.Channels.Slack != nil {
		r.HandleFunc("/channels/slack", b.slackChannelHandler).Methods("POST")
	}

	// Other bot endpoints
	r.HandleFunc("/bot/healthz", b.healthzHandler).Methods("GET")
	r.HandleFunc("/bot/predict", b.predictHandler).Methods("POST")
	r.HandleFunc("/bot/senders/{sender}", b.detailsHandler).Methods("GET")

	b.Router = r
}

func writeAnswer(w http.ResponseWriter, answers []query.Answer) {
	js, err := json.Marshal(cleanAnswers(answers))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func cleanAnswers(answers []query.Answer) []map[string]string {
	finalAnswers := make([]map[string]string, len(answers))
	for i, answer := range answers {
		finalAnswers[i] = make(map[string]string)
		if answer.Text != "" {
			finalAnswers[i]["text"] = answer.Text
		}
		if answer.Image != "" {
			finalAnswers[i]["image"] = answer.Image
		}
	}
	return finalAnswers
}
