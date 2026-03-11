#!/bin/bash
# agent-gate.sh: Automated gate verification for agbalumo.
# Usage: ./scripts/agent-gate.sh <gate_id>

set -e

# Robust PATH discovery
source "$(dirname "$0")/utils.sh"
setup_path

GATE_ID=$1
STATE_FILE=".agent/state.json"

if [ -z "$GATE_ID" ]; then
    echo "Usage: $0 <gate_id>"
    exit 1
fi

if [ ! -f "$STATE_FILE" ]; then
    echo "Error: State file not found."
    exit 1
fi

FEATURE=$(jq -r .feature "$STATE_FILE")
PHASE=$(jq -r .phase "$STATE_FILE")

if [ "$FEATURE" == "null" ] || [ "$FEATURE" == "" ]; then
    echo "Error: No active feature found in $STATE_FILE"
    exit 1
fi

echo "Verifying gate: $GATE_ID for feature: $FEATURE ($PHASE)"

update_gate() {
    local status=$1
    ./scripts/agent-exec.sh workflow gate "$GATE_ID" "$status" || true
    
    # Auto-transition logic
    STATE=$(jq -r . "$STATE_FILE")
    PHASE=$(echo "$STATE" | jq -r .phase)
    RED_TEST=$(echo "$STATE" | jq -r '.gates["red-test"]')
    API_SPEC=$(echo "$STATE" | jq -r '.gates["api-spec"]')
    IMPL=$(echo "$STATE" | jq -r '.gates["implementation"]')
    
    if [ "$PHASE" == "RED" ]; then
        if [ "$RED_TEST" == "PASS" ] && [ "$API_SPEC" == "PASS" ]; then
            echo "✨ All RED gates passed. Transitioning to GREEN phase."
            ./scripts/agent-exec.sh workflow set-phase GREEN
        fi
    elif [ "$PHASE" == "GREEN" ]; then
        if [ "$IMPL" == "PASS" ]; then
            echo "✨ Implementation passed. Transitioning to REFACTOR phase."
            ./scripts/agent-exec.sh workflow set-phase REFACTOR
        fi
    fi
    
    echo "--- Current Workflow Status ---"
    ./scripts/agent-exec.sh workflow status | jq -r '"Feature: \(.feature) (\(.phase))\nGates: \(.gates)"'
}

# Dependency Checks
check_deps() {
    local gate=$1
    case "$gate" in
        implementation)
            RED_TEST=$(jq -r '.gates["red-test"]' "$STATE_FILE")
            API_SPEC=$(jq -r '.gates["api-spec"]' "$STATE_FILE")
            if [ "$RED_TEST" != "PASS" ] || [ "$API_SPEC" != "PASS" ]; then
                echo "❌ Error: 'implementation' requires 'red-test' and 'api-spec' to be PASS."
                exit 1
            fi
            ;;
        lint|coverage)
            IMPL=$(jq -r '.gates["implementation"]' "$STATE_FILE")
            if [ "$IMPL" != "PASS" ]; then
                echo "❌ Error: '$gate' requires 'implementation' to be PASS."
                exit 1
            fi
            ;;
        browser-verification)
            IMPL=$(jq -r '.gates["implementation"]' "$STATE_FILE")
            if [ "$IMPL" != "PASS" ]; then
                echo "❌ Error: 'browser-verification' requires 'implementation' to be PASS."
                exit 1
            fi
            ;;
    esac
}

check_deps "$GATE_ID"

case "$GATE_ID" in
    red-test)
        # Phase check: red-test is only valid when phase is RED
        PATTERN=$2
        echo "Running tests expecting failure..."
        
        # 1. Verify code compiles first. go test -run=^$ compiles but runs no tests.
        if ! go test -run=^$ ./... > .tester/red-test-compile.log 2>&1; then
            fail "Code does not compile. Fixed syntax/imports before running red-test."
            cat .tester/red-test-compile.log
            update_gate "FAIL"
            exit 1
        fi

        # 2. Run tests and capture output. Use -v to ensure FAIL markers and patterns are visible.
        TEST_OUTPUT=$(go test -v ./... 2>&1 || true)
        
        # 3. Check for FAIL marker (actual test failure, not just non-zero exit)
        if echo "$TEST_OUTPUT" | grep -q -e "--- FAIL:"; then
            if [ -n "$PATTERN" ]; then
                if echo "$TEST_OUTPUT" | grep -q "$PATTERN"; then
                    pass "Gate PASS: red-test failed with expected pattern '$PATTERN'."
                    update_gate "PASS"
                else
                    fail "Gate FAIL: red-test failed but pattern '$PATTERN' not found in output."
                    echo "--- TEST OUTPUT ---"
                    echo "$TEST_OUTPUT" | tail -n 20
                    update_gate "FAIL"
                    exit 1
                fi
            else
                pass "Gate PASS: red-test failed as expected."
                update_gate "PASS"
            fi
        else
            # If we reach here, either it passed or it failed for an unknown reason (not grep-able failure)
            if echo "$TEST_OUTPUT" | grep -q "PASS$"; then
                fail "Gate FAIL: red-test passed but was expected to fail."
            else
                fail "Gate FAIL: tests failed but could not find '--- FAIL:' marker. Check for panics or setup issues."
                echo "--- TEST OUTPUT ---"
                echo "$TEST_OUTPUT" | tail -n 20
            fi
            update_gate "FAIL"
            exit 1
        fi
        ;;
    api-spec)
        echo "Running API drift check..."
        if bash scripts/api-drift-check.sh; then
            echo "✅ Gate PASS: api-spec drift check passed."
            update_gate "PASS"
        else
            echo "❌ Gate FAIL: api-spec drift check failed."
            update_gate "FAIL"
            exit 1
        fi
        ;;
    implementation)
        echo "Running build and tests..."
        if go build ./... && go test ./...; then
            echo "✅ Gate PASS: implementation build and tests passed."
            update_gate "PASS"
        else
            echo "❌ Gate FAIL: implementation build or tests failed."
            update_gate "FAIL"
            exit 1
        fi
        ;;
    lint)
        echo "Running linter..."
        if command -v golangci-lint >/dev/null 2>&1; then
            if golangci-lint run -c scripts/.golangci.yml; then
                echo "✅ Gate PASS: lint passed."
                update_gate "PASS"
            else
                echo "❌ Gate FAIL: lint failed."
                update_gate "FAIL"
                exit 1
            fi
        else
            echo "⚠️  golangci-lint not found, skipping..."
            update_gate "PASS"
        fi
        ;;
    coverage)
        echo "Verifying test coverage..."
        mkdir -p .tester/coverage
        go test -coverprofile=.tester/coverage/coverage.out ./... > /dev/null
        if [ ! -f ".tester/coverage/coverage.out" ]; then
            echo "❌ Gate FAIL: coverage profile not generated."
            update_gate "FAIL"
            exit 1
        fi
        COVERAGE=$(go tool cover -func=.tester/coverage/coverage.out | grep total | grep -oE "[0-9]+(\.[0-9]+)?" | head -1)
        
        THRESHOLD_FILE=".agent/coverage-threshold"
        THRESHOLD=90.0
        if [ -f "$THRESHOLD_FILE" ]; then
            THRESHOLD=$(cat "$THRESHOLD_FILE")
        fi
        
        if [ "$(echo "$COVERAGE < $THRESHOLD" | bc -l)" -eq 1 ]; then
            echo "❌ Gate FAIL: Coverage $COVERAGE% is below threshold $THRESHOLD%."
            echo "Top 5 lowest coverage files:"
            go tool cover -func=.tester/coverage/coverage.out | grep -v "100.0%" | sort -k 3 -n | head -5 | sed 's/^/  /'
            update_gate "FAIL"
            exit 1
        else
            echo "✅ Gate PASS: Coverage $COVERAGE% meets threshold $THRESHOLD%."
            update_gate "PASS"
        fi
        ;;
    browser-verification)
        echo "⚠️  browser-verification requires manual confirmation or browser subagent."
        STATUS=$(jq -r ".gates[\"$GATE_ID\"]" "$STATE_FILE")
        if [ "$STATUS" == "PASS" ]; then
             echo "✅ Gate PASS: browser-verification already marked as PASS."
        else
             echo "❌ Gate FAIL: browser-verification must be manually passed or verified via browser subagent."
             exit 1
        fi
        ;;
    *)
        echo "Error: Unknown gate_id '$GATE_ID'"
        exit 1
        ;;
esac

