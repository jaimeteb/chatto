package bot

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/fsm"
)

// Bot models a bot with a Classifier and an FSM
type Bot struct {
	Machines   map[string]*fsm.FSM
	Domain     fsm.Domain
	Classifier clf.Classifier
}

// Answer takes a user input and executes a transition on the FSM if possible
func (b Bot) Answer(mess Message) string {
	if _, ok := b.Machines[mess.Sender]; !ok {
		b.Machines[mess.Sender] = &fsm.FSM{State: 0}
	}

	cmd, _ := b.Classifier.Predict(mess.Text) // Predict command from text using classifier
	return b.Machines[mess.Sender].ExecuteCmd(cmd, b.Domain)
}

func (b Bot) handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mess Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(mess.Sender, mess.Text)

	/////////////////////////////////////////////////////////////////////
	resp := b.Answer(mess)

	ans := Message{Sender: "botto", Text: resp}
	log.Println(b.Machines)
	/////////////////////////////////////////////////////////////////////

	js, err := json.Marshal(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// ServeBot function
func ServeBot() {
	// bot.History = make(map[string][]pkg.Message)
	// myBot := &botHandler{Bot: bot}
	domain := fsm.LoadDomain()
	classifier := clf.GetClassifier()

	machines := make(map[string]*fsm.FSM)
	bot := Bot{machines, domain, classifier}

	http.HandleFunc("/endpoint", bot.handler)
	log.Fatal(http.ListenAndServe(":4770", nil))
}
