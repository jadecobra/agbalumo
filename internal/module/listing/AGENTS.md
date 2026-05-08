# Listing Module Context

This module manages the core domain of the application: African food listings. It handles the full lifecycle of a listing from discovery to detailed view and modification.

## Domain Constraints
*   **Open Status**: DO NOT implement "Open Now" badges. Per ADR 2026-05-04, real-time status is decommissioned until time-zone metadata is reliable.
*   **Signature Dishes**: This field is DECOMMISSIONED (ADR 2026-05-05). Use the general description for highlighting specific food items.
*   **Origin Resolution**: Listings must have a verified origin (Country/Region). Use `internal/repository/origins.go` for lookups.

## Handler & UI Constraints
*   **Geocoding**: All listing updates that change address must trigger a geocoding refresh via the `ListingService`.
*   **Images**: Use the `ImageService` for all processing. Never store raw uploads in the main database; store URLs/Paths.
*   **HTMX Fragments**: Favor returning small HTML fragments for "Save" and "Claim" actions to minimize payload size.

## Testing Standards
*   **ADA Compliance**: Every UI change must be verified with `listing_ada_semantic_test.go` to ensure screen-reader compatibility.
*   **Regression**: UI changes MUST pass `ui_regression_home_test.go`.
