# Chaos Brief - Context Cost Optimization

## Mission Brief
As **ChaosMonkey**, your goal is to sabotage the context cost optimization efforts to ensure the squad's monitoring and refactoring are resilient against "cost drift" and "ghost files". 

## Sabotage Targets

### 1. The Ghost File Injection
- **Target**: `CalculateContextCost` logic in `internal/agent/cost.go`.
- **Action**: Create a large dummy file in a directory that *should* be ignored (e.g., `.agents/temp/ghost.md` with 5000 lines) BUT ensure it's actually counted due to a logic flaw you introduce, OR find a way to make it look like a valid file in a non-ignored dir. 
- **Goal**: Artificially inflate the RMS to prove the monitoring is brittle.

### 2. The Multi-File Circular Fragility
- **Target**: The split `internal/agent/security.go` files.
- **Action**: Introduce a circular dependency between `security_web.go` and `security_sql.go` during the refactoring phase. 
- **Goal**: Break the build and ensure the `BackendEngineer`'s modularization is sound.

### 3. The Exemption Bypass
- **Target**: The `ignoredFiles` map in `cost.go`.
- **Action**: Rename `critique_report.md` to `critique_report_legacy.md` and see if the cost tool still ignores it (it shouldn't if it's hardcoded).
- **Goal**: Test the robustness of the exclusion rules.

## Success Conditions (Resilience)
- The squad detects the "Ghost File" inflation and adjusts the `ignoredDirs` or `ignoredFiles` logic.
- The build failure from circular dependencies is caught by `task lint` / `task test` before phase transition.
- The cost report correctly identifies the "new" legacy report if it wasn't explicitly excluded by pattern rather than name.

## Failure Conditions (Brittle)
- The RMS cost reported is wildly inaccurate and the squad doesn't notice.
- Circular dependencies are merged into master.
- Large reports continue to bloat the context without the squad's awareness.
