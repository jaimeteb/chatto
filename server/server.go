package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jaimeteb/chatto/models"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

type botHandler struct {
	Bot *models.Bot
}

func (bh *botHandler) handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mess models.Message

	err := decoder.Decode(&mess)
	if err != nil {
		panic(err)
	}
	log.Println(mess.Sender, mess.Text)

	// ans := bh.Bot.Answer(mess)
	bh.Bot.Answer(mess)

	// bh.Bot.History.Messages[mess.Sender] = append(bh.Bot.History.Messages[mess.Sender], mess)
	// bh.Bot.History.Print(mess.Sender)
}

// ServeBot function
func ServeBot(bot *models.Bot) {
	bot.History.Messages = make(map[string][]models.Message)
	myBot := &botHandler{Bot: bot}

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/endpoint", myBot.handler)
	log.Fatal(http.ListenAndServe(":4770", nil))
}
