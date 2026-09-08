# agbalumo CLI: Social Content Engine

Generate publication-ready social media copy and distribution drafts directly from the listing database.

## Commands

### social

Manage and generate social media content and distribution assets.

```bash
agbalumo social [command]
```

#### Subcommands

##### draft

Draft high-utility social media copy from listing database across defined storytelling pillars.

```bash
agbalumo social draft [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--pillar` | `-p` | 1 | Content pillar (1: Quality Index, 2: Airport Arrival, 3: Merchant Spotlight, 4: Sub-Metro Corridor, 5: Coverage Gaps) |
| `--platform` | | "facebook" | Target platform (facebook, x) |
| `--city` | | "Dallas" | Target city/corridor |
| `--listing-id` | | "" | Target listing ID (required for pillar 3) |
| `--output` | `-o` | "" | Optional path to write drafted copy |
