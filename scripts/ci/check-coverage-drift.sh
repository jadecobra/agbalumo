#!/bin/bash
# scripts/ci/check-coverage-drift.sh - Verify coverage threshold anti-degradation in CI

set -e

# Support direct argument for BASE to allow local testing
BASE_OVERRIDE=$1

# Default values from Environment Variables (set by GitHub Actions or local caller)
THRESHOLD_FILE="${THRESHOLD_FILE:-.agents/coverage-threshold}"
EVENT_NAME="${GITHUB_EVENT_NAME:-push}"
BASE_REF="${GITHUB_BASE_REF:-main}"
EVENT_BEFORE="${GITHUB_EVENT_BEFORE:-}"

# Colors (if scripts/utils.sh not available)
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'

# Messaging Helper (integrated with scripts/utils.sh if possible)
if [ -f "scripts/utils.sh" ]; then
    source scripts/utils.sh
else
    info() { echo -e "${CYAN}  ℹ️  INFO:${NC} $1"; }
    fail() { echo -e "${RED}  ❌ FAIL:${NC} $1"; }
    pass() { echo -e "${GREEN}  ✅ PASS:${NC} $1"; }
fi

if [ ! -f "$THRESHOLD_FILE" ]; then
    info "No threshold file found at $THRESHOLD_FILE, skipping anti-degradation check"
    exit 0
fi

# Determine base branch/commit
if [ -n "$BASE_OVERRIDE" ]; then
    BASE="$BASE_OVERRIDE"
elif [ "$EVENT_NAME" == "pull_request" ]; then
    # For PRs, compare with the merge base of the base branch
    BASE="origin/$BASE_REF"
else
    # For pushes, compare with the previous commit
    BASE="$EVENT_BEFORE"
    # Fallback if first push or no before SHA
    if [ -z "$BASE" ] || [ "$BASE" == "0000000000000000000000000000000000000000" ]; then
        BASE="HEAD^1"
    fi
fi

info "Comparing current with BASE: $BASE"

# Ensure we have the base commit fetched if it looks like a hash or remote ref
if ! git rev-parse --verify "$BASE" >/dev/null 2>&1; then
    info "Fetching $BASE..."
    git fetch origin "$BASE_REF" || true
fi

if git diff --name-only "$BASE" HEAD | grep -q "^$THRESHOLD_FILE$"; then
    OLD_THRESHOLD=$(git show "$BASE:$THRESHOLD_FILE" 2>/dev/null || echo "0.0")
    NEW_THRESHOLD=$(cat "$THRESHOLD_FILE")
    info "New Threshold: $NEW_THRESHOLD, Old Threshold: $OLD_THRESHOLD"
    
    if awk "BEGIN {exit !($NEW_THRESHOLD < $OLD_THRESHOLD)}"; then
        fail "Coverage threshold cannot be lowered ($NEW_THRESHOLD < $OLD_THRESHOLD)"
        echo "You must write tests to maintain or improve coverage, not lower the threshold!"
        exit 1
    fi
    pass "Coverage threshold anti-degradation check passed"
else
    info "Skipping check: $THRESHOLD_FILE not modified in this change"
fi
