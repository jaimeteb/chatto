package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ajg/form"
	"github.com/gorilla/mux"
	"github.com/kimrgrey/go-telegram"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var chattoPort = getEnv("CHATTO_PORT", "4770")

func (b Bot) restEndpointHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mess Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// log.Println(mess.Sender, mess.Text)
	resp := b.Answer(mess)

	ans := make([]Message, 0)
	switch r := resp.(type) {
	case []interface{}:
		for _, text := range r {
			ans = append(ans, Message{Sender: b.Name, Text: text.(string)})
		}
	case interface{}:
		ans = append(ans, Message{Sender: b.Name, Text: r.(string)})
	default:
		errMsg := fmt.Sprintf("Message type unsupported: %T", r)
		http.Error(w, errors.New(errMsg).Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (b Bot) telegramEndpointHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var telegramMess TelegramMessageIn

	err := decoder.Decode(&telegramMess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(telegramMess)
	sender := strconv.Itoa(telegramMess.Message.From.ID)
	mess := Message{
		Sender: sender,
		Text:   telegramMess.Message.Text,
	}

	resp := b.Answer(mess)

	send := func(s, t string) {
		chatID := []string{s}
		text := []string{t}

		respValues := url.Values{
			"chat_id": chatID,
			"text":    text,
		}
		telegramClient := b.Clients["telegram"].(*telegram.Client)
		apiResp := new(interface{})
		telegramClient.Call("SendMessage", respValues, apiResp)
		log.Println(*apiResp)
	}

	switch r := resp.(type) {
	case []interface{}:
		for _, text := range r {
			send(sender, text.(string))
		}
	case interface{}:
		send(sender, r.(string))
	default:
		errMsg := fmt.Sprintf("Message type unsupported: %T", r)
		http.Error(w, errors.New(errMsg).Error(), http.StatusInternalServerError)
		return
	}
}

func (b Bot) twilioEndpointHandler(w http.ResponseWriter, r *http.Request) {
	decoder := form.NewDecoder(r.Body)
	var twilioMessage TwilioMessageIn
	if err := decoder.Decode(&twilioMessage); err != nil {
		http.Error(w, "Form could not be decoded", http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	log.Println(twilioMessage)
	sender := twilioMessage.From
	text := twilioMessage.Body
	mess := Message{
		Sender: sender,
		Text:   text,
	}

	resp := b.Answer(mess)

	send := func(s, t string) {
		twilio := b.Clients["twilio"].(Twilio)
		msg, err := twilio.Client.Messages.SendMessage(twilio.Number, s, t, nil) // TODO
		log.Println(msg, err)
	}

	switch r := resp.(type) {
	case []interface{}:
		for _, text := range r {
			send(sender, text.(string))
		}
	case interface{}:
		send(sender, r.(string))
	default:
		errMsg := fmt.Sprintf("Message type unsupported: %T", r)
		http.Error(w, errors.New(errMsg).Error(), http.StatusInternalServerError)
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
	var mess Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	inputText := mess.Text.(string)
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
func ServeBot(path *string) {
	bot := LoadBot(path)

	// log.Println("\n" + LOGO)
	log.Println("Server started")

	r := mux.NewRouter()
	r.HandleFunc("/endpoints/rest", bot.restEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/telegram", bot.telegramEndpointHandler).Methods("POST")
	r.HandleFunc("/endpoints/twilio", bot.twilioEndpointHandler).Methods("POST")
	r.HandleFunc("/predict", bot.predictHandler).Methods("POST")
	r.HandleFunc("/senders/{sender}", bot.detailsHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", chattoPort), r))
}
