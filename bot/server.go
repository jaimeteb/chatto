package bot

import (
	"encoding/json"
	"fmt"
	"net/http"

	cmn "github.com/jaimeteb/chatto/common"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func (b Bot) restEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Clients.REST.RecieveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	if err := SendMessages(resp, &b.Clients.REST, mess.Sender, w); err != nil {
		log.Error(err)
		return
	}
}

func (b Bot) telegramEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Clients.Telegram.RecieveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	if err := SendMessages(resp, &b.Clients.Telegram, mess.Sender, w); err != nil {
		log.Error(err)
		return
	}
}

func (b Bot) twilioEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Clients.Twilio.RecieveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	if err := SendMessages(resp, &b.Clients.Twilio, mess.Sender, w); err != nil {
		log.Error(err)
		return
	}
}

func (b Bot) slackEndpointHandler(w http.ResponseWriter, r *http.Request) {
	mess, err := b.Clients.Slack.RecieveMessage(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	resp := b.Answer(mess)

	if err := SendMessages(resp, &b.Clients.Slack, mess.Sender, w); err != nil {
		log.Error(err)
		return
	}
}

func (b Bot) detailsHandler(w http.ResponseWriter, r *http.Request) {
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

func (b Bot) predictHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mess cmn.Message

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

// ServeBot function
func ServeBot(path *string, port *int) {
	bot := LoadBot(path)

	// log.Info("\n" + LOGO)
	log.Info("Server started")

	r := mux.NewRouter()

	// Integration Endpoints
	r.HandleFunc("/endpoints/rest", bot.restEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/telegram", bot.telegramEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/twilio", bot.twilioEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/slack", bot.slackEndpointHandler).Methods("POST")

	// Prediction and Sender Endpoints
	r.HandleFunc("/predict", bot.predictHandler).Methods("POST")
	r.HandleFunc("/senders/{sender}", bot.detailsHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), r))
}
