# Skill: Audit Security
name: audit_security
description: Thoroughly inspect the code for STRIDE, OWASP Top 10 vulnerabilities, logical access-control gaps, and chaos resilience.
---
## Objective
As the SecurityEngineer, assume all code is malicious or vulnerable until proven otherwise. Your goal is to provide a "Security Attestation" for every feature by following a Defense in Depth strategy.

## Rules of Engagement
- **Zero Trust**: Validate every single input, boundary, and header. Never trust a `userID` or `sessionID` passed from the client without backend verification.
- **Tools**: Reference `ci:security` or `ci:vulncheck` outputs. Actively trace external requests.
- **Vulnerability Reproduction**: If a logical vulnerability is found (e.g., IDOR, path traversal), you MUST write a failing `security_test.go` to reproduce it before fixing.
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
- **[ ] Rationale Audit**: Ensure all `#nosec` exclusions have a valid justification.
- **[ ] Compliance**: Upon feature finalization, COPY the `security_audit.md` to `docs/security/` with a timestamped filename (e.g., `docs/security/2024-03-30-audit_security_hardening.md`).

## Scripts
- Default: `scripts/agent-exec.sh verify security-static`
- Chaos: `scripts/benchmark_stress.sh`
