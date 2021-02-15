package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/channels"
	"github.com/jaimeteb/chatto/channels/messages"
	"github.com/jaimeteb/chatto/channels/slack"
	"github.com/jaimeteb/chatto/query"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) restEndpointHandler(w http.ResponseWriter, r *http.Request) {
	b.endpointHandler(w, r, b.Channels.REST)
}

func (b *Bot) telegramEndpointHandler(w http.ResponseWriter, r *http.Request) {
	b.endpointHandler(w, r, b.Channels.Telegram)
}

func (b *Bot) twilioEndpointHandler(w http.ResponseWriter, r *http.Request) {
	b.endpointHandler(w, r, b.Channels.Twilio)
}

func (b *Bot) slackEndpointHandler(w http.ResponseWriter, r *http.Request) {
	b.endpointHandler(w, r, b.Channels.Slack)
}

func (b *Bot) endpointHandler(w http.ResponseWriter, r *http.Request, chnl channels.Channel) {
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

	answers, err := b.Answer(receiveMsg.Question)
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

func (b *Bot) slackMessageEvents() {
	if b.Channels.Slack != nil {
		receiveChan := make(chan messages.Receive)

		go b.Channels.Slack.ReceiveMessages(receiveChan)

		go func() {
			for receiveMsg := range receiveChan {
				answers, err := b.Answer(receiveMsg.Question)
				if err != nil {
					log.Error(err)
					return
				}

				err = b.Channels.Slack.SendMessage(&messages.Response{Answers: answers, ReplyOpts: receiveMsg.ReplyOpts})
				if err != nil {
					log.Error(err)
					return
				}
			}
		}()
	}
}

func (b *Bot) detailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars == nil {
		log.Errorf("unable to get sender from request uri: %s", r.URL.RawPath)
		http.Error(w, "unable to get sender from request uri", http.StatusInternalServerError)
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
	decoder := json.NewDecoder(r.Body)

	var question query.Question

	err := decoder.Decode(&question)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inputText := question.Text
	prediction, prob := b.Classifier.Predict(inputText)
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

// Run starts the bot which is a long running process
func (b *Bot) Run() {
	// log.Info("\n" + LOGO)
	log.Info("Bot started...")

	// Start event listeners
	b.slackMessageEvents()

	// Start web server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", b.Config.Port), b.Router))
}

// RegisterRoutes with the bot router
func (b *Bot) RegisterRoutes() {
	if b.Channels == nil {
		log.Warn("no channels configured, not registering routes")
		return
	}

	r := mux.NewRouter()

	// Channel endpoints
	if b.Channels.REST != nil {
		r.HandleFunc("/endpoints/rest", b.restEndpointHandler).Methods("POST")
	}

	if b.Channels.Telegram != nil {
		r.HandleFunc("/endpoints/telegram", b.telegramEndpointHandler).Methods("POST")
	}

	if b.Channels.Twilio != nil {
		r.HandleFunc("/endpoints/twilio", b.twilioEndpointHandler).Methods("POST")
	}

	if b.Channels.Slack != nil {
		r.HandleFunc("/endpoints/slack", b.slackEndpointHandler).Methods("POST")
	}

	// Prediction and Sender Endpoints
	r.HandleFunc("/predict", b.predictHandler).Methods("POST")
	r.HandleFunc("/senders/{sender}", b.detailsHandler).Methods("GET")

	b.Router = r
}

func writeAnswer(w http.ResponseWriter, answers []query.Answer) {
	js, err := json.Marshal(answers)
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
