package commands

import "github.com/southwestmogrown/pokedexcli/internal/locations"

func commandPokedex(args []string) error {
    return locations.Pokedex(args)
}
