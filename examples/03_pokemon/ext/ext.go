package main

import (
	"fmt"
	"log"

	"github.com/asmcos/requests"
	"github.com/jaimeteb/chatto/ext"
	"github.com/jaimeteb/chatto/fsm"
)

func searchPokemon(req *ext.Request) (res *ext.Response) {
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
		intoState = req.Dom.StateTable["search_pokemon"]
	} else {
		if response.R.StatusCode == 404 {
			message = "Pokémon not found, try with another input."
			intoState = req.Dom.StateTable["search_pokemon"]
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

	return &ext.Response{
		FSM: &fsm.FSM{
			State: intoState,
			Slots: req.FSM.Slots,
		},
		Res: message,
	}
}

var myExtMap = ext.ExtensionMap{
	"ext_search_pokemon": searchPokemon,
}

func main() {
	if err := ext.ServeExtensionRPC(myExtMap); err != nil {
		log.Fatalln(err)
	}
}
