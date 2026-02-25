package commands

import "fmt"

func commandHelp() error {
    fmt.Println("Available commands:")
    for _, cmd := range GetCommands() {
        fmt.Printf(" - %s: %s\n", cmd.Name, cmd.Description)
    }
    return nil
}
