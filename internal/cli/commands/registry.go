package commands

type Command struct {
    Name        string
    Description string
    Callback    func(args []string) error
}

func GetCommands() map[string]Command {
    return map[string]Command{
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
            Callback:    commandMap,
        },
        "mapb": {
            Name:        "mapb",
            Description: "Displays location-area names by following the previous link",
            Callback:    commandMapB,
        },
        "explore": {
            Name:        "explore",
            Description: "Displays Pokemon encountered in a location area by name or ID (usage: explore <location-area-name-or-id>)",
            Callback:    commandExplore,
        },
    }
}
