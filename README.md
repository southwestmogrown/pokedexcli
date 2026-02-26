# Pokedex CLI

A terminal-based Pokedex application written in Go.

The app runs an interactive REPL (`Pokedex >`) and supports browsing location areas, exploring encounters, catching Pokemon, inspecting caught Pokemon, and listing your Pokédex.

## Features

- Interactive CLI prompt (`Pokedex >`)
- Forward and backward location pagination (`map`, `mapb`)
- Location-area exploration by name/id (`explore`)
- Catch mechanics using Pokemon `base_experience` + `math/rand` (`catch`)
- Local in-memory Pokédex for caught Pokemon (`inspect`, `pokedex`)
- Shared HTTP cache with TTL + background reaper

## Architecture

The project follows a conventional Go CLI layout:

- `cmd/pokedex`  
  Executable entrypoint.
- `internal/cli`  
  REPL loop and input normalization.
- `internal/cli/commands`  
  Command registry and one file per command handler.
- `internal/locations`  
  Domain logic for map/mapb/explore/catch/inspect/pokedex.
- `internal/cache`  
  Thread-safe TTL cache used by network-backed location logic.

### High-level flow

1. `main` starts the REPL.
2. REPL parses input into command + args.
3. Command handler delegates to `internal/locations`.
4. Domain logic reads cached data or calls PokeAPI when needed.
5. User-facing output is printed with consistent formatting.

## Setup

### Prerequisites

- Go `1.25+` (matches `go.mod`)

### Install / clone

```bash
git clone <your-repo-url>
cd pokedex
```

## Run

```bash
go run ./cmd/pokedex
```

You should see:

```text
Pokedex >
```

## Commands

- `help`  
  Shows available commands and usage examples.
- `exit`  
  Exits the CLI.
- `map`  
  Shows location-area names, moving forward one page per call.
- `mapb`  
  Shows location-area names, moving backward one page per call.
- `explore <location-area-name-or-id>`  
  Lists Pokemon encountered in that location area.
- `catch <pokemon-name-or-id>`  
  Attempts to catch a Pokemon using base-experience catch odds.
- `inspect <pokemon-name>`  
  Prints details for a caught Pokemon (no API call; reads local Pokédex).
- `pokedex`  
  Lists all caught Pokemon in capture order.

## Example session

```text
Pokedex > catch pidgey
Throwing a Pokeball at pidgey...
pidgey was caught!
You may now inspect it with the inspect command.

Pokedex > inspect pidgey
Name: pidgey
Height: 3
Weight: 18
Stats:
  -hp: 40
  -attack: 45
  -defense: 40
  -special-attack: 35
  -special-defense: 35
  -speed: 56
Types:
  - normal
  - flying

Pokedex > pokedex
Your Pokedex:
 - pidgey
```

## Output conventions

Location and exploration commands print itemized lists first, followed by optional status footers:

- `*** CACHE HIT ***` when response data is served from cache
- `*** REACHED END ***` when `map` reaches the last page
- `*** WRAPPED TO START ***` when `map` wraps back to the first page
- `*** REACHED BEGINNING ***` when `mapb` cannot go further back
- `*** PAGE COMPLETE ***` when command output completes

`pokedex` prints:

- `Your Pokedex:` and then each caught pokemon as ` - <name>`
- `  (empty)` when no Pokemon have been caught yet

## Testing

Run all tests:

```bash
go test ./...
```

Run targeted packages:

```bash
go test ./internal/cli/...
go test ./internal/locations/...
go test ./internal/cache/...
```

## Notes for contributors

- Keep command handlers thin: parse/validate args, then delegate to `internal/locations`.
- Keep business logic in domain packages, not the REPL loop.
- Reuse shared cache utilities for network-backed behavior.
- Preserve user-facing output format unless intentionally changing UX and tests together.
