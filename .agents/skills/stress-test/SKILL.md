---
name: "Stress Test"
description: "Simulate load, benchmark SQLite concurrency, measure HTTP endpoint response times, and preserve baseline profile artifacts for comparative performance tuning."
triggers:
  - "/stress-test"
  - "stress test"
  - "benchmark system"
  - "load test"
  - "performance optimization"
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

## Performance Tuning & Benchmark Comparisons
1. **Preserve Baseline Artifacts**: When tuning performance, capture and preserve baseline benchmark profiles in `.tester/profiles/` (e.g., `search_cpu.pprof`, `search_mem.pprof`). Ensure the directory exists before initiating benchmarks (`mkdir -p .tester/profiles`).
2. **Retain Artifacts Until Comparative Completion**: Do not delete artifact directories before post-optimization benchmarks complete.
3. **Automated Comparative Verification**: Run post-optimization comparisons using identical flags (`-bench`, `-benchtime`, `-cpuprofile`, `-memprofile`) and present side-by-side latency/allocation deltas before declaring the task complete.
