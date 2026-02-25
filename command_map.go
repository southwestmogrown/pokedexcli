package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// baseURL is the initial endpoint used for the first invocation.
// tests can override this value to use a local server.
var baseURL = "https://pokeapi.co/api/v2/location-area"

// currentURL holds the URL that will be fetched when commandMap is
// executed. After each successful call it is updated to the "next"
// property returned by the API; when there are no more pages it is
// reset back to the empty string so that the next invocation will
// start over using baseURL.
var currentURL string

// currentLocations caches the most recent response so that callers
// (or tests) can inspect it if needed.
var currentLocations PokeLocations

// interactionManager keeps a simple history of every call to
// commandMap.  Each invocation stores the decoded PokeLocations
// payload so that callers can later examine what was fetched.
var interactionManager = InteractionManager{
	Interactions: make(map[int]PokeLocations),
}

func commandMap() error {
	// determine which URL to request
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

	// record this invocation for auditing or debugging
	interactionManager.InteractionCount++
	interactionManager.Interactions[interactionManager.InteractionCount] = currentLocations

	// display each location-area name on its own line
	for _, r := range currentLocations.Results {
		fmt.Println(r.Name)
	}

	// advance to the next page for the following invocation
	if currentLocations.Next != "" {
		currentURL = currentLocations.Next
	} else {
		// no more results; reset so a subsequent call hits the base URL
		currentURL = ""
	}

	return nil
}

// resetMapState clears any cached state so that commandMap will
// behave as if it has never been invoked. This is primarily useful
// for unit tests.
func resetMapState() {
	currentURL = ""
	currentLocations = PokeLocations{}
	resetInteractionManager()
}

// resetInteractionManager clears the history of commandMap calls.
func resetInteractionManager() {
	interactionManager = InteractionManager{Interactions: make(map[int]PokeLocations)}
}

type PokeLocations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}