package bot

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/fsm"
)

func (b Bot) endpointHandler(w http.ResponseWriter, r *http.Request) {
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

func (b Bot) detailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	senderObj := b.Machines[vars["sender"]]

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

	extension := fsm.LoadExtension(path)

	machines := make(map[string]*fsm.FSM)
	bot := Bot{machines, domain, classifier, extension}

	log.Println("\n" + LOGO)
	log.Println("Server started")

	r := mux.NewRouter()
	r.HandleFunc("/endpoint", bot.endpointHandler)
	r.HandleFunc("/predict", bot.predictHandler)
	r.HandleFunc("/senders/{sender}", bot.detailsHandler)
	log.Fatal(http.ListenAndServe(":4770", r))
}
