package maps

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// baseURL is the initial endpoint used for the first invocation.
// tests can override this value to use a local server.
var baseURL = "https://pokeapi.co/api/v2/location-area"

var currentURL string
var currentLocations PokeLocations

// interactionManager keeps a simple history of every call to Map.
var interactionManager = InteractionManager{
    Interactions: make(map[int]PokeLocations),
}

// Map fetches from the PokeAPI and prints the names of each location-area.
// Subsequent calls follow the "next" link returned by the API, and when
// no further pages remain it resets to baseURL.
func Map() error {
    if currentURL == "" {
        currentURL = baseURL
    }

    res, err := http.Get(currentURL)
    if err != nil {
        fmt.Println(err)
        return err
    }
    defer res.Body.Close()

    decoder := json.NewDecoder(res.Body)
    if err := decoder.Decode(&currentLocations); err != nil {
        return err
    }

    interactionManager.InteractionCount++
    interactionManager.Interactions[interactionManager.InteractionCount] = currentLocations

    for _, r := range currentLocations.Results {
        fmt.Println(r.Name)
    }

    if currentLocations.Next != "" {
        currentURL = currentLocations.Next
    } else {
        currentURL = ""
    }

    return nil
}

// Reset clears pagination and interaction history.
func Reset() {
    currentURL = ""
    currentLocations = PokeLocations{}
    resetInteractionManager()
}

// InteractionManager records each Map invocation.
type InteractionManager struct {
    InteractionCount int
    Interactions     map[int]PokeLocations
}

func resetInteractionManager() {
    interactionManager = InteractionManager{Interactions: make(map[int]PokeLocations)}
}

// PokeLocations mirrors the structure returned by the API.
type PokeLocations struct {
    Count    int    `json:"count"`
    Next     string `json:"next"`
    Previous any    `json:"previous"`
    Results []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}
