# ADR-20260505: Decommissioning Unreliable 'Signature Dish' Signal

## Context
The "Signature Dish" feature (internally `TopDish`) was implemented as an "Ada Signal" to provide users with high-quality food recommendations automatically extracted from business websites. However, the current extraction logic—which relies on scraping `h1`/`h2` tags for "signature" or "special" keywords—has proven unreliable in production. It frequently extracts generic text or incorrect items, leading to a degraded user experience and violating the "Zero-Cognitive-Load Curation" standard.

## Decision
We will decommission the "Signature Dish" feature across the entire stack. This involves:
1. Removing the `TopDish` field from the `Listing` domain model and database schema mappings.
2. Disabling the scraping logic responsible for populating this signal.
3. Removing the UI components that display "Signature: [Dish]" in cards and modals.
4. Cleaning up form fields and associated tests.

We will NOT remove the database column `top_dish` in this phase to maintain migration compatibility, but it will be ignored by the application code.

## Consequences
- **User Trust**: Improved by removing inaccurate and potentially confusing data.
- **Complexity**: Reduced cognitive and technical overhead in the scraper and data models.
- **Latency**: Minor reduction in template rendering and database scanning time.
- **Data Loss**: Existing "Signature Dish" data will no longer be visible or accessible through the UI/API.
