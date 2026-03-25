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
pass() { if [ "$FMT" != "json" ]; then echo "${GREEN}  ✅ PASS:${NC} $1"; fi }
warn() { if [ "$FMT" != "json" ]; then echo "${YELLOW}  ⚠️  WARN:${NC} $1"; fi }
fail() { if [ "$FMT" != "json" ]; then echo "${RED}  ❌ FAIL:${NC} $1"; fi }
info() { if [ "$FMT" != "json" ]; then echo "${CYAN}  ℹ️  INFO:${NC} $1"; fi }

# output_json_envelope <success_bool> <command_string> <output_string_or_json> [warnings_json_array]
output_json_envelope() {
    local success=$1
    local cmd=$2
    local out=$3
    local warnings=${4:-"[]"}

    # If $out is valid JSON, insert it directly, otherwise as array of strings
    local out_json
    if echo "$out" | jq -e . >/dev/null 2>&1; then
        out_json="$out"
    else
        out_json=$(jq -Rn --arg str "$out" '$str')
    fi

    # If $warnings is valid JSON, insert it directly, otherwise wrap in an array
    local warnings_json
    if echo "$warnings" | jq -e . >/dev/null 2>&1; then
        warnings_json="$warnings"
    else
        warnings_json=$(jq -Rn --arg str "$warnings" '[$str]')
    fi

    # Convert success to explicit boolean 
    local bool_success="true"
    if [ "$success" = "false" ] || [ "$success" = "0" ]; then
        bool_success="false"
    fi

    jq -n \
        --argjson success "$bool_success" \
        --arg cmd "$cmd" \
        --argjson output "$out_json" \
        --argjson warnings "$warnings_json" \
        '{success: $success, command: $cmd, output: $output, warnings: $warnings}'
}

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
        mkdir -p /tmp/.tester/coverage
        log_file="/tmp/.tester/coverage/lint-results.txt"
    fi

    if "$@" > "$log_file" 2>&1; then
        local end=$(date +%s)
        if [ "$FMT" != "json" ]; then echo "  ${GREEN}✅ $task_name passed ($((end - start))s)${NC}"; fi
        return 0
    else
        local end=$(date +%s)
        if [ "$FMT" != "json" ]; then 
            echo "  ${RED}❌ $task_name failed ($((end - start))s)${NC}"
            cat "$log_file"
        fi
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
        if [ "$FMT" != "json" ]; then
            echo "  ${RED}❌ Workflow gate enforcement failed for '$feature' [$workflow_type] ($phase):${NC}"
            echo "  ${RED}   Required gates not PASS:${failed_gates}${NC}"
            
            local doc_link=".agents/workflows/feature-implementation.md"
            if [ "$workflow_type" = "bugfix" ]; then doc_link=".agents/workflows/bugfix.md"; fi
            if [ "$workflow_type" = "refactor" ]; then doc_link=".agents/workflows/refactor.md"; fi

            echo "  ${YELLOW}   See: $doc_link${NC}"
        fi
        return 1
    fi

    if [ "$FMT" != "json" ]; then echo "  ${GREEN}✅ Workflow gates verified ($phase: $(echo $required_gates | wc -w | tr -d ' ') gates)${NC}"; fi
    return 0
}
