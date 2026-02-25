package locations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type locationAreaDetail struct {
    PokemonEncounters []struct {
        Pokemon struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"pokemon"`
    } `json:"pokemon_encounters"`
}

// Explore fetches a location area by name or id and prints encountered
// Pokemon names.
func Explore(args []string) error {
    if len(args) == 0 {
        return errors.New("usage: explore <location-area-name-or-id>")
    }

    identifier := args[0]
    endpoint := fmt.Sprintf("%s/%s", baseURL, url.PathEscape(identifier))

    cacheHit := false
    if _, ok := requestCache.Get(endpoint); ok {
        cacheHit = true
    }

    body, err := fetchURL(endpoint)
    if err != nil {
        return err
    }

    var area locationAreaDetail
    decoder := json.NewDecoder(bytes.NewReader(body))
    if err := decoder.Decode(&area); err != nil {
        return err
    }

    fmt.Printf("Exploring %s...\n", identifier)
    count := len(area.PokemonEncounters)
    if count == 1 {
        fmt.Println("Found 1 Pokemon:")
    } else {
        fmt.Printf("Found %d Pokemon:\n", count)
    }
    for _, encounter := range area.PokemonEncounters {
        fmt.Printf("- %s\n", encounter.Pokemon.Name)
    }

    if cacheHit {
        fmt.Println("*** CACHE HIT ***")
    }
    fmt.Println("*** PAGE COMPLETE ***")

    return nil
}
