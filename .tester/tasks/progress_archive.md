# Core Infrastructure
Database, performance tuning, and API modularization
- [ ] Verify SQLite optimizations, Network/Load headers (gzip, Cache-Control), and Domain-Driven Handler dependency injection (cmd/server/main.go)
# Performance & Benchmarks
Optimizing long-running tasks and CLI benchmarks
- [ ] Verify CLI benchmarks, slow-query logging (SLOW_QUERY_THRESHOLD_MS), stress generation progress UI, and parallel execution (goroutines/errgroup)
# Security
Security posture and secrets management
- [ ] Verify security headers (CSP, HSTS, CSRF), track/scrub test artifacts from git history, and run secret scanner pre-commit gate
# Authentication & Sessions
Manage user authentication flow and secure cookies
- [ ] Verify OAuth/Session cookie policies (SameSite=Lax/Strict) and cross-site OAuth callback validation
# UI Framework & Design System
Global Theme Rollout & Design System Refinement to Earth theme
- [ ] Verify Earth theme palette (bg-earth/text-earth), typography (Inter/Playfair), and standard component styles across all pages
# UI Modals & Interactions
Refactoring modals and adding micro-animations
- [ ] Verify modal_base shell reusability, Dark Theme integration, and global micro-animations (hover:scale-105) across all modals
# Admin UI & Dashboard
Admin panel features and styling
- [ ] Verify Bulk Upload functionality, Admin Users view, and Dark Theme rendering across the Admin Dashboard metrics/listings
# Business Features & Logic
Core application capabilities and user workflows
- [ ] Verify Claim Ownership workflow, Homepage Service listings, and Status Badge rendering across search/categories
# API & UI Route Verification
Verified all public, user, and admin routes are accessible and functional via the UI
- [ ] Verify successful navigation of Public (Home/About/Search), User (Profile/Create), Admin routes, and client form backend URL normalizer validations
# CI/CD & Build Tooling
Continuous integration, vulnerability scanning, and build workflows
- [ ] Verify Tailwind CLI compilation (no CDN), GitHub Actions pipeline (govulncheck, docker scout), and Fly.io deployment config validation
# Modular Services & Domain Utilities
Refactoring features into dedicated services and domain packages
- [ ] Verify refactored ImageService, domain-level Hours of Operation string parsing, decoupled AuthMiddleware, and strict date parsing workflows
# Security Audit & Hardening
Comprehensive security audit fixes and defense-in-depth improvements
- [ ] Verify ADMIN_CODE logic, Admin Login rate limits, Content-Security-Policy compliance (no inline scripts), HTMX bindings (`hx-on`), and local asset hosting
# User & Admin Journeys Mapping
Comprehensive mapping of application routes and workflows for all user personas
- [ ] Verify Visitor, Authenticated User, and Admin comprehensive journeys, and local DB CLI tool maintenance integration
# Standardize output for Agent
Refactor all scripts in the scripts/ directory to output structured JSON for robust agent consumption
- [ ] Verify script outputs (harness, drift check, CI gates, tests, utils) conform to standard JSON envelope contracts
# Agent Context Optimization
Reduce token footprint of agent context by breaking monolithic specs and documentation into granular files.
- [x] Verify monolithic specs (OpenAPI, BRAND_GUIDELINES, CLI commands) are successfully extracted into granular modular files
- [x] Verify codebase documentation (AGENTS.md, Feature Workflows) is fully modularized into granular rules under `.agent/rules/` and `.agent/workflows/`
- [x] Verify feature-implementation workflows (refactor.md) enforce proactive Code Smell evaluation, safe continuous refactoring loops, and strict test preservation
- [x] Verify Context Cost Math: Implement CalculateContextCost and RMS mathematical logic
- [x] Verify Context Cost Tests: Implement unit tests in internal/agent/cost_test.go covering whitelist/blacklist file filtering and correct RMS output
- [x] Verify Context Cost CLI: Integrate harness cost command into cmd/harness/main.go
- [x] Verify Context Cost Execution: Run harness cost command and validate terminal output aesthetic
# Featured Listings Per Category Implementation
Implemented category-specific featured listings with a max limit of 3, persisting across pagination
- [ ] Added a failing test to verify featured prioritization on page 2.
- [ ] Removed `page == 1` restriction in `HandleHome` and `HandleFragment`.
- [ ] Verify ListingStore GetFeaturedListings interface accepts category parameter
- [x] Verify SQLite GetFeaturedListings query filters by category and enforces LIMIT 3 constraint
- [x] Verify Admin HandleToggleFeatured prevents featuring more than 3 items per category
- [x] Verify listing HandleHome calls GetFeaturedListings correctly with empty category
- [x] Verify listing HandleFragment passes category filter to GetFeaturedListings during HTMX pagination
- [x] Verify mock repository implementations match the updated GetFeaturedListings signature
- [ ] Verify failing test implementation enforces admin max 3 featured listings constraint
- [ ] Verify failing test implementation asserts category-specific featured persistence upon pagination
- [ ] Verify listing HandleFragment passes category filter to GetFeaturedListings during HTMX pagination
# Bugfixes
Fixed Admin Featured Listings HTMX row swap turning into JSON
- [ ] Verify if MockImageService in TestHandleCreate_WithImage and TestListingHandler_Upload_Valid to resolve filesystem I/O errors during tests is still needed
- [ ] Verify `HandleToggleFeatured` in `admin_listings.go` returns JSON natively.
- [ ] Verify `setupAdminTestContext` mocking the renderer to expect an admin_listing_table_row. is necessary
- [ ] Verify `admin_actions_test.go` explicitly verifies HTML is returned.
- [ ] Verify Refactored `HandleToggleFeatured` returns `updatedListing` HTML upon successful save.
- [ ] Verify Fixed jq JSON parsing bug in scripts/utils.sh when failures contain unescaped characters
# Harness Security & Anti-Cheat
Fix exploits and bypasses in the 10x Engineer CLI harness to ensure strict TDD adherence
- [x] Verify JSON Casing Exploit: Modify state.json with uppercase keys and ensure LoadState rejects the invalid signature
- [x] Verify Coverage Spoofing Exploit: Attempt to bypass coverage by modifying .agents/coverage-thresholds.json and ensure Harness catches it
- [x] Verify Red-Test Bypass Exploit: Ensure 'gate red-test PASS' cannot be manually executed via CLI
- [x] Verify Panic Evasion Exploit: Ensure tests that os.Exit(0) without outputting failure correctly trigger red-test failure
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
# Harness Utilities (Decommissioned)
- [x] ~~Implement SpawnAgent in internal/agent/util.go~~ (Removed: Environment constraints)
- [x] ~~Add unit tests for SpawnAgent in internal/agent/util_async_test.go~~ (Removed)
- [x] ~~Verify process detachment via SysProcAttr.Setsid~~ (Removed)
- [x] ~~Pass all harness gates for harness-agent-spawn feature (REFACTOR phase)~~ (Removed)
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
# Learning Loop Core Logic
Implemented the infrastructure for capturing and persisting Squad-Decision-Summaries.
- [x] Create internal/history/history.go
- [x] Define SquadDecision struct with all required personas
- [x] Implement Store function with YYYYMMDD_HHMMSS_<feature>.md naming scheme
- [x] Use internal/util/fs wrapper for cross-platform safety
- [x] Verify file creation and content formatting with TDD (internal/history/history_test.go)
# Learning Loop CLI
Implemented the aglog CLI tool for capturing squad decisions.
- [x] Create cmd/aglog/main.go with flag and JSON support
- [x] Create cmd/aglog/main_test.go with unit tests
- [x] Update internal/history/history.go to return saved path
- [x] Reach 100% test coverage for logic and 90%+ for CLI
- [x] Pass all harness gates for aglog-cli feature
# Auth Component Refactoring
Fix SSRF vulnerability in GetUserInfo by using authenticated oauth2 client and header-based authentication.
- [x] Refactor GetUserInfo to use p.config.Client
- [x] Use http.NewRequestWithContext for safer request creation
- [x] Update #nosec comments with proper SSRF rationale (G107, G704)
- [x] Update existing tests to expect Authorization header
- [x] Verify fix with new reproduction test
# Infrastructure Alignment
Stricter local pre-commit gates aligned with production CI to ensure global reliability.
- [x] Research current `Taskfile.yml` linting targets and `_pre-commit-fast` setup.
- [x] Research `progress.md` for `SpawnAgent` decommissioning.
- [x] Research `build-feature` documentation for any `SpawnAgent` remnants.
- [x] Modify `Taskfile.yml`:
- [x] Rename `lint` to `lint:staged`.
- [x] Create `lint` pointing to `ci:lint`.
- [x] Update `_pre-commit-fast`.
- [x] Modify `.tester/tasks/progress.md`:
- [x] Mark `SpawnAgent` as decommissioned.
- [x] Verify:
- [x] Run `task lint` on a clean `main` and verify it project-scans.
- [x] Run `task ci` to verify the full suite.
- [x] Cleanup: Commit changes to `Taskfile.yml` and `progress.md`.
- [x] Unified linting (ci:lint) for all local commits.
- [x] Full codebase build (ci:build) on every pre-commit.
- [x] Comprehensive test suite (ci:test) with race detection in pre-commit-heavy.
# Decommissioning SpawnAgent
Removed SpawnAgent utility due to functional issues and environment constraints.
- [x] Removed SpawnAgent from internal/agent/util.go
- [x] Deleted internal/agent/spawn_test.go and internal/agent/util_async_test.go
- [x] Cleaned up unused imports (fmt, os, syscall) in internal/agent/util.go
- [x] Removed background agent spawning from harness verify command
- [x] Verified system stability with 100% test pass on remaining gates
# Gitleaks Rules Migration
Replaced mock rules with official default rules in .gitleaks.toml.
- [x] Initialized harness for gitleaks-rules
- [x] Merged default rules with project allowlist
- [x] Verified scan passes with new ruleset
# DevOps Optimization
Implemented granular GitHub Actions caching for Go tools using a composite action to reduce CI overhead.
- [x] Create composite action `.github/actions/setup-task-with-cache`
- [x] Refactor `ci.yml` to use the new action across all relevant jobs
- [x] Implement tool caching for `.tester/tmp` based on `go.sum` and `Taskfile.yml`
