package commands

import "github.com/southwestmogrown/pokedexcli/internal/locations"

func commandCatch(args []string) error {
    return locations.Catch(args)
}
