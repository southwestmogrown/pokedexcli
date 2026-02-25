package locations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	cachepkg "github.com/southwestmogrown/pokedexcli/internal/cache"
)

// baseURL is the initial endpoint; tests may override if necessary.
var baseURL = "https://pokeapi.co/api/v2/location-area"

var currentURL string
var currentLocations PokeLocations

// lastFetchedURL holds the URL of the page most recently retrieved by
// either Map or MapB. It's used to avoid redundant fetches when one
// command leaves state for the other.
var lastFetchedURL string

// interactionManager keeps a simple history of every call to Map.
var interactionManager = InteractionManager{
    Interactions: make(map[int]PokeLocations),
}

// requestCache stores raw HTTP responses keyed by URL.
var requestCache = cachepkg.NewPokeCache(1 * time.Minute)

// resetRequestCache clears the global cache (intended for tests).
func resetRequestCache() {
    requestCache = cachepkg.NewPokeCache(1 * time.Minute)
}

// Map fetches from the PokeAPI and prints the names of each location-area.
// Subsequent calls follow the "next" link returned by the API, and when
// no further pages remain it resets to baseURL.
func Map() error {
    // remember whether we've wrapped from the end back to start; we'll
    // report that after printing names rather than up front.
    endWrap := currentURL == "" && lastFetchedURL != ""
    if currentURL == "" {
        currentURL = baseURL
    }

    // if the URL hasn't changed since the last fetch we already have the
    // contents in currentLocations. This situation happens when MapB was
    // called immediately prior, so instead of refetching the same page we
    // bump to the next link now.
    if currentURL != "" && currentURL == lastFetchedURL {
        if currentLocations.Next != "" {
            currentURL = currentLocations.Next
        } else {
            currentURL = ""
            return nil
        }
    }

    // determine if we'll hit the cache before fetching.
    cacheHit := false
    if _, ok := requestCache.Get(currentURL); ok {
        cacheHit = true
    }

    body, err := fetchURL(currentURL)
    if err != nil {
        fmt.Println(err)
        return err
    }
    decoder := json.NewDecoder(bytes.NewReader(body))
    if err := decoder.Decode(&currentLocations); err != nil {
        return err
    }
    lastFetchedURL = currentURL

    // determine whether this page had no next link
    endReached := currentLocations.Next == ""

    interactionManager.InteractionCount++
    interactionManager.Interactions[interactionManager.InteractionCount] = currentLocations

    fmt.Println("Exploring location areas...")
    count := len(currentLocations.Results)
    if count == 1 {
        fmt.Println("Found 1 Location Area:")
    } else {
        fmt.Printf("Found %d Location Areas:\n", count)
    }
    for _, r := range currentLocations.Results {
        fmt.Printf("- %s\n", r.Name)
    }

    // report any special conditions at the bottom
    if cacheHit {
        fmt.Println("*** CACHE HIT ***")
    }
    if endReached {
        fmt.Println("*** REACHED END ***")
    }
    if endWrap {
        fmt.Println("*** WRAPPED TO START ***")
    }

    fmt.Println("*** PAGE COMPLETE ***")

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
    lastFetchedURL = ""
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

// fetchURL retrieves data, using the cache when possible.
func fetchURL(url string) ([]byte, error) {
    if data, ok := requestCache.Get(url); ok {
        return data, nil
    }
    res, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    requestCache.Add(url, body)
    return body, nil
}

// PokeLocations mirrors the structure returned by the API.
type PokeLocations struct {
    Count    int    `json:"count"`
    Next     string `json:"next"`
    Previous any    `json:"previous"`
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}
