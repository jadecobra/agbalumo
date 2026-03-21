# Check for format flag
FMT="json"
if [ "$1" = "--text" ]; then 
    FMT="text"
    shift
fi

# Robust PATH discovery
source "$(dirname "$0")/utils.sh"
setup_path

# Documentation Links
DOC_WORKFLOW=".agent/workflows/feature-implementation.md"
DOC_STANDARDS="docs/CODING_STANDARDS.md"
DOC_API="docs/api.md"
DOC_CLI="docs/cli.md"

START_TIME=$(date +%s)
if [ "$FMT" != "json" ]; then echo "${BLUE}Running 10x Engineer Quality Checks...${NC}"; fi

# 0. Workflow Gate Enforcement (phase-aware)
STATE_FILE=".agent/state.json"
if [ -f "$STATE_FILE" ]; then
    if ! check_workflow_gates "$STATE_FILE"; then
        if [ "$FMT" = "json" ]; then output_json_envelope false "pre-commit.sh" "Workflow gate enforcement failed."; fi
        exit 1
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

# 1. GolangCI-Lint (Smart local optimization)
if [ -n "$STAGED_GO_FILES" ]; then
    if command -v golangci-lint >/dev/null 2>&1; then
        run_task "lint" "GolangCI-Lint" "$LOG_DIR" golangci-lint run -c scripts/.golangci.yml --new-from-rev=HEAD &
    else
        if [ "$FMT" != "json" ]; then echo "  ${YELLOW}⚠️  golangci-lint not found, falling back to gofmt/govet${NC}"; fi
        check_fmt() {
            UNFORMATTED=$(gofmt -l $STAGED_GO_FILES)
            if [ -n "$UNFORMATTED" ]; then
                if [ "$FMT" != "json" ]; then
                    echo "Go Code is not formatted. Run 'gofmt -w $STAGED_GO_FILES'"
                    echo "See: $DOC_STANDARDS (Section 3)"
                fi
                return 1
            fi
        }
        run_task "fmt" "Go Fmt" "$LOG_DIR" check_fmt &
        run_task "vet" "Go Vet" "$LOG_DIR" go vet ./... &
    fi
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping Lint/Fmt (no staged Go files)${NC}"; fi
fi

# 2. Go Mod Tidy Check (Only if mod files changed)
if [ -n "$MOD_FILES_CHANGED" ]; then
    check_mod() {
        go mod tidy
        if ! git diff --exit-code --quiet go.mod go.sum; then
            if [ "$FMT" != "json" ]; then
                echo "go.mod/go.sum are not tidy. Run 'go mod tidy' and commit changes."
                echo "See: $DOC_STANDARDS"
            fi
            return 1
        fi
    }
    run_task "mod" "Go Mod Tidy" "$LOG_DIR" check_mod &
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping Go Mod Tidy (no changes to go.mod/go.sum)${NC}"; fi
fi

# 3. API & CLI Drift Checks
if [ -n "$STAGED_CMD_DOCS" ]; then
    run_task "api_drift" "API Drift" "$LOG_DIR" bash scripts/api-drift-check.sh &
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping API Drift (no relevant changes)${NC}"; fi
fi

if [ -n "$STAGED_CLI_DOCS" ]; then
    run_task "cli_drift" "CLI Drift" "$LOG_DIR" bash scripts/cli-drift-check.sh &
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping CLI Drift (no relevant changes)${NC}"; fi
fi

# 3.1 Agent Drift Check
if [ -n "$STAGED_AGENT_FILES" ]; then
    run_task "agent_drift" "Agent Drift" "$LOG_DIR" bash scripts/agent-drift-check.sh &
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping Agent Drift (no relevant changes)${NC}"; fi
fi

# 3.2 Template Function Drift Check
if [ -n "$STAGED_SEC_FILES" ]; then
    run_task "template_drift" "Template Drift" "$LOG_DIR" bash scripts/template-func-drift-check.sh &
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping Template Draft (no relevant changes)${NC}"; fi
fi

# 4. Performance Audit
if [ -n "$STAGED_PERF_FILES" ]; then
    run_task "perf" "Performance Audit" "$LOG_DIR" sh scripts/performance-audit.sh &
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping Performance Audit (no relevant changes)${NC}"; fi
fi

# 4.1 Brand Juice Generation
if [ -n "$STAGED_BRAND_FILES" ]; then
    run_task "brand" "Brand Juice" "$LOG_DIR" bash scripts/generate-juice.sh &
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping Brand Juice (no relevant changes)${NC}"; fi
fi

# 5. Tests & Coverage
if [ -n "$STAGED_GO_FILES" ]; then
    check_tests() {
        mkdir -p .tester/coverage
        go test -json -race -count=1 -coverprofile=.tester/coverage/coverage.out ./... > /dev/null
        COVERAGE=$(go tool cover -func=.tester/coverage/coverage.out | awk '/^total:/ {print substr($3, 1, length($3)-1)}')
        THRESHOLD_FILE=".agent/coverage-threshold"
        THRESHOLD=90.0
        if [ -f "$THRESHOLD_FILE" ]; then
            THRESHOLD=$(cat "$THRESHOLD_FILE")
        fi
        if awk "BEGIN {exit !($COVERAGE < $THRESHOLD)}"; then
            if [ "$FMT" != "json" ]; then
                echo "  ${RED}❌ Error: Coverage is below threshold: $COVERAGE% < $THRESHOLD%${NC}"
                echo "  ${YELLOW}Top 5 lowest coverage files:${NC}"
                go tool cover -func=.tester/coverage/coverage.out | grep -v "100.0%" | sort -k 3 -n | head -5 | sed 's/^/    /'
                echo "  ${BLUE}See: $DOC_WORKFLOW (Gate: coverage)${NC}"
            fi
            return 2
        fi
        if [ "$FMT" != "json" ]; then echo "Coverage: $COVERAGE%"; fi
    }
    run_task "test" "Tests & Coverage" "$LOG_DIR" check_tests &
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping Tests & Coverage (no staged Go files)${NC}"; fi
fi

# 5.5 Coverage Threshold Anti-Degradation Check
if [ -n "$STAGED_ALL" ]; then
    check_threshold() {
        THRESHOLD_FILE=".agent/coverage-threshold"
        if [ -f "$THRESHOLD_FILE" ] && git ls-files --error-unmatch "$THRESHOLD_FILE" >/dev/null 2>&1; then
            if git diff --cached --name-only | grep -q "^$THRESHOLD_FILE$"; then
                OLD_THRESHOLD=$(git show HEAD:$THRESHOLD_FILE 2>/dev/null || echo "0.0")
                NEW_THRESHOLD=$(cat "$THRESHOLD_FILE")
                if awk "BEGIN {exit !($NEW_THRESHOLD < $OLD_THRESHOLD)}"; then
                    if [ "$FMT" != "json" ]; then
                        echo "  ${RED}❌ Error: Coverage threshold cannot be lowered ($NEW_THRESHOLD < $OLD_THRESHOLD)${NC}"
                        echo "  ${YELLOW}You must write tests to maintain or improve coverage, not lower the threshold!${NC}"
                    fi
                    return 1
                fi
            fi
        fi
        return 0
    }
    run_task "threshold" "Coverage Threshold Check" "$LOG_DIR" check_threshold &
fi

# 6. Security Check
if [ -n "$STAGED_ALL" ]; then
    run_task "security" "Security Check" "$LOG_DIR" sh scripts/security-check.sh &
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}skipping Security Check (no staged files)${NC}"; fi
fi

# 7. Check for ignored files being staged (Add/Modify only)
if [ -n "$STAGED_ALL" ]; then
    STAGED_NEW_OR_MOD=$(git diff --cached --name-only --diff-filter=ACMR || true)
    if [ -n "$STAGED_NEW_OR_MOD" ]; then
        IGNORED_STAGED=$(git check-ignore --stdin <<< "$STAGED_NEW_OR_MOD" || true)
        if [ -n "$IGNORED_STAGED" ]; then
            if [ "$FMT" != "json" ]; then
                echo "  ${RED}❌ Error: The following ignored files are staged for commit:${NC}"
                echo "$IGNORED_STAGED" | sed 's/^/    /'
                echo "  ${YELLOW}Please unstage them with 'git restore --staged <file>' and run 'git rm --cached' if they should not be tracked.${NC}"
                echo "  ${YELLOW}See: $DOC_STANDARDS${NC}"
            else
                output_json_envelope false "pre-commit.sh" "Ignored files staged for commit: $IGNORED_STAGED"
            fi
            exit 1
        fi
    fi
fi

# 8. CI Workflow Toolset Verification
run_task "ci_tools" "CI Toolset" "$LOG_DIR" bash scripts/verify-ci-tools.sh &

# 9. Local CI Verification
run_task "ci_local" "Local CI (act)" "$LOG_DIR" bash scripts/ci-local.sh --list &

# 10. Restart Verification
SKIP_PRE_COMMIT=true run_task "restart" "Verify Restart" "$LOG_DIR" bash scripts/verify_restart.sh &

# Wait for all background tasks
FAILURES=0
for job in $(jobs -p); do
    wait $job || FAILURES=$((FAILURES + 1))
done

END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))

if [ "$FMT" = "json" ]; then
    if [ $FAILURES -eq 0 ]; then
        output_json_envelope true "pre-commit.sh" "Quality Check Passed in ${TOTAL_TIME}s!"
    else
        ERRORS=""
        if [ -d "$LOG_DIR" ]; then
            for log in "$LOG_DIR"/*.log; do
                [ -e "$log" ] || continue
                # Identify failures by checking if file is non-empty and maybe checking grep
                if [ -s "$log" ]; then
                    ERRORS="$ERRORS\n--- $(basename "$log") ---\n$(cat "$log")"
                fi
            done
        fi
        output_json_envelope false "pre-commit.sh" "Quality Check Failed! ($FAILURES failures)" "[\"$ERRORS\"]"
    fi
    exit $FAILURES
fi

if [ $FAILURES -eq 0 ]; then
    echo ""
    echo "${GREEN}${BOLD}Quality Check Passed in ${TOTAL_TIME}s! 🚀${NC}"
    exit 0
else
    echo ""
    echo "${RED}${BOLD}Quality Check Failed! (failures in $FAILURES tasks). Fix issues before committing.${NC}"
    echo "${BLUE}Refer to the following for standards and workflows:${NC}"
    echo "  - ${BLUE}$DOC_STANDARDS${NC}"
    echo "  - ${BLUE}$DOC_WORKFLOW${NC}"
    exit 1
fi
