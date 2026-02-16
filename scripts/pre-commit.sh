#!/bin/sh
set -e
export PATH=$PATH:/opt/homebrew/bin


echo "Running 10x Engineer Quality Checks..."

echo "1. Running Go Vet..."
go vet ./...

echo "2. Running Tests with Coverage..."
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Enforce minimum coverage (e.g., 80%)
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
THRESHOLD=70.0

if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "❌ Coverage is below threshold: $COVERAGE% < $THRESHOLD%"
    exit 1
fi

echo "✅ Coverage is acceptable: $COVERAGE%"

echo "Quality Check Passed! 🚀"
