package commands

import "github.com/southwestmogrown/pokedexcli/maps"

// cliCommand represents a single command that can be executed from the
// REPL.  The repository keeps a registry of available commands so the
// REPL loop can look up and invoke them.
type cliCommand struct {
    Name        string
    Description string
    Callback    func() error
}

// GetCommands returns a map of command names to their metadata. The
// map is intentionally rebuilt on each call so that callers receive a
// fresh instance and may mutate it if necessary.
func GetCommands() map[string]cliCommand {
    return map[string]cliCommand{
        "exit": {
            Name:        "exit",
            Description: "Exit the Pokedex",
            Callback:    commandExit,
        },
        "help": {
            Name:        "help",
            Description: "Displays a help message",
            Callback:    commandHelp,
        },
        "map": {
            Name:        "map",
            Description: "Displays names of 20 location areas in the Pokemon world",
            Callback:    maps.Map,
        },
        "mapb": {
            Name:        "mapb",
            Description: "Displays location-area names by following the previous link",
            Callback:    maps.MapB,
        },
    }
}
