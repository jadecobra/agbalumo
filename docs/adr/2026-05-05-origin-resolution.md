# ADR 2026-05-05: Origin Signal Resolution

## Context
Agbalumo listings currently show two origin signals:
1. A country flag next to the title (hardcoded to Nigeria/Ghana based on `OwnerOrigin`).
2. An "Origin: [RegionalSpecialty]" badge for Food listings.

### Issues
- **Redundancy**: If both are "Nigeria", it's duplicate information.
- **Conflict**: If `OwnerOrigin` is "Ghana" but `RegionalSpecialty` is "Nigeria" (enriched/verified), the flag is misleading.
- **Ada's Need**: Ada wants authentic signals. A verified food origin is more important than the poster's origin.

## Decision
We will unify the origin signal:
1. **Verified Override**: If `EnrichmentAttemptedAt` is set, the `RegionalSpecialty` (or a derived country) will override the `OwnerOrigin` flag.
2. **Unified Helper**: Implement a `countryFlag` helper to support all African countries from `countries.json`.
3. **Smart Badging**: 
   - Move the flag to a prominent position next to the title.
   - Hide the "Origin" badge if it's identical to the country flag name.
   - If `RegionalSpecialty` is more specific (e.g., "Yoruba"), show it as a secondary signal (Specialty).
4. **Verification Mark**: Add a subtle visual indicator (e.g., a checkmark) if the listing has been enriched.

## Consequences
- **Pros**: Reduced visual noise, more accurate authenticity signals, better utility for Ada.
- **Cons**: Requires mapping logic for specialties if they don't match country names directly.
- **Tradeoff**: We will prioritize exact country matches for flags first; non-country specialties will remain as text badges.
