package commands

import "fmt"

func commandHelp(args []string) error {
    fmt.Println("Available commands:")
    for _, cmd := range GetCommands() {
        fmt.Printf(" - %s: %s\n", cmd.Name, cmd.Description)
    }
    fmt.Println("Usage: explore <location-area-name-or-id> (e.g. explore canalave-city-area)")
    fmt.Println("Usage: catch <pokemon-name-or-id> (e.g. catch pidgey)")
    fmt.Println("Usage: inspect <pokemon-name> (e.g. inspect pidgey)")
    return nil
}
