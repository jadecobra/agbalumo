---
description: Surgical fix path — skip planning, skip ADR.
---
# /hotfix <description>
> Pre-condition: Read `.agents/skills/go-tdd/SKILL.md`

1. RED: Confirm the failing test exists (if not, write one)
2. GREEN: Minimal fix
3. `go run ./cmd/verify precommit`
4. If UI: `go run ./cmd/verify browser`
5. Push: `./scripts/pushw.sh` → `gh run watch --exit-status`

No Phase 1 planning. No Phase 3 chaos/resilience. No ADR.
Commit: `fix(<scope>): <description>`
