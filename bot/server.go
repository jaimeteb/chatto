package bot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/channels"
	"github.com/jaimeteb/chatto/channels/messages"
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
	receiveMsg, err := chnl.ReceiveMessage(w, r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (b *Bot) detailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	senderObj := b.Machines.Get(vars["sender"])

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

// ServeBot starts the bot process which starts the long running processes
func ServeBot(path *string, port *int) {
	bot, err := LoadBot(path)
	if err != nil {
		log.Fatal(err)
	}

	// log.Info("\n" + LOGO)
	log.Info("Server started")

	// Event listeners
	bot.slackMessageEvents()

	r := mux.NewRouter()

	// Integration Endpoints
	r.HandleFunc("/endpoints/rest", bot.restEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/telegram", bot.telegramEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/twilio", bot.twilioEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/slack", bot.slackEndpointHandler).Methods("POST")

	// Prediction and Sender Endpoints
	r.HandleFunc("/predict", bot.predictHandler).Methods("POST")
	r.HandleFunc("/senders/{sender}", bot.detailsHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), r))
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
