package commands

import "github.com/southwestmogrown/pokedexcli/internal/locations"

func commandMap(args []string) error {
    return locations.Map()
}
