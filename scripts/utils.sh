#!/bin/bash
# scripts/utils.sh - Common shell utilities

# Colors
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
YELLOW=$(printf '\033[1;33m')
BLUE=$(printf '\033[1;34m')
CYAN=$(printf '\033[0;36m')
BOLD=$(printf '\033[1m')
NC=$(printf '\033[0m')

# Path Discovery
setup_path() {
    for dir in /usr/local/bin /opt/homebrew/bin /usr/bin /bin; do
        case ":$PATH:" in
            *":$dir:"*) ;;
            *) export PATH="$PATH:$dir" ;;
        esac
    done
}

# Messaging
pass() { echo "${GREEN}  ✅ PASS:${NC} $1"; }
warn() { echo "${YELLOW}  ⚠️  WARN:${NC} $1"; }
fail() { echo "${RED}  ❌ FAIL:${NC} $1"; }
info() { echo "${CYAN}  ℹ️  INFO:${NC} $1"; }

# Task Runner
# Usage: run_task "task_id" "Task Name" log_dir command args...
run_task() {
    local task_id=$1
    local task_name=$2
    local log_dir=$3
    shift 3
    local start=$(date +%s)
    local log_file="$log_dir/$task_id.log"
    
    # Special legacy path for lint
    if [ "$task_id" = "lint" ]; then
        mkdir -p .tester/coverage
        log_file=".tester/coverage/lint-results.txt"
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

# Gate Enforcement
# Usage: check_workflow_gates <state_file>
# Returns 0 if all phase-required gates are PASS, 1 otherwise.
# Skips check entirely when feature is "none" or "null".
check_workflow_gates() {
    local state_file=$1

    if [ ! -f "$state_file" ]; then
        return 0
    fi

    local feature
    feature=$(jq -r .feature "$state_file")
    if [ "$feature" = "none" ] || [ "$feature" = "null" ] || [ -z "$feature" ]; then
        return 0
    fi

    local phase
    local workflow_type
    phase=$(jq -r .phase "$state_file")
    workflow_type=$(jq -r '.workflow_type // "feature"' "$state_file")

    local required_gates=""
    case "$phase" in
        RED)
            required_gates="red-test"
            ;;
        GREEN)
            required_gates="red-test api-spec implementation"
            ;;
        REFACTOR)
            required_gates="red-test api-spec implementation lint coverage"
            ;;
        IDLE)
            required_gates="red-test api-spec implementation lint coverage browser-verification"
            ;;
        *)
            return 0
            ;;
    esac

    local failures=0
    local failed_gates=""
    for gate in $required_gates; do
        local status
        status=$(jq -r ".gates[\"$gate\"]" "$state_file")
        if [ "$status" != "PASS" ] && [ "$status" != "PASSED" ]; then
            failures=$((failures + 1))
            failed_gates="$failed_gates $gate($status)"
        fi
    done

    if [ "$failures" -gt 0 ]; then
        echo "  ${RED}❌ Workflow gate enforcement failed for '$feature' [$workflow_type] ($phase):${NC}"
        echo "  ${RED}   Required gates not PASS:${failed_gates}${NC}"
        
        local doc_link=".agent/workflows/feature-implementation.md"
        if [ "$workflow_type" = "bugfix" ]; then doc_link=".agent/workflows/bugfix.md"; fi
        if [ "$workflow_type" = "refactor" ]; then doc_link=".agent/workflows/refactor.md"; fi

        echo "  ${YELLOW}   See: $doc_link${NC}"
        return 1
    fi

    echo "  ${GREEN}✅ Workflow gates verified ($phase: $(echo $required_gates | wc -w | tr -d ' ') gates)${NC}"
    return 0
}
