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
        "catch": {
            Name:        "catch",
            Description: "Attempts to catch a Pokemon by name or ID (usage: catch <pokemon-name-or-id>)",
            Callback:    commandCatch,
        },
        "inspect": {
            Name:        "inspect",
            Description: "Shows details for a caught Pokemon (usage: inspect <pokemon-name>)",
            Callback:    commandInspect,
        },
        "pokedex": {
            Name:        "pokedex",
            Description: "Lists all caught Pokemon",
            Callback:    commandPokedex,
        },
    }
}
