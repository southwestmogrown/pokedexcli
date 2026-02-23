package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func repl() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		word := words[0]
		if command, exists := commands[word]; exists {
			err := command.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}	
	}
}

func cleanInput(text string) []string {
	var cleanWords []string
	words := strings.Split(text, " ")
	
	for _, word := range words {
		if word == "" {
			continue
		}
		trimmed := strings.TrimSpace(word)
		lowerCase := strings.ToLower(trimmed)
		cleanWords = append(cleanWords, lowerCase)
	}

	return cleanWords
}

