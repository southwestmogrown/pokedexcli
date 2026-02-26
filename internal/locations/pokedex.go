package locations

import "fmt"

// Pokedex prints all caught pokemon names in capture order.
func Pokedex(args []string) error {
    fmt.Println("Your Pokedex:")
    if len(pokedexOrder) == 0 {
        fmt.Println("  (empty)")
        return nil
    }
    for _, name := range pokedexOrder {
        fmt.Printf(" - %s\n", name)
    }
    return nil
}
