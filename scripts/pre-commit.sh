#!/bin/bash
set -e

# Robust PATH discovery for macOS and Linux
for dir in /usr/local/bin /opt/homebrew/bin /usr/bin /bin; do
    case ":$PATH:" in
        *":$dir:"*) ;;
        *) export PATH="$PATH:$dir" ;;
    esac
done

START_TIME=$(date +%s)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[1;34m'
NC='\033[0m'

echo "${BLUE}Running 10x Engineer Quality Checks...${NC}"

# 1. Get staged files for efficient checking
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACMR | grep '\.go$' || true)
MOD_FILES_CHANGED=$(git diff --cached --name-only --diff-filter=ACMR | grep -E 'go\.mod$|go\.sum$' || true)

# Create a temporary directory for parallel task outputs
LOG_DIR=$(mktemp -d)
trap 'rm -rf "$LOG_DIR"' EXIT

run_task() {
    local task_id=$1
    local task_name=$2
    shift 2
    local start=$(date +%s)
    echo "  [ ] Running $task_name..."
    if "$@" > "$LOG_DIR/$task_id.log" 2>&1; then
        local end=$(date +%s)
        echo "  ${GREEN}✅ $task_name passed ($((end - start))s)${NC}"
        return 0
    else
        local end=$(date +%s)
        echo "  ${RED}❌ $task_name failed ($((end - start))s)${NC}"
        cat "$LOG_DIR/$task_id.log"
        return 1
    fi
}

# 1. GolangCI-Lint (Smart local optimization)
if [ -n "$STAGED_GO_FILES" ]; then
    if command -v golangci-lint >/dev/null 2>&1; then
        # Use --new-from-rev=HEAD for extremely fast local linting of only changes
        run_task "lint" "GolangCI-Lint" golangci-lint run -c scripts/.golangci.yml --new-from-rev=HEAD &
    else
        # Fallback to standard tools if golangci-lint is not installed
        echo "  ${YELLOW}⚠️  golangci-lint not found, falling back to gofmt/govet${NC}"
        check_fmt() {
            UNFORMATTED=$(gofmt -l $STAGED_GO_FILES)
            if [ -n "$UNFORMATTED" ]; then
                echo "Go Code is not formatted. Run 'gofmt -w $STAGED_GO_FILES'"
                return 1
            fi
        }
        run_task "fmt" "Go Fmt" check_fmt &
        run_task "vet" "Go Vet" go vet ./... &
    fi
else
    echo "  ${YELLOW}skipping Lint/Fmt (no staged Go files)${NC}"
fi

# 2. Go Mod Tidy Check (Only if mod files changed)
if [ -n "$MOD_FILES_CHANGED" ]; then
    check_mod() {
        go mod tidy
        if ! git diff --exit-code --quiet go.mod go.sum; then
            echo "go.mod/go.sum are not tidy. Run 'go mod tidy' and commit changes."
            return 1
        fi
    }
    run_task "mod" "Go Mod Tidy" check_mod &
else
    echo "  ${YELLOW}skipping Go Mod Tidy (no changes to go.mod/go.sum)${NC}"
fi

# 3. API & CLI Drift Checks
run_task "api_drift" "API Drift" bash scripts/api-drift-check.sh &
run_task "cli_drift" "CLI Drift" bash scripts/cli-drift-check.sh &

# 4. Performance Audit
run_task "perf" "Performance Audit" sh scripts/performance-audit.sh &

# 5. Tests & Coverage
check_tests() {
    go test -race -coverprofile=@tester/coverage.out ./... > /dev/null
    COVERAGE=$(go tool cover -func=@tester/coverage.out | grep total | grep -oE "[0-9]+(\.[0-9]+)?" | head -1)
    THRESHOLD=90.0
    if [ "$(echo "$COVERAGE < $THRESHOLD" | bc -l)" -eq 1 ]; then
        echo "Coverage is below threshold: $COVERAGE% < $THRESHOLD%"
        return 2
    fi
    echo "Coverage: $COVERAGE%"
}
run_task "test" "Tests & Coverage" check_tests &

# Wait for all background tasks
FAILURES=0
for job in $(jobs -p); do
    wait $job || FAILURES=$((FAILURES + 1))
done

END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))

if [ $FAILURES -eq 0 ]; then
    echo ""
    echo "${GREEN}${BOLD}Quality Check Passed in ${TOTAL_TIME}s! 🚀${NC}"
    exit 0
else
    echo ""
    echo "${RED}${BOLD}Quality Check Failed! (failures in $FAILURES tasks). Fix issues before committing.${NC}"
    exit 1
fi
