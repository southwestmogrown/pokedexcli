package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/southwestmogrown/pokedexcli/commands"
)

// Repl starts the interactive prompt loop. It exits when the "exit"
// command calls os.Exit.
func Repl() {
    scanner := bufio.NewScanner(os.Stdin)
    commandsMap := commands.GetCommands()

    for {
        fmt.Print("Pokedex > ")
        scanner.Scan()
        words := cleanInput(scanner.Text())
        if len(words) == 0 {
            continue
        }
        word := words[0]
        if command, exists := commandsMap[word]; exists {
            err := command.Callback()
            if err != nil {
                fmt.Println(err)
            }
        } else {
            fmt.Println("Unknown command")
        }
    }
}

// cleanInput splits a line of text into lowercase words, discarding
// extra whitespace.
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
