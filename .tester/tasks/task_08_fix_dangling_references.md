# Goal: Purge Dangling References & Seed UI Component Foundation

## Background
Following the "Pure Go" migration, several core documentation files (`DEVELOPMENT.md`, `CODING_STANDARDS.md`, `README.md`) were left pointing to deleted Bash scripts, creating dangerous documentation drift. Furthermore, `generate-juice.sh` was broken by the deletion of custom `.toon` files, and `internal/maintenance/cost.go` references a deleted `Taskfile.yml`. 

We must also seed the UI component constraint in our workflows to ensure Agents utilize `tailwind.config.js` properly without hallucinating raw CSS.

## Implementation Steps for Gemini 3 Flash

### 1. Delete Broken Bash Scripts
Remove the obsolete brand generation script completely.
```bash
git rm scripts/generate-juice.sh
```

### 2. Update Documentation Drift
Modify the markdown files to point to their native `cmd/verify` Go alternatives:

- **`docs/DEVELOPMENT.md`**:
  - Replace references to `./scripts/ci-local.sh` with `go run cmd/verify/main.go ci`
  - Replace references to `./scripts/performance-audit.sh` with `go run cmd/verify/main.go perf`
  - Replace references to `./scripts/security-check.sh` with `go run cmd/verify/main.go audit`
  - Fully delete references to `./scripts/agent-exec.sh workflow status` and `./scripts/agent-gate.sh` as they are now strictly forbidden paperwork protocols.

- **`README.md` & `docs/CODING_STANDARDS.md`**:
  - Replace instructions to run `./scripts/verify_restart.sh` with `go run cmd/verify/main.go watch` (or standard `go build`).

### 3. Remove Dead Code References
Edit `internal/maintenance/cost.go`:
- Remove `"Taskfile.yml": true` from the `cost.go` exclusion map (around line 47) since the file no longer exists.

### 4. Seed the UI Component Architecture
Set up the environment so Agents stop guessing at CSS utility classes and use pre-compiled templates.
- **Create Sandbox**: `mkdir -p ui/templates/components/`
- **Enforce UI Constraint**: Append the following rule to `.agents/workflows/coding-standards.md`:

> **UI Component Constraint**: Do NOT write raw HTML form, button, or modal elements directly into major page templates. You MUST utilize or propose existing component definitions located inside `ui/templates/components/`. All colors, spacing, and visual effects must be derived exclusively from `tailwind.config.js` logic natively—no arbitrary hex codes or custom CSS outside of the Tailwind engine.

### 5. Finalize the Cleanup
Commit everything to the repository:
```bash
git add -u
git add .agents/workflows/coding-standards.md
git commit -m "docs(agents): purge dangling script references and seed UI component constraint framework"
```
