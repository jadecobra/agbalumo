# CLI Implementation Details

This package contains the command implementations for the Agbalumo CLI.

## Commands Registration
All commands should be registered in `RegisterCommands(rootCmd *cobra.Command)` in `register.go`. 

## Shared Helpers
Shared helpers like `InitRepo`, `ExitOnErr`, and `BindListingFlags` are located in `shared.go`.

## Adding New Commands
1. Create a new file in `cmd/cli/`.
2. Define the command variable (e.g., `NewCmd`).
3. Add subcommands in an `init()` function if necessary.
4. Export the command variable.
5. Register the command in `register.go`.
