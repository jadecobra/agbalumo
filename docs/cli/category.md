# agbalumo CLI: Category Management

Manage categories for listings in the agbalumo directory.

## Commands

### category

Manage categories.

```bash
agbalumo category [command]
```

#### Subcommands

##### add

Add a new category to the agbalumo system.

```bash
agbalumo category add [name] [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--claimable` | `-c` | false | Is this category claimable? |

##### list

List all active categories in the database.

```bash
agbalumo category list
```
