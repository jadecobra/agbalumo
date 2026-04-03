# agbalumo CLI Documentation

The `agbalumo` CLI is a comprehensive tool for managing the directory platform, including listings, admin tasks, and the agent workflow harness.

## Index

- **[Listing Management](cli/listing.md)**
  - Create, update, and manage business listings.
- **[Admin Operations](cli/admin.md)**
  - Approve claims, manage users, and site configuration.
- **[Verification & Maintenance](cli/verify.md)**
  - Documentation drift, template checking, and coverage gates.
- **[Category Management](cli/category.md)**
  - Add and list custom categories.
- **[System Maintenance](cli/maintenance.md)**
  - Serve, Seed, Benchmark, Stress, and logs.

## Global Flags

The following flags are available for all commands:

| Flag | Description |
|------|-------------|
| `--json` | Output in JSON format for machine readability |
| `--verbose` | Enable verbose logging |

## Quick Start

```bash
# Start the web server
agbalumo serve

# Create a new listing
agbalumo listing create --title "Example Business" --city "Lagos"

# Approve a listing (admin)
agbalumo admin approve cli-12345
```
