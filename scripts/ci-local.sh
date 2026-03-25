#!/bin/bash
# scripts/ci-local.sh
# Helper script to run CI checks natively

FMT="json"
if [ "$1" = "--text" ]; then 
    FMT="text"
    shift
fi

source "$(dirname "$0")/utils.sh"
setup_path
export GOCACHE=/tmp/gocache
export GOTMPDIR=/tmp/gotmp
mkdir -p "$GOCACHE" "$GOTMPDIR"

# Ensure we are in the root directory
cd "$(dirname "$0")/.."

LOG_DIR=$(mktemp -d /tmp/ci-local-XXXXXX)
trap 'rm -rf "$LOG_DIR"' EXIT

if [ "$FMT" != "json" ]; then 
    echo "${BLUE}${BOLD}🚀 Running Local CI...${NC}"
    echo "   (matched to pre-commit gates: lint, test, coverage, security)"
fi

# Define tasks matching pre-commit logic
# 1. Lint
run_task "lint" "GolangCI-Lint" "$LOG_DIR" golangci-lint run -c scripts/.golangci.yml &

# 2. Test & Coverage
check_tests() {
    go test -json -race -count=1 -coverprofile=/tmp/coverage.out ./...
}
run_task "test" "Tests & Coverage" "$LOG_DIR" check_tests &

# 3. Security
run_task "security" "Security Check" "$LOG_DIR" sh scripts/security-check.sh &

# 4. Benchmarks (Legacy ci-local functionality)
run_task "bench" "Benchmarks" "$LOG_DIR" \
    go test -json -v -bench=BenchmarkSearchPerformance ./internal/repository/sqlite/search_performance_test.go &

# Wait for all background tasks
FAILURES=0
for job in $(jobs -p); do
    wait $job || FAILURES=$((FAILURES + 1))
done

if [ "$FMT" = "json" ]; then
    SUCCESS=true
    if [ $FAILURES -gt 0 ]; then SUCCESS=false; fi
    
    # Consolidate outputs for the envelope
    COMBINED_OUT=""
    for log in "$LOG_DIR"/*.log; do
        [ -e "$log" ] || continue
        COMBINED_OUT="$COMBINED_OUT\n--- $(basename "$log") ---\n$(cat "$log")"
    done
    
    output_json_envelope "$SUCCESS" "ci-local.sh" "$COMBINED_OUT"
    exit $FAILURES
fi

if [ $FAILURES -eq 0 ]; then
    echo ""
    echo "${GREEN}${BOLD}✅ All CI checks passed!${NC}"
    exit 0
else
    echo ""
    echo "${RED}${BOLD}❌ CI failed with $FAILURES failures.${NC}"
    exit 1
fi
