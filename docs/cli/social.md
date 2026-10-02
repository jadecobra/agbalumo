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
| `--fail-bad-links` | | false | Fail draft immediately if any deep link returns non-2xx (default: omit bad links with loud stderr; pillar 3 always fails) |

**Dual-Gate Data Integrity:**

To guarantee that public deep links (`https://agbalumo.com/listings/{id}`) never 404 in published social drafts, the engine enforces a two-stage verification gate:

1. **Gate 1: Prod-Parity SQLite**:
   - `agbalumo social draft` requires a production-parity database snapshot containing matching live UUIDs.
   - The database path is wired via environment variables:
     - Preferred: `AGBALUMO_SOCIAL_DB`
     - Fallback: documented `DATABASE_URL`
     - Local default: `.tester/data/prod_snapshot.db` (if present)
   - **Prohibited**: `.tester/data/agbalumo.db` is strictly blocked for social drafting because local synthetic UUIDs do not exist on production.
   - Sync the production snapshot using the one-liner routine:
     ```bash
     ./scripts/pull_prod_db.sh
     # or directly via Fly CLI:
     fly ssh sftp get /data/agbalumo.db .tester/data/prod_snapshot.db
     ```

2. **Gate 2: Deep Link Probe**:
   - Before emitting a draft, every listing deep link (`https://agbalumo.com/listings/{id}`) is verified via HTTP `HEAD` (falling back to `GET`).
   - If a link returns non-2xx (e.g., 404):
     - **Pillar 3 (Spotlight)**: Fails draft generation with an explicit error to prevent spotlighting dead links.
     - **Multi-Listing Pillars (1, 2, 4, 5)**: Omits the link from the listing item with a loud warning printed to `stderr` (or fails the draft if `--fail-bad-links` is supplied).
   - No silent bad URLs can be emitted.

**Content Pillars:**

1. **Pillar 1: Quality Index** — The DFW African Food Quality Index with community-reviewed spots.
2. **Pillar 2: Airport Corridor Cities** — Verified spots in corridor cities (Arlington, Grand Prairie, Irving).
3. **Pillar 3: Merchant Spotlight** — Dedicated spotlight on an individual merchant. Specify `--listing-id` to target a venue, or omit to rotate through unfeatured listings.
4. **Pillar 4: Sub-Metro Corridor Guide** — Highlights North DFW and Collin County spots (Plano, Allen, McKinney, Frisco).
5. **Pillar 5: Radical Transparency & Coverage Gaps** — Top spot per city with counts, blind spots built from the data, and an open call for missing spots.

**Exposure & Rotation Engine:**

- **Exposure over Star-Sort**: Draft generation prioritizes exposure over star-sorting so the same few listings do not dominate every week.
- **Pillars 1, 2, 4, 5**: Listings round-robin within the filtered candidate set across successive CLI runs.
- **Pillar 3**: If `--listing-id` is provided, that merchant is spotlighted directly. If omitted, the CLI rotates through listings while excluding recently featured merchants.
- **Honest Attribution**: Generated links use clean, honest UTM attribution (`utm_source=facebook&utm_medium=social&utm_campaign=<pillar>`), with no fabricated platform variants.

**Examples:**

```bash
# 1. Sync production SQLite database snapshot
./scripts/pull_prod_db.sh

# 2. Generate a post for Pillar 1 (Quality Index) in Dallas using prod-parity DB
AGBALUMO_SOCIAL_DB=.tester/data/prod_snapshot.db agbalumo social draft --pillar 1

# 3. Generate an Airport corridor dispatch (Pillar 2)
AGBALUMO_SOCIAL_DB=.tester/data/prod_snapshot.db agbalumo social draft --pillar 2

# 4. Spotlight a specific live merchant by verified ID (Pillar 3)
AGBALUMO_SOCIAL_DB=.tester/data/prod_snapshot.db agbalumo social draft --pillar 3 --listing-id 3b839bd9-5b71-4771-92a1-21093ae9ff8e

# 5. Spotlight next merchant via rotation (Pillar 3)
AGBALUMO_SOCIAL_DB=.tester/data/prod_snapshot.db agbalumo social draft --pillar 3

# 6. Generate Collin County corridor guide (Pillar 4)
AGBALUMO_SOCIAL_DB=.tester/data/prod_snapshot.db agbalumo social draft --pillar 4

# 7. Save draft to file with strict link failure enforcement
AGBALUMO_SOCIAL_DB=.tester/data/prod_snapshot.db agbalumo social draft --pillar 5 --fail-bad-links --output post_draft.txt
```

