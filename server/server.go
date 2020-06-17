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

// func endpoint(w http.ResponseWriter, req *http.Request) {
// 	decoder := json.NewDecoder(req.Body)
// 	var mess models.Message

// 	err := decoder.Decode(&mess)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Println(mess.Sender, mess.Text)
// }

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

	bh.Bot.History.Messages = append(bh.Bot.History.Messages, mess)
	bh.Bot.History.Print()
}

// ServeBot function
func ServeBot(bot *models.Bot) {
	myBot := &botHandler{Bot: bot}

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/endpoint", myBot.handler)
	log.Fatal(http.ListenAndServe(":4770", nil))
}
