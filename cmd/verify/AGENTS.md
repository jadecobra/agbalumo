# Verification Tooling (cmd/verify)

This package contains the CLI entry points for the `verify` tool.

## Constraints
- **Trusted Commands**: Commands here are used for maintenance and run in trusted environments. `exec.Command` should be used with `#nosec G204` if arguments are somewhat dynamic but trusted.
- **Dependency Hygiene**: This package should primarily coordinate calls into `internal/maintenance/`. Avoid implementing heavy logic here.
- **Reporting**: Favor structured output (tables/JSON) over verbose prose.
