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
