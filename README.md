# agbalumo

A robust web application platform for the West African diaspora community, featuring a business directory, job board, event listings, and community requests.

## Quick Start

```bash
# Build and start the server
./scripts/verify_restart.sh

# Access the application
open https://localhost:8443
```

## Documentation

- **[CLI Commands](docs/cli.md)** - Command-line interface for managing listings and admin operations
- **[HTTP API Reference](docs/api.md)** - REST API endpoints and authentication
- **[OpenAPI Specification](docs/openapi.yaml)** - Full API schema (OpenAPI 3.0.3)

## Environment & Setup

### Development Server
The development server runs securely on **HTTPS Port 8443**: `https://localhost:8443`
- The project is configured to use self-signed certificates in `certs/`.
- **Do not** attempt to access `http://localhost:8080` in development if certificates are present; the server will auto-switch to 8443.

### Go Environment (macOS)
- **Go Binary**: Managed via Homebrew.
- **Path**: Ensure `/opt/homebrew/bin` is in your `PATH`.
- **Do not** look for Go in `/usr/local/go`.

### Scripts
Reference the `scripts/` directory for standard operations:
- `scripts/verify_restart.sh`: Rebuilds and restarts the server safely, handling process cleanup and environment variables.
- `scripts/pre-commit.sh`: Runs tests, coverage checks, and linting.

### Database
- Uses SQLite with `.tester/data/agbalumo.db`.
- Database URL defaults to `.tester/data/agbalumo.db` in `env`.

## Running the Project
1. Ensure `.env` is set up.
2. Run `scripts/verify_restart.sh` to build and start the server.
3. Access `https://localhost:8443`.
