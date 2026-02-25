package maps

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// baseURLB parallels baseURL in map.go and may be overridden in tests.
var baseURLB = "https://pokeapi.co/api/v2/location-area"

// the package-level requestCache from map.go is reused for MapB;
// resetRequestCache defined there is also used by MapB tests.

// nothing to define here

// MapB behaves similarly to Map but follows the "previous" link in the
// response instead of "next". When there is no previous value the
// state is reset so a subsequent call starts at baseURLB again.
func MapB() error {
    // choose the URL that corresponds to the previous page of whatever we
    // last loaded.  If we've never loaded anything yet (empty
    // currentLocations) or the previous field is empty, start at baseURLB.
    var url string
    beginReached := false
    if prev, ok := currentLocations.Previous.(string); ok && prev != "" {
        url = prev
    } else {
        beginReached = true
        url = baseURLB
    }

    // pre-check cache similar to Map
    cacheHit := false
    if _, ok := requestCache.Get(url); ok {
        cacheHit = true
    }

    body, err := fetchURL(url)
    if err != nil {
        fmt.Println(err)
        return err
    }
    var locs PokeLocations
    decoder := json.NewDecoder(bytes.NewReader(body))
    if err := decoder.Decode(&locs); err != nil {
        return err
    }

    // update shared state so subsequent calls (map or mapb) see it
    currentLocations = locs
    // remember where we fetched from; tests rely on this value
    currentURL = url
    lastFetchedURL = url

    for _, r := range currentLocations.Results {
        fmt.Println(r.Name)
    }

    if cacheHit {
        fmt.Println("*** CACHE HIT ***")
    }
    if beginReached {
        fmt.Println("*** REACHED BEGINNING ***")
    }

    // matching flair for backwards pagination
    fmt.Println("*** PAGE COMPLETE ***")

    return nil
}

// ResetB is a thin wrapper around Reset so tests can continue to
// call the old function name if desired. It resets shared map state.
func ResetB() {
    Reset()
}

// fetchURL is implemented in map.go and available here via package
// scope; no duplicate definition needed.
