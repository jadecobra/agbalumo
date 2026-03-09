#!/bin/bash
set -e

# Robust PATH discovery for macOS and Linux
for dir in /usr/local/bin /opt/homebrew/bin /usr/bin /bin; do
    case ":$PATH:" in
        *":$dir:"*) ;;
        *) export PATH="$PATH:$dir" ;;
    esac
done

# Verify golangci-lint config is valid (runs always, even with no staged files)
if command -v golangci-lint >/dev/null 2>&1; then
    CONFIG_FILE="scripts/.golangci.yml"
    if [ -f "$CONFIG_FILE" ]; then
        echo "Verifying golangci-lint config..."
        if ! golangci-lint config verify --config="$CONFIG_FILE" >/dev/null 2>&1; then
            echo "Error: golangci-lint config is invalid"
            golangci-lint config verify --config="$CONFIG_FILE" || true
            exit 1
        fi
        echo "✅ golangci-lint config is valid"
    fi
fi

START_TIME=$(date +%s)

START_TIME=$(date +%s)

# Colors
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
YELLOW=$(printf '\033[1;33m')
BLUE=$(printf '\033[1;34m')
BOLD=$(printf '\033[1m')
NC=$(printf '\033[0m')

echo "${BLUE}Running 10x Engineer Quality Checks...${NC}"

# 0. Workflow State Check
STATE_FILE=".agent/state.json"
if [ -f "$STATE_FILE" ]; then
    FEATURE=$(jq -r .feature "$STATE_FILE")
    if [ "$FEATURE" != "none" ] && [ "$FEATURE" != "null" ]; then
        PHASE=$(jq -r .phase "$STATE_FILE")
        echo "${BLUE}  Workflow detected: $FEATURE ($PHASE)${NC}"
        
        # Check mandatory gates for any committed work
        # For now, we enforce that 'lint' and 'red-test' (if RED+) must be PASS
        LINT_GATE=$(jq -r '.gates.lint' "$STATE_FILE")
        if [ "$LINT_GATE" != "PASS" ]; then
            echo "  ${RED}❌ Workflow Error: 'lint' gate must be PASS before committing.${NC}"
            exit 1
        fi
        
        if [ "$PHASE" != "IDLE" ] && [ "$PHASE" != "RED" ]; then
            RED_GATE=$(jq -r '.gates["red-test"]' "$STATE_FILE")
            if [ "$RED_GATE" != "PASS" ]; then
                 echo "  ${RED}❌ Workflow Error: 'red-test' gate must be PASS for $PHASE phase.${NC}"
                 exit 1
            fi
        fi
        echo "  ${GREEN}✅ Workflow gates verified${NC}"
    fi
fi

# 1. Get staged files for efficient checking
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACMR | grep '\.go$' || true)
MOD_FILES_CHANGED=$(git diff --cached --name-only --diff-filter=ACMR | grep -E 'go\.mod$|go\.sum$' || true)
STAGED_CMD_DOCS=$(git diff --cached --name-only --diff-filter=ACMR | grep -E '^cmd/|^docs/api\.md$|^docs/openapi\.yaml$' || true)
STAGED_CLI_DOCS=$(git diff --cached --name-only --diff-filter=ACMR | grep -E '^cmd/cli/|^docs/cli\.md$' || true)
STAGED_PERF_FILES=$(git diff --cached --name-only --diff-filter=ACMR | grep -E '^ui/|^internal/handler/|^internal/repository/|^scripts/' || true)
STAGED_SEC_FILES=$(git diff --cached --name-only --diff-filter=ACMR | grep -E '\.html$|\.go$|\.js$' || true)
STAGED_AGENT_FILES=$(git diff --cached --name-only --diff-filter=ACMR | grep -E '\.agents/agent\.yaml$|docs/CODING_STANDARDS\.md$' || true)
STAGED_BRAND_FILES=$(git diff --cached --name-only --diff-filter=ACMR | grep -E '\.agent/rules/brand\.toon$' || true)
STAGED_ALL=$(git diff --cached --name-only || true)

# Create a temporary directory for parallel task outputs
LOG_DIR=$(mktemp -d)
trap 'rm -rf "$LOG_DIR"' EXIT

run_task() {
    local task_id=$1
    local task_name=$2
    shift 2
    local start=$(date +%s)
    local log_file="$LOG_DIR/$task_id.log"
    if [ "$task_id" = "lint" ]; then
        mkdir -p @tester
        log_file="@tester/lint-results.txt"
    fi

    if "$@" > "$log_file" 2>&1; then
        local end=$(date +%s)
        echo "  ${GREEN}✅ $task_name passed ($((end - start))s)${NC}"
        return 0
    else
        local end=$(date +%s)
        echo "  ${RED}❌ $task_name failed ($((end - start))s)${NC}"
        cat "$log_file"
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
if [ -n "$STAGED_CMD_DOCS" ]; then
    run_task "api_drift" "API Drift" bash scripts/api-drift-check.sh &
else
    echo "  ${YELLOW}skipping API Drift (no relevant changes)${NC}"
fi

if [ -n "$STAGED_CLI_DOCS" ]; then
    run_task "cli_drift" "CLI Drift" bash scripts/cli-drift-check.sh &
else
    echo "  ${YELLOW}skipping CLI Drift (no relevant changes)${NC}"
fi

# 3.1 Agent Drift Check
if [ -n "$STAGED_AGENT_FILES" ]; then
    run_task "agent_drift" "Agent Drift" bash scripts/agent-drift-check.sh &
else
    echo "  ${YELLOW}skipping Agent Drift (no relevant changes)${NC}"
fi

# 4. Performance Audit
if [ -n "$STAGED_PERF_FILES" ]; then
    run_task "perf" "Performance Audit" sh scripts/performance-audit.sh &
else
    echo "  ${YELLOW}skipping Performance Audit (no relevant changes)${NC}"
fi

# 4.1 Brand Juice Generation
if [ -n "$STAGED_BRAND_FILES" ]; then
    run_task "brand" "Brand Juice" bash scripts/generate-juice.sh &
else
    echo "  ${YELLOW}skipping Brand Juice (no relevant changes)${NC}"
fi

# 5. Tests & Coverage
if [ -n "$STAGED_GO_FILES" ]; then
    check_tests() {
        # Removing -race locally to improve test speed in pre-commit
        go test -coverprofile=@tester/coverage.out ./... > /dev/null
        COVERAGE=$(go tool cover -func=@tester/coverage.out | grep total | grep -oE "[0-9]+(\.[0-9]+)?" | head -1)
        THRESHOLD=90.0
        if [ "$(echo "$COVERAGE < $THRESHOLD" | bc -l)" -eq 1 ]; then
            echo "Coverage is below threshold: $COVERAGE% < $THRESHOLD%"
            return 2
        fi
        echo "Coverage: $COVERAGE%"
    }
    run_task "test" "Tests & Coverage" check_tests &
else
    echo "  ${YELLOW}skipping Tests & Coverage (no staged Go files)${NC}"
fi

# 6. Security Check
if [ -n "$STAGED_ALL" ]; then
    run_task "security" "Security Check" sh scripts/security-check.sh &
else
    echo "  ${YELLOW}skipping Security Check (no staged files)${NC}"
fi

# 7. Check for ignored files being staged (Add/Modify only)
if [ -n "$STAGED_ALL" ]; then
    STAGED_NEW_OR_MOD=$(git diff --cached --name-only --diff-filter=ACMR || true)
    if [ -n "$STAGED_NEW_OR_MOD" ]; then
        IGNORED_STAGED=$(git check-ignore --stdin <<< "$STAGED_NEW_OR_MOD" || true)
        if [ -n "$IGNORED_STAGED" ]; then
            echo "  ${RED}❌ Error: The following ignored files are staged for commit:${NC}"
            echo "$IGNORED_STAGED" | sed 's/^/    /'
            echo "  ${YELLOW}Please unstage them with 'git restore --staged <file>' and run 'git rm --cached' if they should not be tracked.${NC}"
            exit 1
        fi
    fi
fi

# 8. CI Workflow Toolset Verification
run_task "ci_tools" "CI Toolset" bash scripts/verify-ci-tools.sh &

# 9. Local CI Verification
run_task "ci_local" "Local CI (act)" bash scripts/ci-local.sh --list &

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
