package locations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
)

var pokemonBaseURL = "https://pokeapi.co/api/v2/pokemon"

// randIntn is a package-level indirection so tests can make catch behavior deterministic.
var randIntn = rand.Intn

// userPokedex stores captured pokemon keyed by pokemon name.
var userPokedex = map[string]Pokemon{}

// pokedexOrder keeps deterministic insertion order for listing in the pokedex command.
var pokedexOrder []string

// Pokemon contains data returned from the pokemon endpoint.
type Pokemon struct {
    ID             int    `json:"id"`
    Name           string `json:"name"`
    BaseExperience int    `json:"base_experience"`
    Height         int    `json:"height"`
    Weight         int    `json:"weight"`
    Stats          []struct {
        BaseStat int `json:"base_stat"`
        Effort   int `json:"effort"`
        Stat     struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"stat"`
    } `json:"stats"`
    Types []struct {
        Slot int `json:"slot"`
        Type struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"type"`
    } `json:"types"`
}

// Catch attempts to catch a pokemon by name or id.
func Catch(args []string) error {
    if len(args) == 0 {
        return errors.New("usage: catch <pokemon-name-or-id>")
    }

    identifier := args[0]
    fmt.Printf("Throwing a Pokeball at %s...\n", identifier)
    endpoint := fmt.Sprintf("%s/%s/", pokemonBaseURL, url.PathEscape(identifier))

    body, err := fetchURL(endpoint)
    if err != nil {
        return err
    }

    var pokemon Pokemon
    decoder := json.NewDecoder(bytes.NewReader(body))
    if err := decoder.Decode(&pokemon); err != nil {
        return err
    }

    caught := canCatchPokemon(pokemon.BaseExperience)

    if !caught {
        fmt.Printf("%s escaped!\n", pokemon.Name)
        return nil
    }

    if _, exists := userPokedex[pokemon.Name]; !exists {
        pokedexOrder = append(pokedexOrder, pokemon.Name)
    }
    userPokedex[pokemon.Name] = pokemon
    fmt.Printf("%s was caught!\n", pokemon.Name)
    fmt.Println("You may now inspect it with the inspect command.")
    return nil
}

func canCatchPokemon(baseExperience int) bool {
    // keep encounters catchable in a reasonable number of tries while still
    // rewarding weaker pokemon with higher catch rates.
    chance := 75 - (baseExperience / 6)
    if chance < 15 {
        chance = 15
    }
    if chance > 90 {
        chance = 90
    }
    roll := randIntn(100)
    return roll < chance
}

// ResetPokedex clears captured pokemon state (primarily for tests).
func ResetPokedex() {
    userPokedex = map[string]Pokemon{}
    pokedexOrder = nil
}

// GetPokedex returns a copy of the current captured pokemon map.
func GetPokedex() map[string]Pokemon {
    cp := make(map[string]Pokemon, len(userPokedex))
    for name, pokemon := range userPokedex {
        cp[name] = pokemon
    }
    return cp
}
