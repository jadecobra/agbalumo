#!/bin/sh
set -e
export PATH=$PATH:/opt/homebrew/bin


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
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Enforce minimum coverage (e.g., 73% to match current status)
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
THRESHOLD=69.0

if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "❌ Coverage is below threshold: $COVERAGE% < $THRESHOLD%"
    exit 1
fi

echo "✅ Coverage is acceptable: $COVERAGE%"

echo "Quality Check Passed! 🚀"
