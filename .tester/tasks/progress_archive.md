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
