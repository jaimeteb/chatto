package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jaimeteb/chatto/pkg"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

type botHandler struct {
	Bot *pkg.Bot
}

func (bh *botHandler) handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mess pkg.Message

	err := decoder.Decode(&mess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(mess.Sender, mess.Text)

	ans := bh.Bot.Answer(mess)

	js, err := json.Marshal(ans)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// ServeBot function
func ServeBot(bot *pkg.Bot) {
	bot.History.Messages = make(map[string][]pkg.Message)
	myBot := &botHandler{Bot: bot}

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/endpoint", myBot.handler)
	log.Fatal(http.ListenAndServe(":4770", nil))
}
