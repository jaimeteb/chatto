package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

type botHandler struct {
	Machines map[string]FSM
}

func (bh botHandler) handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mess Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(mess.Sender, mess.Text)

	/////////////////////////////////////////////////////////////////////
	if _, ok := bh.Machines[mess.Sender]; !ok {
		bh.Machines[mess.Sender] = FSM{State: Initial}
	}

	x := bh.Machines[mess.Sender]
	resp := x.ExecuteCmd(mess.Text)
	bh.Machines[mess.Sender] = x

	ans := Message{Sender: "botto", Text: resp}
	log.Println(bh.Machines)
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

	machines := make(map[string]FSM)
	bot := botHandler{Machines: machines}

	http.HandleFunc("/endpoint", bot.handler)
	log.Fatal(http.ListenAndServe(":4770", nil))
}
