# Triage & Technical Debt Remediation Plan (Agent Optimized)

This plan prioritizes **structural consolidation** to reduce the Agent's context window usage and eliminate the "noise" of duplicate logic.

## User Review Required

> [!IMPORTANT]
> This plan shifts away from "clean logs first" to "clean structure first." We will consolidate large clones (Phase 1) before addressing remaining string literals (Phase 2). This reduces total LoC and token counts more effectively.

## Proposed Changes

### [Phase 1] Structural Consolidation (Clones)
Address major clone groups in tests and business logic to compress the context window.

#### [MODIFY] [admin_claims_test.go](file:///Users/johnnyblase/gym/agbalumo/internal/module/admin/admin_claims_test.go)
*   Consolidate `TestHandleApproveClaim` and `TestHandleRejectClaim` logic into a shared helper within `testutil` or the local test file.
#### [MODIFY] [listing_mutations.go](file:///Users/johnnyblase/gym/agbalumo/internal/module/listing/listing_mutations.go)
*   Deduplicate repeating mutation logic patterns identified by ChiefCritic.

---

### [Phase 2] Literal Deduplication (String Noise)
Consolidate remaining highly repeated strings identified by `goconst` (172 violations).

#### [MODIFY] [constants.go](file:///Users/johnnyblase/gym/agbalumo/internal/domain/constants.go)
*   Add globally reused strings: `".env"`, `"oauth_state"`, `"Listing not found"`.
#### [MODIFY] [main.go](file:///Users/johnnyblase/gym/agbalumo/cmd/verify/main.go)
*   Update verification tool to use domain/testutil constants.

---

### [Phase 3] Auth Module Test Stabilization
Remove legacy test helpers and standardize on `testutil`.

#### [DELETE] [test_helpers_test.go](file:///Users/johnnyblase/gym/agbalumo/internal/module/auth/test_helpers_test.go)
#### [NEW] [auth_mock.go](file:///Users/johnnyblase/gym/agbalumo/internal/testutil/auth_mock.go)
*   Relocate `MockGoogleProvider` here to follow repository standards.
#### [MODIFY] [handler_google_test.go](file:///Users/johnnyblase/gym/agbalumo/internal/module/auth/handler_google_test.go)
#### [MODIFY] [handler_register_test.go](file:///Users/johnnyblase/gym/agbalumo/internal/module/auth/handler_register_test.go)

---

### [Phase 4] Cognitive Complexity
Resolve failure in the maintenance audit.

#### [MODIFY] [perf.go](file:///Users/johnnyblase/gym/agbalumo/internal/maintenance/perf.go)
*   Decompose `RunPerformanceAudit` into smaller functions.

## Open Questions

- None.

## Verification Plan

### Automated Tests
- `go run cmd/verify/main.go ci`: Run the full CI pipeline.
- `go run cmd/verify/main.go critique`: Verify clones and repeated strings are reduced.
