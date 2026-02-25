package locations

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// baseURLB parallels baseURL in map.go and may be overridden in tests.
var baseURLB = "https://pokeapi.co/api/v2/location-area"

// MapB behaves similarly to Map but follows the "previous" link in the
// response instead of "next". When there is no previous value the
// state is reset so a subsequent call starts at baseURLB again.
func MapB() error {
    var url string
    beginReached := false
    if prev, ok := currentLocations.Previous.(string); ok && prev != "" {
        url = prev
    } else {
        beginReached = true
        url = baseURLB
    }

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

    currentLocations = locs
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

    fmt.Println("*** PAGE COMPLETE ***")

    return nil
}

// ResetB is a thin wrapper around Reset.
func ResetB() {
    Reset()
}
