# Chaos Brief - The Ada Context (60-Second Quality)

## Mission Brief
As **ChaosMonkey**, your goal is to sabotage the feature build and infrastructure to ensure its resilience. Specifically, you must target the CI/CD pipeline and the build process to ensure that my fixes are not just "fixing the immediate error" but also robust against future accidental regressions.

## Sabotage Targets

### 1. The Build-Breaker
- **Target**: `.github/workflows/ci.yml` or `Taskfile.yml`.
- **Action**: Inject an invalid go directive or a syntax error in the workflow.
- **Goal**: Verify that the squad's monitoring and your own verification identifies the regression and rejects the build immediately.

### 2. The Dependency Saboteur
- **Target**: `go.mod` or `package.json`.
- **Action**: Force an incompatible version of a dependency (e.g., `golang.org/x/image` version downgrade) that triggers a build failure in CI.
- **Goal**: Ensure the "Dependabot" failures from the logs can be detected and remediated by the squad's SDET.

### 3. The Action-Sha Bypass
- **Target**: Pinned SHAs in `ci.yml`.
- **Action**: Change a pinned SHA for a core action to an older but "valid" version that introduces a security vulnerability.
- **Goal**: Check if the SecurityEngineer persona flags the downgraded/insecure action SHA.

---
*Success Conditions (Resilience)*: 
- Build fails and the harness correctly identifies the error location.
- SDET-Tester can reproduce the failure and pass it back for repair.
- Security scan identifies insecure/stale action versions.

*Failure Conditions (Brittle)*:
- Build passes locally but fails remotely with the same "unable to resolve action" error.
- Regressions go unnoticed by the verification gates.
