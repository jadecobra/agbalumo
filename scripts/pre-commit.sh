#!/bin/sh
set -e
# Robust PATH discovery for macOS and Linux
for dir in /usr/local/bin /opt/homebrew/bin /usr/bin /bin; do
    case ":$PATH:" in
        *":$dir:"*) ;;
        *) export PATH="$PATH:$dir" ;;
    esac
done




echo "Running 10x Engineer Quality Checks..."

echo "1. Running Go Fmt..."
if [ -n "$(gofmt -l .)" ]; then
    echo "❌ Go Code is not formatted. Run 'gofmt -w .'"
    exit 1
fi

echo "2. Running Go Mod Tidy Check..."
go mod tidy
# Check for unstaged changes to go.mod or go.sum
if ! git diff --exit-code --quiet go.mod go.sum; then
    echo "❌ go.mod/go.sum are not tidy (unstaged changes detected). Run 'go mod tidy' and commit changes."
    exit 1
fi

echo "3. Running Go Vet..."
go vet ./...

echo "4. Running Tests with Race Detection & Coverage..."
# Run tests on all packages
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Enforce minimum coverage
# Using grep -oE to extract the first decimal number for better robustness
COVERAGE=$(go tool cover -func=coverage.out | grep total | grep -oE "[0-9]+(\.[0-9]+)?" | head -1)
THRESHOLD=90.0

if [ -z "$COVERAGE" ]; then
    echo "❌ Could not parse coverage value"
    exit 1
fi

if [ "$(echo "$COVERAGE < $THRESHOLD" | bc -l)" -eq 1 ]; then
    echo "❌ Coverage is below threshold: $COVERAGE% < $THRESHOLD%"
    exit 1
fi

echo "✅ Coverage is acceptable: $COVERAGE%"

# Note: Secret Scanner moved to scripts/security-check.sh


# Secret scanner logic removed from here and moved to security-check.sh
# which is executed as part of the pre-commit hook via setup-hooks.sh


echo "5. Running API Drift Check..."
bash scripts/api-drift-check.sh

echo "6. Running CLI Drift Check..."
bash scripts/cli-drift-check.sh

echo "Quality Check Passed! 🚀"

echo "7. Running Performance Audit..."
# Exit 2 = critical failures (block commit). Exit 1 = warnings only (allow through).
# We use `|| true` to prevent set -e from aborting on exit code 1 (warnings).
sh scripts/performance-audit.sh || PERF_EXIT=$?
PERF_EXIT="${PERF_EXIT:-0}"
if [ "$PERF_EXIT" -eq 2 ]; then
    echo "❌ Performance audit has critical failures. Fix them before committing."
    exit 1
fi
