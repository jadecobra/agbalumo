# Goal: Obliterate Obsolete Agent Workflows

## Background
The recent migration to a "Pure Go" verification toolchain (`cmd/verify/main.go`) and the strict enforcement of the `AGENTS.md` rulebook (e.g., "No Paperwork" for explicit progress tracking) have made several legacy markdown workflow constraints totally obsolete, redundant, or broken.

To prevent agents from utilizing decommissioned constraints, running undefined shell scripts, or attempting chaotic local SAST strategies outside the native `ci` phase, we must delete them.

## Implementation Steps for Gemini 3 Flash

1. **Delete Obsolete Workflow Documentation:**
   Run the following command to permanently remove the deprecated workflow constraints from the `.agents/workflows` directory:
   ```bash
   git rm .agents/workflows/janitor.md .agents/workflows/audit.md .agents/workflows/update-coverage.md .agents/workflows/restart-server.md
   ```

2. **Verify Deletion:**
   Ensure the directory `.agents/workflows/` only contains active, Go-integrated workflows (like `build-feature.md`, `coding-standards.md`, `stress-test.md`, and `deploy-secrets.md`).

3. **Commit:**
   Commit the deletion:
   ```bash
   git commit -m "docs(agents): obliterate deprecated and broken workflow modules"
   ```

4. **Verify Toolchain Health:**
   Run the native CI pipeline one last time to ensure documentation drift and legacy dependency removals did not break any constraints:
   ```bash
   go run cmd/verify/main.go ci
   ```
