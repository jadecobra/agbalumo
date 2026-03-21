#!/bin/bash
# test_gate_enforcement.sh: Tests for phase-aware gate enforcement in utils.sh
# Usage: bash scripts/test_gate_enforcement.sh

set -e

# Check for format flag
FMT="json"
if [ "$1" = "--text" ]; then 
    FMT="text"
    shift
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/utils.sh"

PASS_COUNT=0
FAIL_COUNT=0
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

COLLECTED_WARNINGS=()

assert_pass() {
    local test_name=$1
    local state_file=$2
    if check_workflow_gates "$state_file" >/dev/null 2>&1; then
        PASS_COUNT=$((PASS_COUNT + 1))
        if [ "$FMT" != "json" ]; then echo "  ${GREEN}✅ PASS: $test_name${NC}"; fi
    else
        FAIL_COUNT=$((FAIL_COUNT + 1))
        if [ "$FMT" != "json" ]; then echo "  ${RED}❌ FAIL: $test_name (expected PASS, got FAIL)${NC}"; fi
        COLLECTED_WARNINGS+=("FAIL: $test_name (expected PASS, got FAIL)")
    fi
}

assert_fail() {
    local test_name=$1
    local state_file=$2
    if check_workflow_gates "$state_file" >/dev/null 2>&1; then
        FAIL_COUNT=$((FAIL_COUNT + 1))
        if [ "$FMT" != "json" ]; then echo "  ${RED}❌ FAIL: $test_name (expected FAIL, got PASS)${NC}"; fi
        COLLECTED_WARNINGS+=("FAIL: $test_name (expected FAIL, got PASS)")
    else
        PASS_COUNT=$((PASS_COUNT + 1))
        if [ "$FMT" != "json" ]; then echo "  ${GREEN}✅ PASS: $test_name${NC}"; fi
    fi
}

make_state() {
    local file="$TMP_DIR/$1.json"
    local feature=$2
    local phase=$3
    local red_test=$4
    local api_spec=$5
    local impl=$6
    local lint=$7
    local coverage=$8
    local browser=$9
    cat > "$file" <<EOF
{
  "feature": "$feature",
  "persona": "none",
  "phase": "$phase",
  "gates": {
    "red-test": "$red_test",
    "api-spec": "$api_spec",
    "implementation": "$impl",
    "lint": "$lint",
    "coverage": "$coverage",
    "browser-verification": "$browser"
  },
  "updated_at": "2026-01-01T00:00:00Z"
}
EOF
    echo "$file"
}

if [ "$FMT" != "json" ]; then
    echo "${BOLD}Running Gate Enforcement Tests${NC}"
    echo "================================"
    echo ""
    echo "${BLUE}1. No active feature${NC}"
fi

STATE=$(make_state "no_feature" "none" "IDLE" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING")
assert_pass "feature=none skips gate checks" "$STATE"

STATE=$(make_state "null_feature" "null" "IDLE" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING")
assert_pass "feature=null skips gate checks" "$STATE"

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BLUE}2. RED phase${NC}"
fi
STATE=$(make_state "red_pass" "search" "RED" "PASS" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING")
assert_pass "RED phase with red-test=PASS" "$STATE"

STATE=$(make_state "red_fail" "search" "RED" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING")
assert_fail "RED phase with red-test=PENDING" "$STATE"

STATE=$(make_state "red_fail2" "search" "RED" "FAIL" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING")
assert_fail "RED phase with red-test=FAIL" "$STATE"

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BLUE}3. GREEN phase${NC}"
fi
STATE=$(make_state "green_pass" "search" "GREEN" "PASS" "PASS" "PASS" "PENDING" "PENDING" "PENDING")
assert_pass "GREEN phase with all required gates PASS" "$STATE"

STATE=$(make_state "green_fail_impl" "search" "GREEN" "PASS" "PASS" "PENDING" "PENDING" "PENDING" "PENDING")
assert_fail "GREEN phase with implementation=PENDING" "$STATE"

STATE=$(make_state "green_fail_api" "search" "GREEN" "PASS" "PENDING" "PASS" "PENDING" "PENDING" "PENDING")
assert_fail "GREEN phase with api-spec=PENDING" "$STATE"

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BLUE}4. REFACTOR phase${NC}"
fi
STATE=$(make_state "refactor_pass" "search" "REFACTOR" "PASS" "PASS" "PASS" "PASS" "PASS" "PENDING")
assert_pass "REFACTOR phase with all required gates PASS" "$STATE"

STATE=$(make_state "refactor_fail_lint" "search" "REFACTOR" "PASS" "PASS" "PASS" "PENDING" "PASS" "PENDING")
assert_fail "REFACTOR phase with lint=PENDING" "$STATE"

STATE=$(make_state "refactor_fail_cov" "search" "REFACTOR" "PASS" "PASS" "PASS" "PASS" "PENDING" "PENDING")
assert_fail "REFACTOR phase with coverage=PENDING" "$STATE"

if [ "$FMT" != "json" ]; then
    echo ""
    echo "${BLUE}5. IDLE phase with active feature${NC}"
fi
STATE=$(make_state "idle_all_pass" "search" "IDLE" "PASS" "PASS" "PASS" "PASS" "PASS" "PASS")
assert_pass "IDLE+feature with all gates PASS" "$STATE"

STATE=$(make_state "idle_fail_browser" "search" "IDLE" "PASS" "PASS" "PASS" "PASS" "PASS" "PENDING")
assert_fail "IDLE+feature with browser-verification=PENDING" "$STATE"

STATE=$(make_state "idle_all_pending" "search" "IDLE" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING" "PENDING")
assert_fail "IDLE+feature with all gates PENDING" "$STATE"

TOTAL=$((PASS_COUNT + FAIL_COUNT))

if [ "$FMT" = "json" ]; then
    COMBINED_HINTS="[]"
    if [ ${#COLLECTED_WARNINGS[@]} -gt 0 ]; then
        COMBINED_HINTS=$(printf '%s\n' "${COLLECTED_WARNINGS[@]}" | jq -R . | jq -s .)
    fi

    if [ "$FAIL_COUNT" -eq 0 ]; then
        output_json_envelope true "test_gate_enforcement.sh" "All $TOTAL tests passed!" "$COMBINED_HINTS"
        exit 0
    else
        output_json_envelope false "test_gate_enforcement.sh" "$FAIL_COUNT of $TOTAL tests failed." "$COMBINED_HINTS"
        exit 1
    fi
fi

echo ""
echo "================================"
if [ "$FAIL_COUNT" -eq 0 ]; then
    echo "${GREEN}${BOLD}All $TOTAL tests passed! ✅${NC}"
    exit 0
else
    echo "${RED}${BOLD}$FAIL_COUNT of $TOTAL tests failed ❌${NC}"
    exit 1
fi
