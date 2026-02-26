package locations

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCatchRequiresArgument(t *testing.T) {
    err := Catch([]string{})
    if err == nil {
        t.Fatal("expected usage error when catch has no arguments")
    }
    if !strings.Contains(err.Error(), "usage: catch") {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestCatchEscapeThenCatch(t *testing.T) {
    originalRandIntn := randIntn
    defer func() { randIntn = originalRandIntn }()

    calls := 0
    randIntn = func(n int) int {
        calls++
        if calls == 1 {
            return 99 // force fail for first attempt
        }
        return 0 // force success for second attempt
    }

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/pokemon/pikachu/" {
            t.Fatalf("unexpected path: %s", r.URL.Path)
        }
        fmt.Fprintln(w, `{"id":25,"name":"pikachu","base_experience":112,"height":4,"weight":60,"stats":[],"types":[]}`)
    }))
    defer ts.Close()

    pokemonBaseURL = ts.URL + "/pokemon"
    ResetPokedex()
    resetRequestCache()

    out1 := captureOutput(func() {
        if err := Catch([]string{"pikachu"}); err != nil {
            t.Fatalf("first catch failed: %v", err)
        }
    })
    if !strings.Contains(out1, "Throwing a Pokeball at pikachu...") {
        t.Fatalf("missing throw log on first attempt: %s", out1)
    }
    if !strings.Contains(out1, "pikachu escaped!") {
        t.Fatalf("expected escape output on first attempt: %s", out1)
    }
    if _, exists := GetPokedex()["pikachu"]; exists {
        t.Fatalf("pikachu should not be in pokedex after escape")
    }

    out2 := captureOutput(func() {
        if err := Catch([]string{"pikachu"}); err != nil {
            t.Fatalf("second catch failed: %v", err)
        }
    })
    if !strings.Contains(out2, "Throwing a Pokeball at pikachu...") {
        t.Fatalf("missing throw log on second attempt: %s", out2)
    }
    if !strings.Contains(out2, "pikachu was caught!") {
        t.Fatalf("expected caught output on second attempt: %s", out2)
    }
    if !strings.Contains(out2, "You may now inspect it with the inspect command.") {
        t.Fatalf("missing inspect hint after successful catch: %s", out2)
    }

    pokedex := GetPokedex()
    pokemon, exists := pokedex["pikachu"]
    if !exists {
        t.Fatalf("expected pikachu in pokedex after catch")
    }
    if pokemon.ID != 25 || pokemon.BaseExperience != 112 {
        t.Fatalf("unexpected pokemon data stored: %+v", pokemon)
    }
}

func TestCatchByID(t *testing.T) {
    originalRandIntn := randIntn
    defer func() { randIntn = originalRandIntn }()
    randIntn = func(n int) int { return 0 }

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/pokemon/25/" {
            t.Fatalf("expected id path, got: %s", r.URL.Path)
        }
        fmt.Fprintln(w, `{"id":25,"name":"pikachu","base_experience":112,"height":4,"weight":60,"stats":[],"types":[]}`)
    }))
    defer ts.Close()

    pokemonBaseURL = ts.URL + "/pokemon"
    ResetPokedex()
    resetRequestCache()

    if err := Catch([]string{"25"}); err != nil {
        t.Fatalf("catch by id failed: %v", err)
    }
    if _, exists := GetPokedex()["pikachu"]; !exists {
        t.Fatalf("expected pikachu to be added to pokedex by id catch")
    }
}
