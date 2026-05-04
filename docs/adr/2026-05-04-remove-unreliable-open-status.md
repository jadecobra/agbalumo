# ADR: 2026-05-04 Remove Unreliable Open Status

## Status
Accepted

## Context
The "Open Now / Closed" status indicator was implemented to provide real-time operational status for businesses. However, the system currently lacks per-listing time zone metadata. This led to incorrect "Closed" labels for businesses in different time zones or during midnight shift overlaps, even after fixing specific edge cases. Providing inaccurate data violates the "Zero-Cognitive-Load Curation" standard and degrades user trust.

## Decision
We will remove the "Open Now / Closed" status badge from all UI components (Listing Cards and Detail Modals) and retire the `IsCurrentlyOpen` field from the domain model.

The underlying hours parsing logic will be retained for data enrichment and backfill purposes, but will no longer be surfaced as a live UI signal.

## Consequences
- **Positive**: Eliminated misinformation and reduced UI visual noise.
- **Positive**: Reduced cognitive complexity in the `ListingHandler`.
- **Negative**: Users can no longer see at a glance if a business is currently open without checking the hours manually.
- **Future**: This feature may be re-introduced once a robust time zone engine and per-listing location-aware time context are implemented.
