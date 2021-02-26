package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

	response, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon))
	if err != nil {
		message = "Something went wrong..."
		intoState = req.Domain.StateTable["search_pokemon"]

		return &extension.Response{
			FSM: &fsm.FSM{
				State: intoState,
				Slots: req.FSM.Slots,
			},
			Answers: []query.Answer{{Text: message}},
		}
	}

	if response.StatusCode == 404 {
		message = "Pokémon not found, try with another input."
		intoState = req.Domain.StateTable["search_pokemon"]

		return &extension.Response{
			FSM: &fsm.FSM{
				State: intoState,
				Slots: req.FSM.Slots,
			},
			Answers: []query.Answer{{Text: message}},
		}
	}

	var pokemonResp map[string]interface{}

	body, readAllErr := ioutil.ReadAll(response.Body)
	if readAllErr != nil {
		message = "Pokémon not found, try with another input."
		intoState = req.Domain.StateTable["search_pokemon"]

		return &extension.Response{
			FSM: &fsm.FSM{
				State: intoState,
				Slots: req.FSM.Slots,
			},
			Answers: []query.Answer{{Text: message}},
		}
	}

	unmarshalErr := json.Unmarshal(body, &pokemonResp)
	if unmarshalErr != nil {
		message = "Pokémon not found, try with another input."
		intoState = req.Domain.StateTable["search_pokemon"]

		return &extension.Response{
			FSM: &fsm.FSM{
				State: intoState,
				Slots: req.FSM.Slots,
			},
			Answers: []query.Answer{{Text: message}},
		}
	}

	pokemonName := pokemonResp["name"].(string)
	pokemonID := pokemonResp["id"].(float64)
	pokemonHeight := pokemonResp["height"].(float64)
	pokemonWeight := pokemonResp["weight"].(float64)
	message = fmt.Sprintf("Name: %s \nID: %.2f \nHeight: %.2f \nWeight: %.2f", pokemonName, pokemonID, pokemonHeight, pokemonWeight)

	return &extension.Response{
		FSM: &fsm.FSM{
			State: intoState,
			Slots: req.FSM.Slots,
		},
		Answers: []query.Answer{{Text: message}},
	}
}

var RegisteredCommandFuncs = extension.RegisteredCommandFuncs{
	"search_pokemon": searchPokemon,
}

func main() {
	if err := extension.ServeRPC(RegisteredCommandFuncs); err != nil {
		log.Fatalln(err)
	}
}
