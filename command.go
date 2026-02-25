package main

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the Pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "Displays a help message",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "Displays names of 20 location areas in the Pokemon world",
			callback: commandMap,
		},
	}
}


type cliCommand struct {
	name string
	description string
	callback func() error
}