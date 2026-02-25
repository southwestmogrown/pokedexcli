package commands

import "github.com/southwestmogrown/pokedexcli/internal/locations"

func commandExplore(args []string) error {
    return locations.Explore(args)
}
