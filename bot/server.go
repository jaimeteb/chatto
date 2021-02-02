package bot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/channels"
	"github.com/jaimeteb/chatto/message"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) restEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Channels.REST.ReceiveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	ans, err := channels.SendMessages(resp, b.Channels.REST, mess.Sender)
	if err != nil {
		log.Error(err)
		return
	}

	writeAnswer(w, ans)
}

func (b *Bot) telegramEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Channels.Telegram.ReceiveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	ans, err := channels.SendMessages(resp, b.Channels.Telegram, mess.Sender)
	if err != nil {
		log.Error(err)
		return
	}

	writeAnswer(w, ans)
}

func (b *Bot) twilioEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Channels.Twilio.ReceiveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	ans, err := channels.SendMessages(resp, b.Channels.Twilio, mess.Sender)
	if err != nil {
		log.Error(err)
		return
	}

	writeAnswer(w, ans)
}

func (b *Bot) slackEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Channels.Slack.ReceiveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	} else if (mess == message.Message{}) {
		return
	}

	resp := b.Answer(mess)

	ans, err := channels.SendMessages(resp, b.Channels.Slack, mess.Sender)
	if err != nil {
		log.Error(err)
		return
	}

	writeAnswer(w, ans)
}

func (b *Bot) slackMessageEvents() {
	messageChan := make(chan message.Message)

	go b.Channels.Slack.ReceiveMessages(messageChan)

	go func() {
		for mess := range messageChan {
			resp := b.Answer(mess)

			_, err := channels.SendMessages(resp, b.Channels.Slack, mess.Sender)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (b *Bot) predictHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mess message.Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	inputText := mess.Text
	prediction, prob := b.Classifier.Predict(inputText)
	ans := Prediction{inputText, prediction, prob}

	js, err := json.Marshal(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// ServeBot starts the bot process which starts the long running processes
func ServeBot(path *string, port *int) {
	bot := LoadBot(path)

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

func writeAnswer(w http.ResponseWriter, ans []map[string]string) {
	js, err := json.Marshal(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
