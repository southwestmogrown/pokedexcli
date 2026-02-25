package commands

import "fmt"

func commandHelp(args []string) error {
    fmt.Println("Available commands:")
    for _, cmd := range GetCommands() {
        fmt.Printf(" - %s: %s\n", cmd.Name, cmd.Description)
    }
    fmt.Println("Usage: explore <location-area-name-or-id> (e.g. explore canalave-city-area)")
    return nil
}
