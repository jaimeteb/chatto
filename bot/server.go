package bot

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jaimeteb/chatto/clf"
	"github.com/jaimeteb/chatto/fsm"
)

func (b Bot) handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mess Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(mess.Sender, mess.Text)
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

// ServeBot function
func ServeBot() {
	domain := fsm.Create()
	classifier := clf.Create()

	machines := make(map[string]*fsm.FSM)
	bot := Bot{machines, domain, classifier}

	http.HandleFunc("/endpoint", bot.handler)
	log.Fatal(http.ListenAndServe(":4770", nil))
}
