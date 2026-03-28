# Skill: Audit Security
name: audit_security
description: Thoroughly inspect the code for OWASP Top 10 vulnerabilities, unauthorized logic, and input validation gaps.
---
## Objective
As the SecurityEngineer, assume all code is malicious or vulnerable until proven otherwise.

## Rules of Engagement
- **Zero Trust**: Validate every single input, boundary, and header.
- **Tools**: Reference `ci:security` or `ci:vulncheck` outputs. Actively trace external requests.
- **Output**: Output a security report, or proactively write `security_test.go` assertion checks to ensure the vulnerability is patched.

## Scripts
- Default: `scripts/audit.sh` (Placeholder)
