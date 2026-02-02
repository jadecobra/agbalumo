# Agbalumo Coding Standards (10x Engineer Edition)

> "We do not write code that breaks. We do not write code without tests. We are 10x."

## 1. The Golden Rule: Verified TDD
*   **Protocol**: Write the test FIRST. Watch it fail (Red). Write the Code. Watch it pass (Green). Refactor.
*   **Mandatory Check**: You must run `./scripts/pre-commit.sh` before submitting any PR or artifact.
*   **No "flaky" tests**: Tests must be deterministic. Use `go test -count=1` to bypass cache if needed.

## 2. Directory Structure & Architecture
*   **`cmd/`**: Main applications. Minimal code.
*   **`internal/`**: Private application and library code.
    *   **`domain/`**: Pure business logic (Structs, Interfaces, Validation methods). NO external dependencies here.
    *   **`handler/`**: HTTP Transport layer (Gin/Echo).
    *   **`repository/`**: Database access layer.
*   **`pkg/`**: Library code ok to use by external apps (if any).

## 3. Go Best Practices
*   **Error Handling**: Return errors, don't panic. Wrap errors with context.
*   **Concurrency**: Use Goroutines for ALL external API calls (Gemini, Database, etc) where appropriate.
*   **Linter**: `go vet` is the minimum standard.

## 4. Domain Specific Rules
*   **Requests Listing**: MUST have a valid `OwnerOrigin` (West African country) and `Deadline`.
*   **Contact Info**: MUST have at least one valid method (WhatsApp, Email, or Phone).
*   **Cultural Context**: All placeholder data and "mock" content must reflect West African culture (Nigerian/Ghanaian focus).

## 5. Agent Protocol
*   **Lead Architect**: Enforces this document.
*   **SDET Agent**: Owns the `*_test.go` files. Backend Agent CANNOT edit test logic without SDET approval.
*   **Backend Agent**: Writes code to pass SDET tests.
