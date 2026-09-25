# ADR: Social Draft Dual-Gate Architecture
**Date**: 2026-09-25 **Status**: Accepted

## 1. Context & User Problem
The local development and testing SQLite database (`.tester/data/agbalumo.db`) uses synthetic UUIDs that do not exist on the live Fly.io production deployment. Generating social media marketing drafts from this local database resulted in dead public deep links (`https://agbalumo.com/listings/{id}` returning HTTP 404). Furthermore, relying only on HTTP link verification against a test database would draft from the wrong listing set and could spotlight merchants that are not live or miss merchants that are.

## 2. Decision
We implement a deterministic **Dual Gate** for the `agbalumo social draft` engine:

1. **Gate 1 (Prod-Parity SQLite)**:
   - `social draft` queries a production-parity database snapshot containing live UUIDs.
   - Database path resolution prefers `AGBALUMO_SOCIAL_DB`, falling back to documented `DATABASE_URL`, and finally to `.tester/data/prod_snapshot.db` if present.
   - `.tester/data/agbalumo.db` is strictly prohibited and fails execution with an actionable error.
   - Snapshot synchronization is codified in `scripts/pull_prod_db.sh` using `fly ssh sftp get /data/agbalumo.db`.

2. **Gate 2 (Deep Link Probe)**:
   - Before emitting a post draft, each listing deep link (`https://agbalumo.com/listings/{id}`) is verified via HTTP `HEAD` (falling back to `GET` for endpoints where HEAD is not explicitly routed).
   - If a link returns non-2xx:
     - Pillar 3 (Single Merchant Spotlight) immediately fails draft generation.
     - Multi-listing pillars (1, 2, 4, 5) omit the link from the listing item and log a loud warning to `stderr` (or fail immediately if `--fail-bad-links` is supplied).
   - No silent bad URLs can be emitted.

Out of scope: Graph API integrations, automated posting, ID remapping translation tables, and rewriting production state from local.

## 3. The Complexity Kill-Switch (Rationale)
How does this decision respect the 60-second find goal?
* **User Value**: Eliminates broken 404 entry points into the platform from social diaspora groups. Users clicking social recommendations land directly on live merchant profiles in <60 seconds.
* **Performance Budget**: CLI link verification runs concurrently/fast before emit (<500ms total) with zero impact on production runtime HTTP request latencies.
* **Minimalism Check**: Avoided introducing complex ID mapping tables, database slug layers, or headless browser automation for link validation.

## 4. Consequences
* **Technical Tradeoffs**: Developers drafting social copy must pull a fresh production snapshot (`./scripts/pull_prod_db.sh`) periodically.
* **Observability**: Stderr warnings loudly surface omitted non-2xx deep links; CLI exits non-zero on broken spotlight links or prohibited tester databases.
* **SQLite Impact**: The snapshot is read-only locally; Fly production database WAL mode and concurrent access are unaffected by `sftp get` read snapshots.

## 5. Alternatives Considered
* **ID Remapping Table**: Maintained a mapping between local tester UUIDs and Fly production UUIDs. Rejected as high-maintenance complexity creep with severe drift risk.
* **Slug Layer**: Rewrote URLs to use merchant title slugs instead of UUIDs. Rejected as out-of-scope architectural churn altering stable contract `/listings/{id}`.
* **HEAD-only Verification without Prod DB**: Probed public URLs using local database records. Rejected because local DB contains synthetic listings absent from production, leading to empty/failed drafts.
