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
*   **SDET Agent**: Owns the `*_test.go` files (Functional/Integration). Backend Agent CANNOT edit test logic without SDET approval.
*   **Security Engineer**: Owns `security_test.go` and security policy. specific audits.
*   **Backend Agent**: Writes code to pass SDET tests.

## 6. Security Best Practices (Zero Trust)
*   **Input Validation**: Validate ALL inputs at the boundary. No implicit trust.
*   **Secure Headers**: Ensure all HTTP responses include secure headers (HSTS, CSP, etc.).
*   **Dependency Auditing**: Regularly check for vulnerabilities in dependencies.
*   **Least Privilege**: Application parts should only have the permissions they absolutely need.

## 7. 10x Operational Standards
*   **Minimal/Modular**: Changes must be the smallest possible increment that passes the test. Avoid over-engineering.
*   **Validate**: A change is not done until it is verified to pass the test AND perform as expected (manual or automated check).
*   **Test First**: No code is written without a failing test. This is non-negotiable.
