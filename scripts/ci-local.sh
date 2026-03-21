#!/bin/bash
# scripts/ci-local.sh
# Helper script to run GitHub Actions locally using act

# Check for format flag
FMT="json"
if [ "$1" = "--text" ]; then 
    FMT="text"
    shift # Remove the flag so it isn't passed to act
fi

source "$(dirname "$0")/utils.sh"

# Ensure we are in the root directory
cd "$(dirname "$0")/.."

# Check if act is installed
if ! command -v act &> /dev/null; then
    if [ "$FMT" = "text" ]; then
        echo "❌ act is not installed. Please install it with 'brew install act'."
    else
        output_json_envelope false "ci-local.sh" "act is not installed. Please install it with 'brew install act'."
    fi
    exit 1
fi

# Detect Apple M-series (arm64) and apply architecture flag
ARCH_FLAG=()
if [[ $(uname -m) == "arm64" ]]; then
    if [ "$FMT" = "text" ]; then echo "Detected Apple M-series chip. Using --container-architecture linux/amd64"; fi
    ARCH_FLAG=("--container-architecture" "linux/amd64")
fi

# Run act
if [ "$FMT" = "text" ]; then echo "🚀 Running local CI with act..."; fi

SUCCESS=true
if [ "$FMT" = "text" ]; then
    act "${ARCH_FLAG[@]}" "$@" || SUCCESS=false
else
    # Capture output for JSON but allow failure
    ACT_OUT=$(act "${ARCH_FLAG[@]}" "$@" 2>&1) || SUCCESS=false
fi

# Run full performance benchmarks
if [ "$FMT" = "text" ]; then
    echo ""
    echo "📊 Running full search performance benchmarks (10,000 listings)..."
    go test -json -v -bench=BenchmarkSearchPerformance ./internal/repository/sqlite/search_performance_test.go || SUCCESS=false
else
    BENCH_OUT=$(go test -json -v -bench=BenchmarkSearchPerformance ./internal/repository/sqlite/search_performance_test.go 2>&1) || SUCCESS=false
fi

if [ "$FMT" = "json" ]; then
    if [ "$SUCCESS" = "true" ]; then
        output_json_envelope true "ci-local.sh" "CI and benchmarks passed successfully."
    else
        # To avoid massive output in JSON, we just pass the captured streams 
        COMBINED_OUT="ACT OUTPUT:\n$ACT_OUT\n\nBENCHMARK OUTPUT:\n$BENCH_OUT"
        output_json_envelope false "ci-local.sh" "$COMBINED_OUT"
    fi
fi

if [ "$SUCCESS" = "false" ]; then exit 1; fi
