# Goal: The Final Polish (Linter Conventions & Stubborn Shell Dependencies)

## Background
While hunting down broken references, a few deeper structural inconsistencies were identified that directly impact developer experience (DX) and violate Go conventions. 

Specifically:
1. The `.golangci.yml` config is hiding inside the `scripts/` directory. This breaks standard IDE integrations (like VSCode's Go extension), which natively scan the *root* of the project to enforce rules as you type. 
2. `verify_restart.sh` and `browser_audit.sh` are fully orphaned ghost scripts. Scanning the entire repository reveals that absolutely nothing calls them (they were likely artifacts of the now-deleted `audit.md` manual paperwork workflow).

## Implementation Steps for Gemini 3 Flash

### 1. Root the Linter Configuration
Move the GolangCI-Lint configuration file out of the scripts directory into the project root where the language server (`gopls` / VSCode) expects it automatically.
```bash
mv scripts/.golangci.yml .golangci.yml
```

### 2. Remove Filepath Hardcoding
Since `.golangci.yml` is now in the root dir, update the hardcoded overrides in testing constraints:
- In `.github/workflows/ci.yml`, within the formatting/linting step (`golangci/golangci-lint-action`), remove the flag `--config=scripts/.golangci.yml`.
- In `cmd/verify/main.go`, inside the `runLinter()` func, update the shell execution arguments: Change `.args("run", "-c", "scripts/.golangci.yml")` to simply `.args("run")` (or `"run", "--new-from-rev=HEAD"` where applicable) so it picks it up naturally.

### 3. Obliterate the Ghost Scripts
Since nothing in the repository calls `verify_restart.sh` or `browser_audit.sh`, delete them permanently to stop the bleeding of outdated shell constraints.
```bash
git rm scripts/verify_restart.sh scripts/browser_audit.sh
```

### 4. Validate & Commit
Ensure all unit tests and quality gates still pass with the new layout.
```bash
go run cmd/verify/main.go ci
git add -A
git commit -m "chore: move golangci config to root and delete orphaned browser_audit / verify_restart scripts"
```
