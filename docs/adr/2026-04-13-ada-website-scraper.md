# ADR: Heuristic-based Website Scraper for Ada Enrichment

## Status
Accepted

## Context
Ada discovery requires high-trust "cultural signals" (Heat Level, Payment Methods, Menu URL) that are often missing from initial listings. Restaurant websites are highly irregular (SPA, PDF menus, image-only menus). We need a way to enrich listings without heavy external dependencies or manual entry.

## Decision
We implemented a lightweight, heuristic-based scraper built into the `internal/service` layer.
- **Standard Library Only**: Used `net/http` and `golang.org/x/net/html` for speed and stability.
- **Sensory Heuristics**:
    - **Heat Level**: Token frequency mapping of "spicy" keywords to a 1-5 scale.
    - **Payment Methods**: Literal scanning for P2P payment platforms (Zelle, Venmo).
    - **Top Dish**: Heading extraction combined with "signature" keyword proximity.
- **Worker Pattern**: Implemented `ScraperJob` to orchestrate enrichment for listings missing signals.
- **Defense in Depth**: 
    - 15-second timeouts to prevent hung jobs.
    - `io.LimitReader` (512KB) to prevent OOM from malicious or oversized restaurant pages.
    - Cognitive complexity handled by splitting token processing into specific tag handlers.

## Consequences
- **Pros**: Zero overhead, no external scraping bill, fast execution.
- **Cons**: Will miss data on heavily obfuscated or JS-only websites.
- **Maintenance**: Heuristics may need tuning as restaurant site patterns evolve.
