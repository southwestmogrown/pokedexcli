package locations

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestInspectRequiresArgument(t *testing.T) {
    err := Inspect([]string{})
    if err == nil {
        t.Fatal("expected usage error when inspect has no arguments")
    }
    if !strings.Contains(err.Error(), "usage: inspect") {
        t.Fatalf("unexpected error: %v", err)
    }
}

func TestInspectNotCaught(t *testing.T) {
    ResetPokedex()

    out := captureOutput(func() {
        if err := Inspect([]string{"pidgey"}); err != nil {
            t.Fatalf("inspect failed: %v", err)
        }
    })

    if !strings.Contains(out, "you have not caught that pokemon") {
        t.Fatalf("expected not-caught message, got: %s", out)
    }
}

func TestInspectCaughtPokemonDetails(t *testing.T) {
    ResetPokedex()

    var pidgey Pokemon
    payload := `{
      "id": 16,
      "name": "pidgey",
      "base_experience": 50,
      "height": 3,
      "weight": 18,
      "stats": [
        {"base_stat": 40, "effort": 0, "stat": {"name": "hp", "url": ""}},
        {"base_stat": 45, "effort": 0, "stat": {"name": "attack", "url": ""}},
        {"base_stat": 40, "effort": 0, "stat": {"name": "defense", "url": ""}},
        {"base_stat": 35, "effort": 0, "stat": {"name": "special-attack", "url": ""}},
        {"base_stat": 35, "effort": 0, "stat": {"name": "special-defense", "url": ""}},
        {"base_stat": 56, "effort": 1, "stat": {"name": "speed", "url": ""}}
      ],
      "types": [
        {"slot": 1, "type": {"name": "normal", "url": ""}},
        {"slot": 2, "type": {"name": "flying", "url": ""}}
      ]
    }`
    if err := json.Unmarshal([]byte(payload), &pidgey); err != nil {
        t.Fatalf("failed to unmarshal pokemon fixture: %v", err)
    }

    userPokedex[pidgey.Name] = pidgey

    out := captureOutput(func() {
        if err := Inspect([]string{"pidgey"}); err != nil {
            t.Fatalf("inspect failed: %v", err)
        }
    })

    expectedLines := []string{
        "Name: pidgey",
        "Height: 3",
        "Weight: 18",
        "Stats:",
        "  -hp: 40",
        "  -attack: 45",
        "  -defense: 40",
        "  -special-attack: 35",
        "  -special-defense: 35",
        "  -speed: 56",
        "Types:",
        "  - normal",
        "  - flying",
    }

    for _, line := range expectedLines {
        if !strings.Contains(out, line) {
            t.Fatalf("expected line %q in output, got:\n%s", line, out)
        }
    }
}
