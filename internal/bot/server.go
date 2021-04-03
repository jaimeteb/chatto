package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/internal/channels"
	"github.com/jaimeteb/chatto/internal/channels/message"
	"github.com/jaimeteb/chatto/internal/channels/slack"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

// ErrValidationFailed happens when a channel cannot validate an incoming callback
var ErrValidationFailed error = errors.New("the callback token is invalid")

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

	body, err := ioutil.ReadAll(r.Body)
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

	err = b.SubmitMessageRequest(receiveMsg)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("message request submitted"))
	if err != nil {
		log.Error(err)
	}
}

func (b *Bot) slackEventHandler() {
	receiveChan := make(chan message.Request)
	b.ChannelEventHandler(b.Channels.Slack, receiveChan)
}

// ChannelEventHandler takes a message.Request event and passes it to the bot
func (b *Bot) ChannelEventHandler(chnl channels.Channel, msgRequest chan message.Request) {
	if chnl == nil {
		return
	}

	go chnl.ReceiveMessages(msgRequest)

	go func() {
		for receiveMsg := range msgRequest {
			r := receiveMsg

			err := b.SubmitMessageRequest(&r)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}()
}

func (b *Bot) messageResponseEventHandler() {
	b.MessageResponseQueue = make(chan message.Response)

	go func() {
		for messageResponse := range b.MessageResponseQueue {
			res := messageResponse

			chnl := b.Channels.Get(res.Channel)
			if chnl == nil {
				log.Errorf("channel name is invalid: %s", res.Channel)
				continue
			}

			err := chnl.SendMessage(&res)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}()
}

func (b *Bot) messageResponseHandler(w http.ResponseWriter, r *http.Request) {
	if err := b.authorize(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var res message.Response

	err := decoder.Decode(&res)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chnl := b.Channels.Get(res.Channel)
	if chnl == nil {
		errMsg := fmt.Sprintf("channel name is invalid: %s", res.Channel)
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	err = chnl.SendMessage(&res)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func (b *Bot) sendersHandler(w http.ResponseWriter, r *http.Request) {
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

	// Start event handlers
	b.slackEventHandler()
	b.messageResponseEventHandler()

	// Start web server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", b.Config.Port), b.Router))
}

// RegisterRoutes with the bot router
func (b *Bot) RegisterRoutes() {
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

	// Bot endpoints
	r.HandleFunc("/bot/healthz", b.healthzHandler).Methods("GET")
	r.HandleFunc("/bot/message/response", b.messageResponseHandler).Methods("POST")
	r.HandleFunc("/bot/predict", b.predictHandler).Methods("POST")
	r.HandleFunc("/bot/senders/{sender}", b.sendersHandler).Methods("GET")
	r.HandleFunc("/bot/ws", b.WebSocket.ExtensionWebsocketHandler)

	b.Router = r
}
