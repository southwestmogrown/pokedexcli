package commands

import "github.com/southwestmogrown/pokedexcli/internal/locations"

func commandInspect(args []string) error {
    return locations.Inspect(args)
}
