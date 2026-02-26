package locations

import (
	"errors"
	"fmt"
)

// Inspect prints details for a previously-caught pokemon.
// It only reads from the in-memory pokedex and makes no API calls.
func Inspect(args []string) error {
    if len(args) == 0 {
        return errors.New("usage: inspect <pokemon-name>")
    }

    name := args[0]
    pokemon, ok := userPokedex[name]
    if !ok {
        fmt.Println("you have not caught that pokemon")
        return nil
    }

    fmt.Printf("Name: %s\n", pokemon.Name)
    fmt.Printf("Height: %d\n", pokemon.Height)
    fmt.Printf("Weight: %d\n", pokemon.Weight)
    fmt.Println("Stats:")
    for _, stat := range pokemon.Stats {
        fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
    }
    fmt.Println("Types:")
    for _, typeInfo := range pokemon.Types {
        fmt.Printf("  - %s\n", typeInfo.Type.Name)
    }

    return nil
}
