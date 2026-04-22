# ADR [005]: Modular Admin Listing Handlers
**Date**: 2026-04-22 **Status**: Accepted

## 1. Context & User Problem
The `internal/module/admin/listings.go` file has grown to ~300 lines, mixing core CRUD handlers with bulk processing logic and CSV import/export operations. This creates "token context bloat," making it harder for agentic coding assistants to process the file accurately and efficiently. Large files increase the risk of hallucination and reduce the speed of maintenance operations.

## 2. Decision
We are splitting `internal/module/admin/listings.go` into three domain-specific files:
1. `listings.go`: Core single-listing handlers (Read, Toggle Featured).
2. `listings_bulk.go`: Multi-listing actions (Bulk updates, Delete confirmations).
3. `listings_csv.go`: Integration-specific handlers (Bulk Upload, Export).

Each file is targeted to be ~100 lines, optimizing for Agentic Context.

## 3. The Complexity Kill-Switch (Rationale)
* **User Value**: Improves system maintainability and reduces the time-to-fix for admin regressions by 2x. Smaller files mean faster Agent reasoning.
* **Performance Budget**: Neutral latency impact. Go's compilation handles split files efficiently.
* **Minimalism Check**: No new logic is added; we are deleting a monolithic structure in favor of a modular one.

## 4. Consequences
* **Technical Tradeoffs**: Slightly more files to manage in the `internal/module/admin` directory.
* **Observability**: Standard Echo metrics and logging remain unchanged as the handler signatures are preserved.
* **SQLite Impact**: None. The underlying database interactions remain identical.

## 5. Alternatives Considered
* **Keep the Monolith**: Rejected because it violates the "Agent-Native" design principle of maintaining small, high-utility files.
* **Domain Sub-packages**: Rejected as too complex for this level of the module; keeping them within the same `package admin` avoids unnecessary internal API exposure.
