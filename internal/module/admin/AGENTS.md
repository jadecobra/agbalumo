# Admin Module Context

This module provides high-privilege operations for site moderators and owners. 

## Architectural Constraints (Agent-Native)
*   **File Splitting**: Per ADR 2026-04-22, listing operations are split to reduce token context bloat:
    *   `listings.go`: Core single-listing CRUD.
    *   `listings_bulk.go`: Multi-listing actions and confirmations.
    *   `listings_csv.go`: Import/Export and external integrations.
*   **Thin Handlers**: Keep admin handlers focused on request binding and permission checks. Business logic belongs in services.

## Security & Operations
*   **RBAC**: Every handler MUST use the `AdminOnly` middleware.
*   **Bulk Safety**: Bulk delete or update operations MUST require a confirmation step (modal or separate page).
*   **Modals**: Use `modals.go` patterns for consistent HTMX-driven administrative dialogs.

## Data Management
*   **CSV Integrity**: Bulk uploads must be validated for required fields (Name, Origin, Address) before any database writes occur.
*   **Exports**: Use streaming CSV generation for large listing exports to avoid memory exhaustion.
