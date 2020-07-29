package pkg

import (
	"encoding/json"
	"log"
	"net/http"
)

type botHandler struct {
	Machines map[string]FSM
	Domain   Domain
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
		bh.Machines[mess.Sender] = FSM{State: 0}
	}

	machine := bh.Machines[mess.Sender]
	resp := machine.ExecuteCmd(mess.Text, &bh.Domain)
	bh.Machines[mess.Sender] = machine

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
	domain := LoadDomain()

	machines := make(map[string]FSM)
	bot := botHandler{machines, domain}

	http.HandleFunc("/endpoint", bot.handler)
	log.Fatal(http.ListenAndServe(":4770", nil))
}
