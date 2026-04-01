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

TEST_OPTS="-buildvcs=false"

# Use race detector unless SKIP_RACE is set
if [ "$SKIP_RACE" != "true" ]; then
    TEST_OPTS="$TEST_OPTS -race"
fi

# Use -count=1 (disable cache) only in STRICT_MODE
if [ "$STRICT_MODE" == "true" ]; then
    TEST_OPTS="$TEST_OPTS -count=1"
fi

mkdir -p .tester/coverage
TEST_LOG=".tester/test_output.log"
rm -f "$TEST_LOG"

# Execute tests with streaming heartbeats
# We use PIPESTATUS to check for go test failure while piping to awk
set +e 
go test -v $TEST_OPTS -coverprofile=.tester/coverage/coverage.out $PKG 2>&1 | \
    awk -v logfile="$TEST_LOG" '
    /^[[:space:]]*(=== RUN|--- PASS)/ { printf "."; fflush() }
    { print >> logfile }
    '
TEST_EXIT_CODE=${PIPESTATUS[0]}
set -e

# Ensure newline after heartbeats
echo ""

if [ $TEST_EXIT_CODE -ne 0 ]; then
    if [ -f "$TEST_LOG" ]; then
        cat "$TEST_LOG"
    fi
    exit $TEST_EXIT_CODE
fi

# Coverage analysis (Silent unless failed)
COVERAGE=$(go tool cover -func=.tester/coverage/coverage.out | awk '/^total:/ {print substr($3, 1, length($3)-1)}')

if awk "BEGIN {exit !($COVERAGE < $THRESHOLD)}"; then
    if [ "$FMT" != "json" ]; then
        echo -e "${RED}❌ Error: Coverage is below threshold: $COVERAGE% < $THRESHOLD%${NC}"
        echo -e "${YELLOW}Top 5 lowest coverage files:${NC}"
        go tool cover -func=.tester/coverage/coverage.out | grep -v "100.0%" | sort -k 3 -n | head -5 | sed 's/^/    /'
    fi
    exit 2
fi

if [ "$FMT" != "json" ]; then echo -e "✅ Tests Passed. Coverage: ${GREEN}$COVERAGE%${NC} (Threshold: $THRESHOLD%)"; fi
