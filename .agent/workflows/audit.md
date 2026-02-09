---
description: comprehensive project audit (tests, security, ui, performance)
---
1. **Test Coverage Analysis**
    - Run the full test suite with coverage:
      ```bash
      go test -coverprofile=coverage.out ./...
      ```
    - Display coverage function report to identify weak spots:
      ```bash
      go tool cover -func=coverage.out
      ```
    - **Goal**: Ensure overall coverage is >80% and critical paths (domain/handlers) are covered.

2. **Security Audit**
    - **Custom Security Tool**:
      ```bash
      go run cmd/security-audit/main.go
      ```
      - Verify HTTPS (HSTS), CSP, and X-Frame-Options.
    
    - **CVE Scanning (Docker Scout)**:
      ```bash
      # Ensure image is built
      docker build -t agbalumo-secure:latest .
      
      # Run CVE scan
      docker scout cves local://agbalumo-secure:latest
      ```
      - **Goal**: Identify and remediate Critical/High vulnerabilities.

3. **UI & UX Review**
    - Start the server if not running:
      ```bash
      go run cmd/server/main.go
      ```
    - **Browser Verification**: Use the `browser_subagent` to:
      - Visit `https://localhost:8443`
      - Take a screenshot of the home page.
      - Check console for errors (CSP, 404s, JS crashes).
      - Verify responsive layout on mobile/desktop viewports.
      - Click through critical flows (Create Listing, View Detail, Search).

4. **Performance Check**
    - Monitor server logs during browser navigation for slow requests (>500ms).
    - Check browser network tab for large assets (images >500KB, JS bundles).
    - ensure database queries are efficient (no N+1 issues visible in logs).

5. **Reporting**
    - Compile findings into a markdown report (e.g., `audit_report.md`).
    - Create tasks for any regressions or failures found.
