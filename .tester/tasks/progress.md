# Testing CLI Harness
Comprehensive test metric improvements and testing tools
- [ ] Verify CLI testing harness (api-spec, cli-drift validations), baseline test coverage (>90%), and module-specific coverage targets across the codebase
- [x] Verify os.Exit(1) calls in gate.go, set_phase.go, and status.go are replaced with proper RunE error returns to enable unit testing error paths
- [x] Verify normalizePath utility is extracted to a shared helper in internal/agent to eliminate duplication between drift.go and ast.go
- [x] Verify hardcoded gate strings (e.g. 'red-test', 'api-spec') are extracted into exported constants in internal/agent/state.go
- [x] Verify internal/agent/verify.go test coverage for VerifyRedTest: Write tests covering UI bypass logic, compilation failure, expected failure, and panic/exit evasion
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 1: Route in Code only (fail, missing in OpenAPI and MD)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 2: Route in OpenAPI only (fail, missing in Code and MD)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 3: Route in MD only (fail, missing in Code and OpenAPI)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 4: Route in Code and OpenAPI, but missing in MD (fail)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 5: Route in Code and MD, but missing in OpenAPI (fail)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 6: Route in OpenAPI and MD, but missing in Code (fail)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 7: Command in Code only (fail, missing in Docs)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 8: Command in MD only (fail, missing in Code)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 9: Drift during 'refactor' workflow (assert special error message)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 10: Drift during 'bugfix' workflow (assert special error message)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 11: Drift during 'feature' workflow (assert generic error message)
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 12: ExtractRoutes error (assert false, log: 'Error extracting routes from code')
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 13: swagger-cli bundle command fails (assert false, log: 'Error bundling docs/openapi.yaml')
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 14: ExtractOpenAPIRoutes fails on malformed YAML (assert false, log: 'Error extracting openapi routes')
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 15: os.ReadFile('docs/api.md') fails (assert false, log: 'Error reading docs/api.md')
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 16: ExtractMarkdownRoutes error (assert false, log: 'Error extracting md routes')
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 17: ExtractCLICodeCommands error (assert false, log: 'Error extracting CLI code cmds')
- [ ] Verify internal/agent/verify.go test coverage for VerifyApiSpec - Case 18: ExtractCLIMarkdownCommands error (assert false, log: 'Error extracting CLI md cmds')
- [x] Verify internal/agent/verify.go test coverage for VerifyImplementation and VerifyLint: Write tests asserting correct mock execution of build, test, and vet commands
- [x] Verify internal/agent/verify.go test coverage for VerifyCoverage: Write tests covering profile parsing, threshold enforcement, and violation reporting
- [ ] Verify os.Exit(1) calls in root.go and init.go are replaced with proper RunE error returns to enable unit testing error paths
- [ ] Verify os.Exit(1) calls in gate.go, set_phase.go, and status.go are replaced with proper RunE error returns to enable unit testing error paths
- [ ] Verify os.Exit(1) calls in verify.go, cost.go, and update_coverage.go are replaced with proper RunE error returns to enable unit testing error paths
- [ ] Refactor summarizeProgress and checkAndApplyProgressUpdate to return errors
- [ ] Update init.go and verify.go to handle errors
- [ ] Add unit tests for error paths in root_test.go
- [ ] Created verify_apispec_test.go with 9 tests covering success, API drift, CLI drift, workflow-specific messages (refactor/bugfix), and error paths
- [ ] Coverage improved from 61.9% to 85.7% for VerifyApiSpec
- [ ] Created 7 new test cases covering UI bypass (clean/rejected), compilation failure, compilation-failed-from-JSON, all-tests-pass gate rejection, pattern-matched, and pattern-not-matched branches
- [ ] Refactored mock exec factories into reusable makeMockExec helper
- [x] Verify internal/agent/state.go test coverage for SaveState and calculateSignature: Write tests ensuring proper file creation, permissions, and valid SHA256 generation
- [x] Verify internal/agent/state.go test coverage for LoadState: Write tests covering valid load, missing file, structural mismatch anti-cheat, and signature mismatch
- [x] Verify map[string]interface{} usage in checkAndApplyProgressUpdate (cmd/harness/commands/root.go) is replaced with strongly-typed structs
# Code Refactoring
Reducing duplication and streamlining API handlers
- [ ] Verify Pagination Standardization
- [ ] Verify standard ErrorResponse JSON structure on 400 Error
- [ ] Verify Response Helpers (Pending): Check for usage of `RenderSuccess`, `JSONSuccess`, and `ErrorResponse` handlers across controllers
- [ ] Verify Base Handler Struct (Pending): Review handler method signatures to confirm they use a base struct to eliminate session parsing duplication
- [ ] Verify Service Layer Consolidation (Pending): Inspect service layer structures and ensure validation middleware runs consistently
- [x] Verify Route Grouping (Completed): Inspect code and ensure auth-required routes are fully separated from public routes via grouped middleware
# Phase 3 Roadmap (Pending)
Features and architectural improvements planned for the next phase
- [ ] Verify Apple/Facebook Platforms (Pending): Click on 'Apple' and 'Facebook' login buttons in the identity section and configure dev accounts
- [ ] Verify DevLogin Hardening (Pending): Ensure hitting the DevLogin route in a prod-like condition executes the real claim workflow instead of a skip mechanism
- [ ] Verify Admin Customization (Pending): Navigate to Admin settings and inject a custom dashboard color or font, then confirm the page updates dynamically
- [ ] Verify Expiration Ticker (Pending): Confirm that the background service correctly tags listings older than their configured TTL config as 'expired'
- [ ] Verify Location Filter (Pending): Perform a map-based search query like `?filterLocation=City` and assert that the returned result list matches accurately
- [ ] Verify Browser Integration Test (Pending): Trigger the browser integration test suite and verify 'Create Listing' finishes flawlessly start to end
- [ ] Verify Security Hardening (Pending): Identify dependencies utilizing alpine:latest and ensure that CVE-2025-60876 patches are updated
- [ ] Verify Fat Handler Structs (Pending): Use structural tag libraries rather than custom parsing to validate object bindings
# Verification
Verified necessity of MockImageService in listing tests and removed it where possible.
- [ ] Replaced MockImageService with real LocalImageService using t.TempDir() in TestHandleCreate_WithImage to avoid I/O pollution.
- [ ] Replaced MockImageService with real LocalImageService using t.TempDir() in TestListingHandler_Upload_Valid.
- [ ] Removed unused testifyMock import from listing_upload_test.go.
- [ ] Confirmed MockImageService is still needed in TestHandleCreate_ImageUploadError to simulate upload failures.
- [ ] Implemented TestVerifyApiSpec_FullDriftReport for complete drift aggregation
- [ ] Implemented TestVerifyApiSpec_EmptyFiles for minimal docs edge cases
- [ ] Implemented TestVerifyApiSpec_MarkdownPathEdgeCases for path normalization
- [ ] Implemented TestVerifyApiSpec_ExtractMarkdownRoutes_Direct for unit testing extraction logic
- [ ] Verified all tests pass with GOCACHE/GOTMPDIR /tmp configuration
- [ ] Passed all harness gates: red-test, api-spec, implementation, lint, coverage, browser-verification
- [x] Research legitimate stable SHA for setup-task (v2.0.0: b91d5d2c96a56797b48ac1e0e89220bf64044611)
- [x] Replace invalid SHA b9ed8b34f8a84c3563456885eb0515156a6451df (v2.1.0) with valid v2.0.0 SHA in .github/workflows/ci.yml
- [x] Verify no remaining invalid SHA references
# Context Cost Reduction
Granular refactoring tasks to reduce the context cost of the top 10 most expensive files
- [x] Verify monolithic test files (e.g., handler_test.go, sqlite_test.go, listing_write_test.go), frontend scripts (app.js), and CLI command definitions have been successfully extracted into modular, context-efficient files while maintaining test suite integrity
- [ ] Verify Interface Segregation (CQRS): Extract ListingStore into ListingReader and ListingWriter interfaces in internal/domain to eliminate sqlite_listing.go structural fragmentation
- [ ] Verify Test Suite Isolation: Refactor residual large test files like sqlite_listing_test.go and listing_validation_test.go to use testify/suite or targeted sub-test groupings to further reduce context bloat
- [ ] Verify Domain Map Documentation: Create a domain mapping reference in .agents/rules/ indicating where specific business validations (Create, Geocoding, Jobs) are physically located to prevent context discovery thrashing
# CI/CD Pipeline Improvements
Unified CI/CD pipeline using Go-Task with security-first patterns, coverage anti-degradation, and deep local-to-cloud parity.
- [x] Verify CI/CD: Taskfile optimizations (binary checks, dynamic GOBIN), optional GOPATH isolation, unified coverage anti-degradation, and hardened environment respect in agent scripts verified via 'task ci'
# Go-Task Migration
Replace orchestrator scripts with declarative Taskfile while migrating legacy checkers to AST.
- [x] Verify Go-Task Migration: Replaced orchestrator scripts with declarative Taskfile, migrated regex-based security checks (XSS/SQLi) to AST-based logic, and verified via 'task ci:security'
# Architectural Compliance
Implemented harness context optimization and progress archiving.
- [x] Verify Architectural Compliance: Truncated verification outputs, automated progress archiving, and silenced Taskfile/script noise implemented
# Feedback Loop Optimization
Sequential dev loop with auto-formatting, watch mode, and agent protocols.
- [x] Verify 10x Feedback Loop: Executing 'task fmt', 'task pre-commit', and './scripts/watch.sh' confirms sequential Fast vs Heavy orchestration and agent protocols in GEMINI.md
# ChiefCritic: User Journey Audit
Enhanced ChiefCritic persona with browser-based user journey validation and conversion metrics. Blueprint: file:///Users/johnnyblase/.gemini/antigravity/brain/e7d5db52-4f58-4a5b-97d6-f41db4f6ed0c/implementation_plan.md
- [x] Verify ChiefCritic Journey Audit: Integrated browser-based user journey validation (user_journeys.yaml, audit skill, persona update) with MockGoogleProvider auth and orchestration script (scripts/browser_audit.sh) established
# Harness Utilities
Implemented autonomous agent spawning capability to enable background orchestration.
- [x] Implement SpawnAgent in internal/agent/util.go
- [x] Add unit tests for SpawnAgent in internal/agent/util_async_test.go
- [x] Verify process detachment via SysProcAttr.Setsid
- [x] Pass all harness gates for harness-agent-spawn feature (REFACTOR phase)
# Security Hardening
Validated and excluded G117 struct-secret-pattern linter warnings for configuration fields.
- [x] Applied #nosec G117 to SessionSecret and GoogleMapsAPIKey in config.go
- [x] Applied #nosec G117 to APIKey in geocoding.go
- [x] Verified mandatory rationale comments with hyphen separator
- [x] Passed all security static analysis and internal scanners
- [x] Verified zero regressions in build and test suite
# Audit Security Hardening
Implemented a 10x Security Reasoning Framework for the audit_security skill.
- [x] Hardened audit_security skill with STRIDE and Chaos Injection mandates.
- [x] Implemented mandatory security_audit.md artifact and compliance archiving.
- [x] Verified skill compliance via automated TDD-driven validation.
