#!/bin/bash
# scripts/test.sh - Optimized test runner for agbalumo.
set -e
source "$(dirname "$0")/utils.sh"
setup_path

# Configuration
THRESHOLD_FILE=".agents/coverage-threshold"
THRESHOLD=80.0
if [ -f "$THRESHOLD_FILE" ]; then
    THRESHOLD=$(cat "$THRESHOLD_FILE")
fi

PKG="./..."
if [ "$1" != "" ]; then
    PKG="$1"
fi

TEST_OPTS="-json -buildvcs=false"

# Use race detector unless SKIP_RACE is set
if [ "$SKIP_RACE" != "true" ]; then
    TEST_OPTS="$TEST_OPTS -race"
fi

# Use -count=1 (disable cache) only in STRICT_MODE
if [ "$STRICT_MODE" == "true" ]; then
    TEST_OPTS="$TEST_OPTS -count=1"
fi

mkdir -p .tester/coverage
# Execute tests
go test $TEST_OPTS -coverprofile=.tester/coverage/coverage.out $PKG > /dev/null

# Coverage analysis
COVERAGE=$(go tool cover -func=.tester/coverage/coverage.out | awk '/^total:/ {print substr($3, 1, length($3)-1)}')

if awk "BEGIN {exit !($COVERAGE < $THRESHOLD)}"; then
    if [ "$FMT" != "json" ]; then
        echo "  ${RED}❌ Error: Coverage is below threshold: $COVERAGE% < $THRESHOLD%${NC}"
        echo "  ${YELLOW}Top 5 lowest coverage files:${NC}"
        go tool cover -func=.tester/coverage/coverage.out | grep -v "100.0%" | sort -k 3 -n | head -5 | sed 's/^/    /'
        echo "  ${BLUE}See: .agents/workflows/feature-implementation.md (Gate: coverage)${NC}"
    fi
    exit 2
fi

if [ "$FMT" != "json" ]; then echo "Coverage: $COVERAGE%"; fi
