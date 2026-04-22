# ADR [00X]: [Short, descriptive title]
**Date**: YYYY-MM-DD **Status**: [Proposed | Accepted | Deprecated]

## 1. Context & User Problem
What is the specific user problem we are solving? (e.g., "Users in London can't find Jollof Rice in < 60s due to X"). Why do existing patterns fail to solve this?

## 2. Decision
What exactly are we implementing? State the technology, pattern, or constraint clearly.

## 3. The Complexity Kill-Switch (Rationale)
How does this decision respect the 60-second find goal?
* **User Value**: Why is this 2x better than the current state?
* **Performance Budget**: What is the estimated latency impact? (Goal: < 100ms impact).
* **Minimalism Check**: What component or abstraction was deleted to make room for this change?

## 4. Consequences
* **Technical Tradeoffs**: What becomes harder to maintain?
* **Observability**: How will we monitor if this decision is failing in production?
* **SQLite Impact**: Are there specific locking or concurrency considerations for our disk-parity strategy?

## 5. Alternatives Considered
List at least one simpler approach that was rejected and why it failed the Kill-Switch gate.