# Infrastructure Alignment
Stricter local pre-commit gates aligned with production CI to ensure global reliability.
- [x] Unified linting (ci:lint) for all local commits.
- [x] Full codebase build (ci:build) on every pre-commit.
- [x] Comprehensive test suite (ci:test) with race detection in pre-commit-heavy.
