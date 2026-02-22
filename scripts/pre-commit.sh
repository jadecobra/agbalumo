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
# Run tests on all packages
go test -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Enforce minimum coverage (89.2%)
COVERAGE=$(go tool cover -func=coverage.out | grep -v "mock" | grep total | awk '{print substr($3, 1, length($3)-1)}')
THRESHOLD=89.2

if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "❌ Coverage is below threshold: $COVERAGE% < $THRESHOLD%"
    exit 1
fi

echo "✅ Coverage is acceptable: $COVERAGE%"

echo "5. Running Secret Scanner..."
FAILED=0

# 1. Check filenames (staged)
# Matches .env, .pem, .key, .db, .db-shm, .db-wal
if git diff --cached --name-only | grep -E "\.env$|\.pem$|\.key$|\.db$|\.db-shm$|\.db-wal$"; then
    echo "❌ Secret/Artifact Leak: Sensitive file extension detected!"
    FAILED=1
fi

# 2. Check content (staged)
# We use git grep --cached to search the index.
# -I ignores binary files.
# -n prints line numbers.
# Patterns: Private Key, OpenAI Key, Google API Key
if git grep --cached -I -n -E "[B]EGIN PRIVATE KEY|sk-[a-zA-Z0-9]{20,}|[A]Iza[a-zA-Z0-9_-]{30,}"; then
    echo "❌ Secret Leak: Sensitive pattern detected in staged content!"
    FAILED=1
fi

if [ $FAILED -ne 0 ]; then
    exit 1
fi

echo "Quality Check Passed! 🚀"
