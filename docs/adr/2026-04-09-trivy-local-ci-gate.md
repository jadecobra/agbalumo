# ADR: Add Trivy Image Scan to Local CI Gate

* **Status:** Decided
* **Date:** 2026-04-09

## Context

Three classes of CVEs have reached production CI because local CI only scans the app's own module graph (`govulncheck`). Third-party binaries in the Docker image (such as `litestream`) have their own dependency trees that are invisible to `govulncheck`. This gap allows container-level regressions to reach production undetected.

## Decision

Add a Trivy image scan after the Docker build step in the `ci --with-docker` command. The scan will use identical flags to those used in production CI (`ci.yml`):
- `--exit-code 1`
- `--ignore-unfixed`
- `--vuln-type os,library`
- `--severity CRITICAL,HIGH`

Trivy becomes a hard dependency when the `--with-docker` flag is used.

## Consequences

- Developers must install Trivy (`brew install trivy`) to use the `--with-docker` flag.
- Build and scan process takes approximately 5–10 minutes locally.
- No CVE in any compiled binary within the image can reach production CI without first failing the local validation gate.
