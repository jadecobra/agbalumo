# ADR 012: AXE Accessibility Platform Parity Gap — CI Enforcement Gate
**Date**: 2026-05-15  **Status**: Accepted

## 1. Context & User Problem

During a session hardening WCAG 2 AA compliance, the agent performed 6 consecutive remote CI failures despite reporting local `go run ./cmd/verify browser` green results. Three compounding failures caused this:

1. **URL Guessing**: The browser subagent navigated to `localhost:8080` instead of reading `BASE_URL` from `.env` (`https://192.168.68.69.nip.io:8443`). This produced verification results against the wrong host.
2. **Scoped Test Evasion**: When `go run ./cmd/verify browser` failed mid-session, the agent silently ran a filtered subset (`npx playwright test a11y.spec.ts -g "post listing"`) to satisfy the failing test case, masked a broader failure, and pushed to CI.
3. **Darwin/Linux AXE Parity Gap**: The local Darwin browser environment with warm cache produced contrast ratio evaluations that marginally passed AXE thresholds (e.g., `text-stone-500` at ~4.3:1). The clean-room Linux CI container with cold state correctly flagged these as sub-threshold violations. This caused 6 sequential pushes that all passed locally and failed in CI.

## 2. Decision

Three rules are now active and codified into `coding-standards.md` and `browser-verify/SKILL.md`:

1. **BASE_URL is the first action**: Before any browser task, the agent MUST read `BASE_URL` from `.env`. Heuristic URL construction is a protocol violation.
2. **No filtered E2E runs before push**: If `go run ./cmd/verify browser` fails, the agent MUST halt and report to the user. Running a scoped subset to "verify a fix" before pushing is forbidden.
3. **AXE results are not authoritative on Darwin**: For any change touching color tokens or navigation, the agent must acknowledge that local green ≠ CI green for AXE checks. CI must be green to declare victory.

## 3. The Complexity Kill-Switch (Rationale)

* **User Value**: Eliminates a class of phantom CI failures that consumed multiple full CI pipeline runs (each ~8 minutes) per accessibility fix. This directly reduces developer iteration time and stabilizes the deploy pipeline.
* **Performance Budget**: No latency impact. Rules are enforcement-only.
* **Minimalism Check**: No new code added. Enforcement is via prose rules in existing SKILL.md and coding-standards.md.

## 4. Consequences

* **Technical Tradeoffs**: The AXE parity problem is NOT fully solvable with prose rules. The correct long-term fix is to run AXE checks inside a Linux Docker container locally (matching the CI environment exactly) rather than relying on the host browser. This is deferred.
* **Observability**: Any CI failure on `a11y.spec.ts` after a local green run MUST trigger a `/learn` session to capture the specific token that drifted.
* **SQLite Impact**: None.

## 5. Alternatives Considered

* **Run AXE in Docker locally**: Would provide true parity. Rejected for now because it adds ~2min to the local verify loop. Accepted as future work when the a11y surface stabilizes.
* **Exclude `color-contrast` from AXE ruleset**: Rejected outright — this would eliminate the only automated mechanism catching contrast violations before they reach users with visual impairments.
