# ADR-20260409: Pivot to 'Ada' User Persona

**Date**: 2026-04-09 **Status**: Accepted

## Context
Agbalumo currently serves four distinct pillars (Business Directory, Job Board, Events, and Community Requests). This broad scope has led to increased cognitive complexity in the codebase, fragmented UI navigation, and a lack of clear value proposition for the end-user. We need a way to focus development efforts on a specific, high-value user segment to ensure the product remains maintainable and useful.

## Decision
We have decided to pivot the platform's primary focus to **Ada, the Quality-Obsessed Consultant**. 

- **Primary Use Case**: High-quality ethnic food discovery in < 60 seconds.
- **Secondary Use Case**: Finding specialized cultural services (specifically Tailors for African weddings).
- **Demoted Pillars**: Job Board, Events, and generic Community Requests will be hidden from the primary navigation and landing page to simplify the user journey.

## Consequences
- **Easier**: The UI will become significantly cleaner, focusing only on Search and Discovery for Food/Services.
- **Easier**: Complexity in the `Listing` logic can be streamlined to prioritize quality signals (verified status, peer trust) over generic metadata for jobs/events.
- **Maintenance**: Unused features (Jobs/Events) will still exist in the database and background logic for now but will not receive feature updates.
- **Risk**: Existing users (if any) relying on Job/Event pillars will find them harder to access. However, at this stage of the project, focusing on a single "winning" persona (Ada) is prioritized over broad utility.
