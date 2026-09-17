# agbalumo CLI: Social Content Engine

Generate publication-ready social media copy and distribution drafts directly from the listing database following agbalumo core story principles.

## Commands

### social

Manage and generate social media content and distribution assets.

```bash
agbalumo social [command]
```

#### Subcommands

##### draft

Draft high-utility social media copy from the listing database across defined storytelling pillars.

```bash
agbalumo social draft [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--pillar` | `-p` | 1 | Content pillar (1-5) |
| `--listing-id` | `-l` | "" | Target listing ID (for pillar 3 spotlight; rotates if omitted) |
| `--city` | | "Dallas" | Target city or metro anchor |
| `--output` | `-o` | "" | Optional path to write drafted copy |

**Content Pillars:**

1. **Pillar 1: Quality Index** — The DFW African Food Quality Index with community-reviewed spots.
2. **Pillar 2: Airport Corridor Cities** — Verified spots in corridor cities (Arlington, Grand Prairie, Irving).
3. **Pillar 3: Merchant Spotlight** — Dedicated spotlight on an individual merchant. Specify `--listing-id` to target a venue, or omit to rotate through unfeatured listings.
4. **Pillar 4: Sub-Metro Corridor Guide** — Highlights North DFW and Collin County spots (Plano, Allen, McKinney, Frisco).
5. **Pillar 5: Radical Transparency & Coverage Gaps** — City-by-city breakdown of mapped spots and open call for missing community gems.

**Exposure & Rotation Engine:**

- **Exposure over Star-Sort**: Draft generation prioritizes exposure over star-sorting so the same few listings do not dominate every week.
- **Pillars 1, 2, 4, 5**: Listings round-robin within the filtered candidate set across successive CLI runs.
- **Pillar 3**: If `--listing-id` is provided, that merchant is spotlighted directly. If omitted, the CLI rotates through listings while excluding recently featured merchants.
- **Honest Attribution**: Generated links use clean, honest UTM attribution (`utm_source=cli&utm_medium=social&utm_campaign=<pillar>`), with no fabricated platform variants.

**Examples:**

```bash
# Generate a post for Pillar 1 (Quality Index) in Dallas
agbalumo social draft --pillar 1

# Generate an Airport corridor dispatch (Pillar 2)
agbalumo social draft --pillar 2

# Spotlight a specific merchant by ID (Pillar 3)
agbalumo social draft --pillar 3 --listing-id cli-12345

# Spotlight next merchant via rotation (Pillar 3)
agbalumo social draft --pillar 3

# Explore Collin County corridor guide (Pillar 4) in Plano
agbalumo social draft --pillar 4 --city Plano

# Save the draft directly to a text file for editing
agbalumo social draft --pillar 5 --output post_draft.txt
```
