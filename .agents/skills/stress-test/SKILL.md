---
name: "Stress Test"
description: "Simulate load, benchmark SQLite concurrency, and measure HTTP endpoint response times under pressure."
triggers:
  - "/stress-test"
  - "stress test"
  - "benchmark system"
  - "load test"
mutating: false
---

# Stress Test Skill

Benchmark system performance limits, measure endpoint latency, and stress test SQLite concurrency.

## Execution
1. Run `bash scripts/benchmark_stress.sh` to execute the load test framework.
2. Inspect output metrics for 95th and 99th percentile response times, throughput, and error rates.
3. Compare measured latencies against the 500 millisecond time-to-first-result target.

## Bottleneck Analysis
1. If SQLite lock contention causes spikes, evaluate transaction sizes and WAL configuration pragmas.
2. If template rendering is the bottleneck, evaluate component fragment caching opportunities.
3. Present findings and specific optimization recommendations to the user.
