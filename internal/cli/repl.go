package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/southwestmogrown/pokedexcli/internal/cli/commands"
)

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

        name := words[0]
        args := words[1:]
        if command, exists := commandsMap[name]; exists {
            if err := command.Callback(args); err != nil {
                fmt.Println(err)
            }
            continue
        }
        fmt.Println("Unknown command")
    }
}
