package main

import (
	"fmt"

	"github.com/jaimeteb/chatto/models"
	"github.com/jaimeteb/chatto/server"
)

func main() {
	bot := models.Bot{
		ID:   0,
		Name: "Botto",
	}
	fmt.Println(bot.ID, bot.Name)

	server.ServeBot(&bot)
}
