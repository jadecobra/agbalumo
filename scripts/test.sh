#!/bin/bash
# scripts/test.sh
set -e
source "$(dirname "$0")/utils.sh"
setup_path

mkdir -p .tester/coverage
go test -json -race -count=1 -coverprofile=.tester/coverage/coverage.out ./... > /dev/null
COVERAGE=$(go tool cover -func=.tester/coverage/coverage.out | awk '/^total:/ {print substr($3, 1, length($3)-1)}')
THRESHOLD_FILE=".agents/coverage-threshold"
THRESHOLD=90.0
if [ -f "$THRESHOLD_FILE" ]; then
    THRESHOLD=$(cat "$THRESHOLD_FILE")
fi
if awk "BEGIN {exit !($COVERAGE < $THRESHOLD)}"; then
    if [ "$FMT" != "json" ]; then
        echo "  ${RED}❌ Error: Coverage is below threshold: $COVERAGE% < $THRESHOLD%${NC}"
        echo "  ${YELLOW}Top 5 lowest coverage files:${NC}"
        go tool cover -func=/tmp/.tester/coverage/coverage.out | grep -v "100.0%" | sort -k 3 -n | head -5 | sed 's/^/    /'
        echo "  ${BLUE}See: .agents/workflows/feature-implementation.md (Gate: coverage)${NC}"
    fi
    exit 2
fi
if [ "$FMT" != "json" ]; then echo "Coverage: $COVERAGE%"; fi
