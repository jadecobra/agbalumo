# Config Layer (internal/config)

Application configuration loaded from environment variables.

## Constraints
- Single file: `config.go` with `LoadConfig()` entry point.
- All env var keys should use `domain.EnvKey*` constants where available.
- `TraceMode` is toggled by `AGBALUMO_ENV=trace`.
