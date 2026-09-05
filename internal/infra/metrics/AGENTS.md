# Metrics Infrastructure Context

Provides observability and performance monitoring.

## Constraints
*   **Persistence**: Persist metrics through `domain.ListingRepository.SaveMetric` using `domain.Metric`.
*   **Logging**: Tag log output with `[ADA-METRIC]` prefix via `slog` for external log aggregation.
*   **Tracing**: Per ADR 2026-05-05, use the `Trace` utility for cross-layer observability without polluting business logic.

## Usage
*   `service.LogAndSave(ctx, eventType, value, metadata)`: Captures a metric, persists to DB, and logs to stdout.
