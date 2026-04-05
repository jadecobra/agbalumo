# Goal: Purge Lingering Artifacts and Consolidate Workspace Entropy

## Background
The user noticed workspace entropy: leftover unstructured artifacts from manual test runs, disjointed binary directories (`.bin` vs `bin`), Python scripts lingering in the root, and orphaned Bash generator scripts. 

To maintain the rigorous "No Paperwork/Zero Entropy" standard, we need to execute a targeted sweep of the root workspace, move stray outputs to designated temporary/metrics channels, consolidate binaries, and reorganize our generation workflows.

## Implementation Steps for Gemini 3 Flash

### 1. Consolidate Binary Directories
We have both `.bin` (which contains `verify`) and `bin` (which contains legacy binaries). Pure Go modules should output to a standard `bin/` layer.
```bash
mv .bin/verify bin/verify
rm -rf .bin
```

### 2. Purge Root Build & Test Litter
Agents and test frameworks have been dropping compilation binaries and output logs directly into the root folder.
Delete the orphaned root binaries and test log litter:
```bash
rm -f harness server admin.test
rm -f test_fail.log test_out.log test_output.txt test_results.txt pre-commit-out.txt
```

### 3. Route Coverage Output
Coverage and security type outputs natively belong in the `.metrics` observability folder, not root.
```bash
mv c.out coverage.out sectypes.out .metrics/ 2>/dev/null || true
```

### 4. Relocate `test_ui.py`
The Playwright validation script (`test_ui.py`) running in root causes clutter. It belongs in the structured `.tester` hierarchy or `ui/tests`.
```bash
mkdir -p .tester/scripts
mv test_ui.py .tester/scripts/test_ui.py
```

### 5. Review & Migrate `scripts/generate-juice.sh`
The `generate-juice.sh` script parses `.agents/rules/brand.toon` to generate CSS variables in `ui/static/css/brand-tokens.css` and Go constants in `internal/domain/brand.go`.

**Immediate Action:**
Because this is a deterministic generation engine, it belongs under a `go:generate` directive or a native Go tool, but for this immediate task:
1. Ensure the generated targets (`ui/static/css/brand-tokens.css` and `internal/domain/brand.go`) are actively tracked in Git.
2. If the user permits, this Bash script should be scheduled for porting into a `cmd/juice/main.go` interface. For this pass, do not delete it, but recognize it as the final frontier of the "Pure Go" migration phase.

### 6. Commit the Sweep
Once the workspace is swept, commit the consolidation constraints:
```bash
git add -A bin .tester .metrics
git add -u
git commit -m "chore(workspace): consolidate binaries, route test artifacts to .tester, and purge root runtime litter"
```
