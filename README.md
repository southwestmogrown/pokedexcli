# Pokedex CLI

A terminal-based Pokedex application written in Go.

It provides an interactive REPL for exploring Pokemon location areas using the public PokeAPI location-area endpoint, with in-memory caching and pagination support.

## Features

- Interactive CLI prompt (`Pokedex >`)
- Command registry with argument-aware callbacks
- Forward pagination through location areas (`map`)
- Backward pagination (`mapb`)
- Explore a location area by name or id to list encountered Pokemon (`explore <name-or-id>`)
- In-memory cache with expiration and background reaper
- User-friendly footer messages for cache hits and pagination boundaries

## Architecture

The project follows a conventional Go CLI layout:

- `cmd/pokedex`  
  Executable entrypoint.
- `internal/cli`  
  REPL loop and input normalization.
- `internal/cli/commands`  
  Command registry + one file per command handler.
- `internal/locations`  
  Domain logic for location-area pagination and exploration.
- `internal/cache`  
  Thread-safe TTL cache used by location requests.

### High-level flow

1. `main` starts the REPL.
2. REPL parses input into command + args.
3. Command handler delegates to location domain logic.
4. Location logic fetches data from cache or API.
5. Results and status messages are printed to the CLI.

## Setup

### Prerequisites

- Go `1.25+` (matches `go.mod`)

### Install / clone

```bash
git clone <your-repo-url>
cd pokedex
```

## Run

### Start the CLI

```bash
go run ./cmd/pokedex
```

You should see:

```text
Pokedex >
```

## Usage

Type `help` in the REPL to see all commands.

### Commands

- `help`  
  Shows available commands and explore usage.
- `exit`  
  Exits the CLI.
- `map`  
  Shows location-area names, moving forward one page per call.
- `mapb`  
  Shows location-area names, moving backward one page per call.
- `explore <location-area-name-or-id>`  
  Lists Pokemon names encountered in that location area.

### Example session

```text
Pokedex > help
Pokedex > map
Pokedex > map
Pokedex > mapb
Pokedex > explore canalave-city-area
Pokedex > explore 1
Pokedex > exit
```

## Output conventions

Location commands print result names first, then optional status footers:

- `*** CACHE HIT ***` when response data is served from cache
- `*** REACHED END ***` when `map` reaches the last page
- `*** WRAPPED TO START ***` when `map` starts again from the beginning
- `*** REACHED BEGINNING ***` when `mapb` cannot go further back
- `*** PAGE COMPLETE ***` after each command completes

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
- Keep API and pagination behavior in domain packages, not in REPL code.
- Reuse the shared cache layer for network-backed commands.
- Preserve user-facing output format unless intentionally changing UX.
