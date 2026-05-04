# ADR [2026-05-04]: Full-Stack Stress Protocol
**Date**: 2026-05-04 **Status**: Accepted

## 1. Context & User Problem
Currently, the stress test suite (`scripts/benchmark_stress.sh`) only measures **Data Persistence** (SQLite write/read latency). However, the "60-second find" goal depends on the **Full-Stack** (Request -> Handler -> Service -> Repository -> Template Rendering -> Response). Without HTTP-level benchmarks, we have zero visibility into concurrency bottlenecks in the Echo server or template rendering pipeline.

## 2. Decision
Integrate `ab` (Apache Benchmark) into the `benchmark_stress.sh` script to measure **Requests per Second (RPS)** and **P99 Latency**. The script is updated to respect system-wide invariants defined in `.agents/invariants.json` (specifically protocol: https and default port: 8443).

The script will:
1. Seed the database with the requested volume of listings.
2. Start a temporary background instance of the application server.
3. Execute `ab` against the target endpoints.
4. Kill the temporary server instance upon completion.

## 3. The Complexity Kill-Switch (Rationale)
* **User Value**: Ensures the system remains responsive under load, preventing search failures during peak traffic.
* **Performance Budget**: Target P99 latency for SSR should be < 200ms under 50 concurrent requests.
* **Minimalism Check**: Reuses existing `ab` tooling instead of introducing a heavy k6 or Gatling dependency.

## 4. Consequences
* **Technical Tradeoffs**: Requires a running server instance during benchmarks, which adds setup/teardown time to the script.
* **Observability**: Establishes a baseline for local performance that must be maintained as features are added.
* **SQLite Impact**: Forces testing of SQLite concurrency in WAL mode under real HTTP handler load.

## 5. Alternatives Considered
* **k6**: Rejected as it requires a separate JS runtime and is more complex to integrate into simple shell scripts.
* **Go Benchmarks**: Useful for unit-level, but `ab` provides a better simulation of real-world HTTP concurrency and network stack overhead.
