# CI Performance Metrics

| Step | Environment | Baseline (Before) | Optimized (After) | Improvement |
| :--- | :--- | :--- | :--- | :--- |
| **Full Local CI** | M1 Pro | ~1m 52s | ~1m 47s | ~5s (4.5%) |
| **Core Checks (Local)** | M1 Pro | ~40s | ~7s | ~33s (82%) |
| **Remote CI (Total)** | GitHub Runner | - | - | - |

## Local Details (M1 Pro)

- **Total Time (Baseline):** 112s (1m 52s)
- **Total Time (Optimized):** 107s (1m 47s)
- **Note:** Local bottleneck is `go test ./...` with `-race` which takes ~105s regardless of parallelism. Non-test checks dropped from ~40s to ~7s.

## Remote Details (GitHub Runner)

- **Prepare:** N/A
- **Lint:** _(pending)_
- **Tests:** _(pending)_
- **Docker Build:** _(pending)_
