# Maintenance Engine (internal/maintenance)

The business logic for project verification, linting, and health checks.

## Constraints
- **Statelessness**: Verification functions should be stateless and side-effect free where possible.
- **Regex Purity**: When performing static analysis on templates, use robust regex patterns that account for whitespace and common HTML variations.
- **Parallelism**: Use `RunParallelCI` for long-running checks to maximize CPU utilization during CI/pre-commit.
