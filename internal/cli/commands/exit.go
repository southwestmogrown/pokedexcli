package commands

import (
	"fmt"
	"os"
)

var exitFn = os.Exit

func commandExit(args []string) error {
    fmt.Println("Bye!")
    exitFn(0)
    return nil
}
