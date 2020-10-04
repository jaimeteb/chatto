package bot

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/kimrgrey/go-telegram"
)

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
	ans := Message{Sender: "botto", Text: resp}

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

	chatID := []string{sender}
	text := []string{resp}
	respValues := url.Values{
		"chat_id": chatID,
		"text":    text,
	}
	telegramClient := b.Endpoints["telegram"].(*telegram.Client)
	apiResp := new(interface{})
	telegramClient.Call("SendMessage", respValues, apiResp)
	log.Println(*apiResp)
}

func (b Bot) detailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	senderObj := fsm.FSM{
		State: b.Machines.GetState(vars["sender"]),
		// Slots: , TODO: retreive all slots
	}

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

	prediction, prob := b.Classifier.Predict(mess.Text)
	ans := Prediction{mess.Text, prediction, prob}

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
	domain := fsm.Create(path)
	classifier := clf.Create(path)

	extension, err := fsm.LoadExtension(path)
	if err != nil {
		log.Println(err)
	}

	endpoints := make(map[string]interface{})
	// TELEGRAM
	if telegramKey := os.Getenv("TELEGRAM_BOT_KEY"); telegramKey != "" {
		client := telegram.NewClient(telegramKey)
		endpoints["telegram"] = client

		log.Printf("Added Telegram client: %v\n", client.GetMe())
	}

	var machines fsm.StoreFSM
	// REDIS
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		machines = &fsm.RedisStoreFSM{R: fsm.RDB}
	} else {
		machines = &fsm.CacheStoreFSM{}
	}

	// machines := make(map[string]*fsm.FSM)
	bot := Bot{machines, domain, classifier, extension, endpoints}

	// log.Println("\n" + LOGO)
	log.Println("Server started")

	r := mux.NewRouter()
	r.HandleFunc("/endpoints/rest", bot.restEndpointHandler)
	r.HandleFunc("/endpoints/telegram", bot.telegramEndpointHandler)
	r.HandleFunc("/predict", bot.predictHandler)
	r.HandleFunc("/senders/{sender}", bot.detailsHandler)
	log.Fatal(http.ListenAndServe(":4770", r))
}
