---
description: Simulate load and benchmark system constraints
---

When the user asks to stress test the system, benchmark an endpoint, or find load constraints:

1. **Protocol**: You MUST execute `bash scripts/benchmark_stress.sh` to trigger the stress test framework.
2. **Analysis**: Analyze the output for P95/P99 latency spikes, HTTP status failures, and connection drops.
3. **Architectural Handoff**: Depending on bottleneck (e.g., DB locks vs CPU), suggest relevant optimizations (connection pooling, fast-path caching, SQLite WAL adjustments) to the user natively.
