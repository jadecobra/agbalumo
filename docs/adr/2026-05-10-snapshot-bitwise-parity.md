# ADR 2026-05-10: Bitwise Snapshot Parity Enforcement

**Date**: 2026-05-10 **Status**: Accepted

## 1. Context & User Problem
The project mandates visual snapshot parity between Darwin (local development) and Linux (CI environment). However, the existing `verify snapshot-parity` tool only checked for file existence. This created a loophole where agents or developers could "fake" parity by copying `*-darwin.png` to `*-linux.png`. Since rendering engines and font antialiasing differ between macOS and Linux, these "fake" snapshots caused non-deterministic failures in production CI, leading to wasted compute and delayed deployments.

## 2. Decision
We are implementing a strict **Bitwise Identity Check** in the `verify snapshot-parity` tool.
*   The tool will now read both the Darwin and Linux variants of a snapshot.
*   If the files are bitwise identical, the tool will emit a failure: "Snapshot is bitwise identical to its darwin counterpart; it was likely copied rather than generated on Linux."
*   Linux snapshots MUST be generated in a genuine Linux environment (e.g., via Docker) before pushing.

## 3. The Complexity Kill-Switch (Rationale)
*   **User Value**: Guarantees that the UI looks correct in production (Linux) without relying on "hope" or manual cross-referencing. Prevents the "it works on my machine" anti-pattern for UI changes.
*   **Performance Budget**: Impact is < 10ms per snapshot (local file I/O). Negligible relative to the cost of a failed CI run.
*   **Minimalism Check**: Replaced a weak prose rule ("Please don't copy snapshots") with a deterministic, automated gate.

## 4. Consequences
*   **Technical Tradeoffs**: Developers on macOS MUST have Docker installed to generate Linux snapshots locally. This increases the "setup cost" but ensures 100% CI parity.
*   **Observability**: Failures are clearly reported in the `pre-push` hook and `verify` CLI.
*   **SQLite Impact**: None.

## 5. Alternatives Considered
*   **LLM Visual Audit**: Rejected due to non-determinism and high latency/token cost.
*   **Darwin-only Snapshots**: Rejected because Linux rendering is the source of truth for our production users.
