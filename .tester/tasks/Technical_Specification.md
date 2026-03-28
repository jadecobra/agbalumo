# Technical Specification: AST XSS Expansion

Migrate `onclick` and "Dangerous JS" (eval, Function, innerHTML) regex checks from global raw scanning (`checkStructuralRaw`) into AST-based logic (`checkInsecurePatternsGo`) for Go files.

## Objective
Increase the precision of security scans in Go code by moving away from line-based regex and utilizing the Abstract Syntax Tree (AST). This allows for context-aware detection within string literals and respects per-expression suppression comments.

## Architecture

### Component: `internal/agent/security.go`

#### New Function: `checkInsecurePatternsGo`
- Signature: `func checkInsecurePatternsGo(node *ast.File, fset *token.FileSet) []SecurityViolation`
- Logic:
    - Use `ast.Inspect` to find `*ast.BasicLit` nodes where `Kind == token.STRING`.
    - Pattern matching (within string literals):
        - `onclick\s*=`
        - `eval\s*\(`
        - `Function\s*\(`
        - `innerHTML\s*=`
        - `https?://(unpkg\.com|cdn\.jsdelivr\.net|cdn\.tailwindcss\.com|jsdelivr\.net)`
    - Use `isIgnored(n, node, fset)` to check for suppression comments (`#nosec` or `antigravity:allow`).

#### Modification: `checkFile`
- Integration: Call `checkInsecurePatternsGo` within the Go-specific file block.
- Deduplication: Ensure new violations are merged correctly.

#### Modification: `checkStructuralRaw`
- Intent: Prevent duplicate reporting and prioritize the more accurate AST check for Go files.
- Change: Skip `onclick` and `Dangerous JS` regex checks if the file has `.go` extension.

## Data Contracts

### Violation Types
- `Violation.Type`: Keep as "Structural" or transition to "XSS" for these specific patterns.
- *Decision*: Keep as "Structural" to maintain parity with original regex checks, or "XSS" if they are clearly XSS-related. Let's use "Structural" for now to keep it consistent with the existing `structuralPatterns` reporting.

## Verification Requirements

### Red (SDET-Tester) Phase
- Generate failing tests that expect AST-based detection for:
    - Backtick strings containing `onclick`.
    - Concatenated strings containing `eval`.
    - Strings with `#nosec` correctly ignored.

### Green (Implementation) Phase
- All tests in `security_test.go` must pass.
- `harness verify security-static` must execute correctly.

## Approval Required
> [!IMPORTANT]
> - This migration means Go files will no longer be scanned by the `checkStructuralRaw` regexes for `onclick`/`eval`.
> - If these patterns occur outside of string literals (unlikely in valid Go, but possible in malformed code), they might be missed by AST but would have been caught by regex. We accept this as AST focus is on valid, compiled code.
