package main

import (
	"fmt"
	"log"

	"github.com/asmcos/requests"
	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

func searchPokemon(req *extension.Request) (res *extension.Response) {
	m := req.FSM

	pokemon := m.Slots["pokemon"]

	var message string
	var intoState int

	intoState = req.FSM.State

	r := requests.Requests()
	r.Header.Set("Content-Type", "application/json")
	response, err := r.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", pokemon))

	if err != nil {
		message = "Something went wrong..."
		intoState = req.DB.StateTable["search_pokemon"]
	} else {
		if response.R.StatusCode == 404 {
			message = "Pokémon not found, try with another input."
			intoState = req.DB.StateTable["search_pokemon"]
		} else {
			var json map[string]interface{}
			response.Json(&json)
			pokemonName := json["name"].(string)
			pokemonID := json["id"].(float64)
			pokemonHeight := json["height"].(float64)
			pokemonWeight := json["weight"].(float64)
			message = fmt.Sprintf("Name: %v \nID: %v \nHeight: %v \nWeight: %v", pokemonName, pokemonID, pokemonHeight, pokemonWeight)
		}
	}

	return &extension.Response{
		FSM: &fsm.FSM{
			State: intoState,
			Slots: req.FSM.Slots,
		},
		Answers: []query.Answer{{Text: message}},
	}
}

var myExtMap = extension.RegisteredFuncs{
	"ext_search_pokemon": searchPokemon,
}

func main() {
	if err := extension.ServeRPC(myExtMap); err != nil {
		log.Fatalln(err)
	}
}
