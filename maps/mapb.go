package maps

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var baseURLB = "https://pokeapi.co/api/v2/location-area"
var currentURLB string
var currentLocationsB PokeLocations

// MapB behaves similarly to Map but follows the "previous" link in the
// response instead of "next". When there is no previous value the
// state is reset so a subsequent call starts at baseURLB again.
func MapB() error {
    if currentURLB == "" {
        currentURLB = baseURLB
    }

    res, err := http.Get(currentURLB)
    if err != nil {
        fmt.Println(err)
        return err
    }
    defer res.Body.Close()

    decoder := json.NewDecoder(res.Body)
    if err := decoder.Decode(&currentLocationsB); err != nil {
        return err
    }

    for _, r := range currentLocationsB.Results {
        fmt.Println(r.Name)
    }

    if prev, ok := currentLocationsB.Previous.(string); ok && prev != "" {
        currentURLB = prev
    } else {
        currentURLB = ""
    }

    return nil
}

// ResetB clears the internal state used by MapB.
func ResetB() {
    currentURLB = ""
    currentLocationsB = PokeLocations{}
}
