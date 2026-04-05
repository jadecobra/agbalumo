# Phase 15: Resolve CI Pipeline Failure & Dependabot Alert

## Objective
Fix a cross-platform path bug causing `TestCheckGosecRationale` to fail in GitHub Actions (Linux) but pass on macOS, and resolve the newly discovered Dependabot vulnerability in `golang.org/x/image`.

## Context (The Bugs)
1. **The CI Failure (Linux vs macOS)**: 
   In `internal/maintenance/gosec.go`, the scanner intentionally ignores files in the `/tmp/` directory to avoid parsing temp build artifacts:
   ```go
   if strings.Contains(path, "/tmp/") { return nil }
   ```
   Locally on your Mac, `os.MkdirTemp()` creates directories at `/var/folders/...`, so the test cases pass harmlessly. However, on GitHub Actions (Linux runners), `os.MkdirTemp()` creates directories in `/tmp/`. The `CheckGosecRationale` function ends up ignoring its own mock test files, causing the unit test to fail.
2. **Dependabot Alert**:
   GitHub found a moderate vulnerability (`CVE-2026-33809`) in `golang.org/x/image` causing Memory/OOM exhaustion via crafted TIFF files. It requires bumping to `v0.38.0`.

## Steps for Execution

### Step 1: Fix the CI Directory Check
1. Open `internal/maintenance/gosec.go`.
2. Navigate to line 26 where the directory exclusions are defined.
3. Change the blanket string containment check for `/tmp/` to ensure it only skips the *local project* `tmp` folder (e.g. `strings.Contains(path, "/agbalumo/tmp/")`) or remove `/tmp/` from the exclusion list entirely if we don't store project code there. We will simply check if the base directory or any segment equals `tmp`:
   ```go
   // Replace the strings.Contains check with this safer approach:
   pathParts := strings.Split(filepath.ToSlash(path), "/")
   for _, part := range pathParts {
       if part == "vendor" || part == ".tester" || part == "tmp" || part == ".go" {
           return nil // skip
       }
   }
   ```

### Step 2: Patch the Dependabot Alert
1. In your terminal, run the following to explicitly upgrade the vulnerable nested dependency:
   ```bash
   go get golang.org/x/image@v0.38.0
   go mod tidy
   ```

### Step 3: Verification
1. Run `go run cmd/verify/main.go ci` locally one final time.
2. Commit natively: `fix(ci): resolve linux temp path bug in gosec test and patch image CVE`
