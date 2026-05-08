# Metrics Infrastructure Context

Provides observability and performance monitoring.

## Constraints
*   **Performance**: Metrics collection MUST be non-blocking. Use buffered channels or atomic counters.
*   **Cardinality**: Avoid high-cardinality labels (e.g., UserID, RequestID) in Prometheus metrics. Stick to Routes, Methods, and Status Codes.
*   **Tracing**: Per ADR 2026-05-05, use the `Trace` utility for cross-layer observability without polluting business logic.

## Usage
*   `metrics.IncListingView(listingID)`: Increment view counts.
*   `metrics.RecordLatency(method, duration)`: Track handler performance.
