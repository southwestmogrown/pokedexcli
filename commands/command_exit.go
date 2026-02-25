package commands

import (
	"fmt"
	"os"
)

func commandExit() error {
    fmt.Println("Bye!")
    os.Exit(0)
    return nil
}
