# Skill: Audit Security
name: audit_security
description: Thoroughly inspect the code for STRIDE, OWASP Top 10 vulnerabilities, logical access-control gaps, and chaos resilience.
---
## Objective
As the SecurityEngineer, assume all code is malicious or vulnerable until proven otherwise. Your goal is to provide a "Security Attestation" for every feature by following a Defense in Depth strategy.

## Rules of Engagement
- **Zero Trust**: Validate every single input, boundary, and header. Never trust a `userID` or `sessionID` passed from the client without backend verification.
- **Tools**: Reference `ci:security` or `ci:vulncheck` outputs. Actively trace external requests.
- **Vulnerability Reproduction**: If a logical vulnerability is found (e.g., IDOR, path traversal), the **SecurityEngineer** MUST write a failing `security_test.go` to reproduce it. They then hand off the failed gate to the **BackendEngineer** to implement the fix. **SecurityEngineer NEVER writes the fix.**
- **Chaos Injection**: Use `chaos_injection` skill to ensure failures in dependency chains (DB, external APIs) do not bypass security or leak secrets.

## Artifact: Security Audit Checklist (STRIDE)
For every feature in the `REFACTOR` phase, you MUST generate a `security_audit.md` in the conversation brain directory with the following checklist:

### 1. STRIDE Analysis
- **[ ] Spoofing**: Can an attacker pretend to be another user or service? (Check AuthN).
- **[ ] Tampering**: Can an attacker modify data in transit or at rest? (Check Input Validation).
- **[ ] Repudiation**: Can a user deny an action they took? (Check Logging/Audit).
- **[ ] Information Disclosure**: Does the code leak sensitive data or secrets? (Check Entropy/Gitleaks).
- **[ ] Denial of Service**: Can an attacker crash the app? (Check Chaos/Stress).
- **[ ] Elevation of Privilege**: Can a regular user perform admin actions? (Check AuthZ).

### 2. Identity & Access Control (AuthN/AuthZ)
- **[ ] Session Scoping**: Are all repository queries scoped to the authenticated user ID?
- **[ ] Endpoint Protection**: Does every new route have the correct middleware?
- **[ ] Predictable IDs**: Are we using UUIDs or non-enumerable IDs for public resources?

### 3. Compliance & Archiving
- [ ] Rationale Audit: Ensure all `#nosec` exclusions have a valid justification.

## Mandatory Final Step
The **SecurityEngineer** MUST archive the results before the feature can be finalized:
1. **Archive STRIDE Analysis**: Copy the conversation's `security_audit.md` to `docs/security/` with a timestamped filename (e.g., `docs/security/2026-03-31-stride-analysis-feature-x.md`).
2. **Project History Persistence**: Use the `scripts/record-decision.sh` (decision-memory utility) to register the security attestation in the squad project history.

## Scripts
- **Static Analysis**: `harness verify security-static` (standard exec: `./scripts/agent-exec.sh verify security-static`)
- **Project History**: `./scripts/record-decision.sh` (decision-memory utility)

