package locations

import (
	"strings"
	"testing"
)

func TestPokedexOutputEmpty(t *testing.T) {
    ResetPokedex()

    out := captureOutput(func() {
        if err := Pokedex([]string{}); err != nil {
            t.Fatalf("pokedex command failed: %v", err)
        }
    })

    if !strings.Contains(out, "Your Pokedex:") {
        t.Fatalf("missing pokedex header: %s", out)
    }
    if !strings.Contains(out, "(empty)") {
        t.Fatalf("expected empty-state message, got: %s", out)
    }
}

func TestPokedexOutput(t *testing.T) {
    ResetPokedex()

    userPokedex["pidgey"] = Pokemon{Name: "pidgey"}
    userPokedex["caterpie"] = Pokemon{Name: "caterpie"}
    pokedexOrder = []string{"pidgey", "caterpie"}

    out := captureOutput(func() {
        if err := Pokedex([]string{}); err != nil {
            t.Fatalf("pokedex command failed: %v", err)
        }
    })

    if !strings.Contains(out, "Your Pokedex:") {
        t.Fatalf("missing pokedex header: %s", out)
    }

    expected := []string{" - pidgey", " - caterpie"}
    for _, line := range expected {
        if !strings.Contains(out, line) {
            t.Fatalf("expected line %q in output, got:\n%s", line, out)
        }
    }
}
