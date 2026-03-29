#!/bin/bash
# scripts/ci/check-coverage-drift.sh - Unified coverage threshold anti-degradation guard
# Supports local pre-commit (staged) and CI (PR/Push) workflows.

set -e

# Support direct argument for BASE to allow local testing or --local flag
MODE=$1
BASE_OVERRIDE=$2

# Paths
THRESHOLD_FILE=".agents/coverage-threshold"
UTILS_SCRIPT="scripts/utils.sh"

# Colors (fallback if utils.sh not available)
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

# Messaging Helper
if [ -f "$UTILS_SCRIPT" ]; then
    source "$UTILS_SCRIPT"
else
    info() { echo -e "${CYAN}  ℹ️  INFO:${NC} $1"; }
    fail() { echo -e "${RED}  ❌ FAIL:${NC} $1"; }
    pass() { echo -e "${GREEN}  ✅ PASS:${NC} $1"; }
fi

if [ ! -f "$THRESHOLD_FILE" ]; then
    info "No threshold file found at $THRESHOLD_FILE, skipping check"
    exit 0
fi

# Determine BASE and compare logic
if [ "$MODE" == "--local" ]; then
    # Local mode: Check if threshold file is staged
    if ! git diff --cached --name-only | grep -q "^$THRESHOLD_FILE$"; then
        info "Skipping: $THRESHOLD_FILE not modified in staged changes"
        exit 0
    fi
    BASE="HEAD"
    TARGET="staged"
    CURRENT_VAL=$(cat "$THRESHOLD_FILE")
    PREVIOUS_VAL=$(git show "HEAD:$THRESHOLD_FILE" 2>/dev/null || echo "0.0")
else
    # CI/Direct mode
    EVENT_NAME="${GITHUB_EVENT_NAME:-push}"
    BASE_REF="${GITHUB_BASE_REF:-main}"
    EVENT_BEFORE="${GITHUB_EVENT_BEFORE:-}"

    if [ -n "$BASE_OVERRIDE" ]; then
        BASE="$BASE_OVERRIDE"
    elif [ "$EVENT_NAME" == "pull_request" ]; then
        BASE="origin/$BASE_REF"
    else
        BASE="$EVENT_BEFORE"
        if [ -z "$BASE" ] || [ "$BASE" == "0000000000000000000000000000000000000000" ]; then
            BASE="HEAD^1"
        fi
    fi

    # Ensure BASE exists
    if ! git rev-parse --verify "$BASE" >/dev/null 2>&1; then
        info "Fetching $BASE..."
        git fetch origin "$BASE_REF" --depth=1 || true
    fi

    # Check if threshold modified in diff
    if ! git diff --name-only "$BASE" HEAD | grep -q "^$THRESHOLD_FILE$"; then
        info "Skipping check: $THRESHOLD_FILE not modified against $BASE"
        exit 0
    fi

    TARGET="HEAD"
    CURRENT_VAL=$(cat "$THRESHOLD_FILE")
    PREVIOUS_VAL=$(git show "$BASE:$THRESHOLD_FILE" 2>/dev/null || echo "0.0")
fi

info "Comparing threshold: CURRENT($CURRENT_VAL) vs PREVIOUS($PREVIOUS_VAL) [Mode: ${MODE:---direct}]"

# Perform Comparison
if awk "BEGIN {exit !($CURRENT_VAL < $PREVIOUS_VAL)}"; then
    fail "Coverage threshold cannot be lowered ($CURRENT_VAL < $PREVIOUS_VAL)"
    echo "Check your tests. If you must lower the threshold, it requires SystemsArchitect approval."
    exit 1
fi

pass "Coverage threshold anti-degradation check passed"
