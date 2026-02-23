package main

import (
	"fmt"
)

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	commands := getCommands()

	for _, command := range commands {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}

	return nil
}