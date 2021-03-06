# Pokemon

The [**Pokémon**](https://github.com/jaimeteb/chatto/tree/master/examples/03_pokemon) bot is a simple integration with the [PokéAPI](https://pokeapi.co/). The purpose of this example is to demonstrate the [RPC Extension Server](/extensions#go).

## Diagram

This bot's Finite State Machine can be visualized like this:

![Pokemon](/img/chatto_pokemon.svg)

This bot also demonstrates:

* The state [**any**](/finitestatemachine/#any), that is used in this case to transition from any state into the initial state, when the **faq** command is executed.
* The command [**any**](/finitestatemachine/#any), that is used to always run the **search_pokemon** extension, when in the **search_pokemon** state.
* The [slot](/finitestatemachine/#slots) **pokemon** is saved when transitioning from **search_pokemon** into the intial state.

## Run it

To run this example:

```bash
go run examples/03_pokemon/ext/ext.go
```

And in other terminal:

```bash
chatto -path examples/03_pokemon/
```
