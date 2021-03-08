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

func searchPokemon(req *extension.ExecuteExtensionRequest) (res *extension.ExecuteExtensionResponse) {
	m := req.FSM

	pokemon := m.Slots["pokemon"]

	var message string
	var intoState int

	intoState = req.FSM.State

	response, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon))
	if err != nil {
		message = "Something went wrong..."
		intoState = req.Domain.StateTable["search_pokemon"]

		return &extension.ExecuteExtensionResponse{
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

		return &extension.ExecuteExtensionResponse{
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

		return &extension.ExecuteExtensionResponse{
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

		return &extension.ExecuteExtensionResponse{
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

	return &extension.ExecuteExtensionResponse{
		FSM: &fsm.FSM{
			State: intoState,
			Slots: req.FSM.Slots,
		},
		Answers: []query.Answer{{Text: message}},
	}
}

var registeredExtensions = extension.RegisteredExtensions{
	"search_pokemon": searchPokemon,
}

func main() {
	if err := extension.ServeRPC(registeredExtensions); err != nil {
		log.Fatalln(err)
	}
}
