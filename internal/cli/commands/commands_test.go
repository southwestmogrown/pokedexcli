package commands

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func captureCommandOutput(f func()) string {
    var buf bytes.Buffer
    originalStdout := os.Stdout
    readPipe, writePipe, _ := os.Pipe()
    os.Stdout = writePipe

    f()

    writePipe.Close()
    _, _ = io.Copy(&buf, readPipe)
    os.Stdout = originalStdout
    return buf.String()
}

func TestGetCommandsIncludesAllCommands(t *testing.T) {
    commandMap := GetCommands()

    expected := []string{"help", "exit", "map", "mapb", "explore", "catch", "inspect", "pokedex"}
    for _, name := range expected {
        cmd, ok := commandMap[name]
        if !ok {
            t.Fatalf("missing command %q", name)
        }
        if cmd.Callback == nil {
            t.Fatalf("command %q has nil callback", name)
        }

        callbackType := reflect.TypeOf(cmd.Callback)
        expectedType := reflect.TypeOf(func(args []string) error { return nil })
        if callbackType != expectedType {
            t.Fatalf("command %q has wrong callback signature: got %v", name, callbackType)
        }
    }
}

func TestHelpCommand(t *testing.T) {
    commandMap := GetCommands()
    out := captureCommandOutput(func() {
        if err := commandMap["help"].Callback([]string{"ignored"}); err != nil {
            t.Fatalf("help command failed: %v", err)
        }
    })

    if !strings.Contains(out, "Available commands:") {
        t.Fatalf("help output missing header: %s", out)
    }
    for _, name := range []string{"help", "exit", "map", "mapb", "explore", "catch", "inspect", "pokedex"} {
        if !strings.Contains(out, name) {
            t.Fatalf("help output missing %q: %s", name, out)
        }
    }
    if !strings.Contains(out, "explore <location-area-name-or-id>") {
        t.Fatalf("help output missing explore usage example: %s", out)
    }
    if !strings.Contains(out, "catch <pokemon-name-or-id>") {
        t.Fatalf("help output missing catch usage example: %s", out)
    }
    if !strings.Contains(out, "inspect <pokemon-name>") {
        t.Fatalf("help output missing inspect usage example: %s", out)
    }
}

func TestCatchCommandRequiresArgument(t *testing.T) {
    commandMap := GetCommands()

    err := commandMap["catch"].Callback([]string{})
    if err == nil {
        t.Fatal("expected catch command to require one argument")
    }
    if !strings.Contains(err.Error(), "usage: catch") {
        t.Fatalf("unexpected catch error: %v", err)
    }
}

func TestInspectCommandRequiresArgument(t *testing.T) {
    commandMap := GetCommands()

    err := commandMap["inspect"].Callback([]string{})
    if err == nil {
        t.Fatal("expected inspect command to require one argument")
    }
    if !strings.Contains(err.Error(), "usage: inspect") {
        t.Fatalf("unexpected inspect error: %v", err)
    }
}

func TestExitCommand(t *testing.T) {
    originalExitFn := exitFn
    defer func() { exitFn = originalExitFn }()

    exitCode := -1
    exitFn = func(code int) {
        exitCode = code
    }

    out := captureCommandOutput(func() {
        if err := commandExit([]string{"ignored"}); err != nil {
            t.Fatalf("exit command failed: %v", err)
        }
    })

    if exitCode != 0 {
        t.Fatalf("expected exit code 0, got %d", exitCode)
    }
    if !strings.Contains(out, "Bye!") {
        t.Fatalf("exit output missing goodbye: %s", out)
    }
}

func TestExploreCommandRequiresArgument(t *testing.T) {
    commandMap := GetCommands()

    err := commandMap["explore"].Callback([]string{})
    if err == nil {
        t.Fatal("expected explore command to require one argument")
    }
    if !strings.Contains(err.Error(), "usage: explore") {
        t.Fatalf("unexpected explore error: %v", err)
    }
}
