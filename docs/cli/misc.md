# agbalumo CLI: Miscellaneous Commands

Utility and data management commands for the agbalumo platform.

## Commands

### seed

Seed the database with initial data.

```bash
agbalumo seed [database-path]
```

### backfill-cities

Backfill missing cities from addresses in listings.
Uses standard data sources to match cities from partial addresses.

```bash
agbalumo backfill-cities
```

### category

Manage listing categories.

#### Subcommands

##### add
Add a new category.
```bash
agbalumo category add [name] [flags]
```

##### list
List all categories.
```bash
agbalumo category list
```

### aglog

Capture squad decisions for the Learning Loop.

```bash
aglog [flags]
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--feature` | Feature name (required) |
| `--arch` | Systems Architect name |
| `--po` | Product Owner name |
| `--sdet` | SDET name |
| `--be` | Backend Engineer name |
| `--summary` | Decision summary |

**Example:**
```bash
aglog --feature "aglog-cli" --arch "Gemini" --summary "Refactored to cobra"
```
