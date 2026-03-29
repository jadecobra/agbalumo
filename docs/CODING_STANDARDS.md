# agbalumo Coding Standards (10x Engineer Edition)

> "We do not write code that breaks. We do not write code without tests. We are 10x."

## 1. The Golden Rule: Verified TDD
*   **Protocol**: Write the test FIRST. Watch it fail (Red). Write the Code. Watch it pass (Green). Refactor.
*   **Mandatory Check**: You must run `./scripts/verify_restart.sh` before submitting any PR or artifact. This script runs tests, checks coverage, and restarts the server.
*   **Coverage Threshold**: Enforced value lives in `.agents/coverage-threshold`. NEVER lower this value, even to bypass failing tests. If coverage drops, you MUST write more tests. Lowering the threshold is strictly forbidden.
*   **No "flaky" tests**: Tests must be deterministic. Use `go test -count=1` to bypass cache if needed.
*   **Write small, single-purpose functions by default (SRP, clean code)**

## 2. Directory Structure & Architecture
*   **`cmd/`**: Main applications. Minimal code.
*   **`internal/`**: Private application and library code.
    *   **`domain/`**: Pure business logic (Structs, Interfaces, Validation methods). NO external dependencies here.
    *   **`handler/`**: HTTP Transport layer (Gin/Echo).
    *   **`repository/`**: Database access layer.
*   **`pkg/`**: Library code ok to use by external apps (if any).

## 3. Go Best Practices
*   **Error Handling**: Return errors, don't panic. Wrap errors with context.
*   **Concurrency**: Use Goroutines for ALL external API calls (Gemini, Database, etc) where appropriate.
*   **Linter**: `go vet` is the minimum standard.

## 4. Domain Specific Rules
*   **Requests Listing**: MUST have a valid `OwnerOrigin` (West African country) and `Deadline`.
*   **Contact Info**: MUST have at least one valid method (WhatsApp, Email, or Phone).
*   **Cultural Context**: All placeholder data and "mock" content must reflect West African culture (Nigerian/Ghanaian focus).

## 5. Agent Protocol (Preferred Models)
*   **ProductOwner** (`GEMINI_3_1_PRO`): Value & Usefulness. Owns the "Why" and "MVDP" strategy. Ultimate say on Two-Way Doors.
*   **SystemsArchitect** (`GEMINI_3_1_PRO`): Technical Integrity. Owns the TDD loop, Go structure, and One-Way Door decisions.
*   **SDET-Tester** (`GEMINI_3_FLASH`): Owns `*_test.go`. Writes failing tests (Red) BEFORE implementation.
*   **BackendEngineer** (`GEMINI_3_FLASH`): Implements logic to pass tests. **Strict Rule: No code without a failing test.**
*   **SecurityEngineer** (`GEMINI_3_1_PRO`): Owns `security_test.go` and audits. "Trust but verify."
*   **UIUXDesigner** (`GEMINI_3_1_PRO`): Owns look/feel (HIG/Material 3). Ensures "User Delight" and FCP < 1.0s. Reads browser testing URLs from `.agents/rules/browser-url.md`.
*   **CulturalModerator** (`GEMINI_3_FLASH`): Ensures content relevance via Gemini. "Validate first."
*   **ChiefCritic** (`GEMINI_3_1_PRO`): Relentless reviewer. Tears down code, architectures, and UI to ensure absolute adherence to 10x standards and high quality UX.
*   **ChaosMonkey** (`GEMINI_3_1_PRO`): Resiliency Lead. Proactively injects faults, corrupts state, and simulates "Cheating" to verify that the squad's safety gates and anti-cheat mechanisms are truly unbreakable.

## 6. Security Best Practices (Zero Trust)
*   **Input Validation**: Validate ALL inputs at the boundary. No implicit trust.
*   **Secure Headers**: Ensure all HTTP responses include secure headers (HSTS, CSP, etc.).
*   **Dependency Auditing**: Regularly check for vulnerabilities in dependencies.
*   **Least Privilege**: Application parts should only have the permissions they absolutely need.

## 7. 10x Operational Standards
*   **Minimal/Modular**: Changes must be the smallest possible increment that passes the test. Avoid over-engineering.
*   **Validate**: A change is not done until it is verified to pass the test AND perform as expected (manual or automated check).
*   **Test First**: No code is written without a failing test. This is non-negotiable.

## 8. UI/UX Standards (Premium & Delightful)
*   **Standards**: Follow [Apple HIG](https://developer.apple.com/design/human-interface-guidelines) and [Material 3](https://m3.material.io/).
*   **Delight Rule**: Every interaction (click, hover, focus) MUST have visual feedback.
*   **Premium Rule**: Pixel-perfect alignment, consistent spacing (8pt grid), and thoughtful typography.
*   **Fun Rule**: The app should feel alive. Use transitions and micro-animations to surprise and delight (e.g., confetti on success, smooth ease-in/out).

## 9. Performance Standards (Latency is the Enemy)
*   **Backend**: Critical path operations (Validation, Parsing) must be benchmarked. Budget: < 1000ns/op for strict logic.
*   **API**: 99p response time must be under 100ms.
*   **Frontend**: First Contentful Paint (FCP) < 1.0s. Minimize client-side JS (Use HTMX).
*   **Database**: No N+1 queries. Use `EXPLAIN` on all complex queries.

## 10. Drift Prevention Protocol
*   **Traceability**: Every Pull Request or major artifact update must cite the specific Standard it adheres to (e.g., "Implements Standard 8.2: Delight Rule").
*   **Enforcement**: The Lead Architect validates that the implementation matches the Plan and the Persona.
*   **Pre-flight**: Agents must self-correct by reviewing their specific instructions before determining the plan.
*   **Sync Rule**: Changes to `.agents/config.yaml` or `.agents/personas/*.yaml` MUST be mirrored in Section 5 of this document immediately. "Double-Commit" is required.
