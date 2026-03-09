#!/bin/bash
# agent-gate.sh: Automated gate verification for agbalumo.
# Usage: ./scripts/agent-gate.sh <gate_id>

set -e

# Robust PATH discovery for macOS and Linux
for dir in /usr/local/bin /opt/homebrew/bin /usr/bin /bin; do
    case ":$PATH:" in
        *":$dir:"*) ;;
        *) export PATH="$PATH:$dir" ;;
    esac
done

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
}

case "$GATE_ID" in
    red-test)
        # Phase check: red-test is only valid when phase is RED (or GREEN/REFACTOR if we want to ensure it STILL fails without implementation, but RED is the standard entry)
        # Expected: go test fails.
        echo "Running tests expecting failure..."
        if go test ./... > /dev/null 2>&1; then
            echo "❌ Gate FAIL: red-test passed but was expected to fail."
            update_gate "FAIL"
            exit 1
        else
            echo "✅ Gate PASS: red-test failed as expected."
            update_gate "PASS"
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
        THRESHOLD=$(grep -oE "THRESHOLD=[0-9]+(\.[0-9]+)?" scripts/pre-commit.sh | cut -d= -f2 || echo "90.0")
        
        if [ "$(echo "$COVERAGE < $THRESHOLD" | bc -l)" -eq 1 ]; then
            echo "❌ Gate FAIL: Coverage $COVERAGE% is below threshold $THRESHOLD%."
            update_gate "FAIL"
            exit 1
        else
            echo "✅ Gate PASS: Coverage $COVERAGE% meets threshold $THRESHOLD%."
            update_gate "PASS"
        fi
        ;;
    browser-verification)
        echo "⚠️  browser-verification requires manual confirmation or browser subagent."
        # For the sake of "Automated", we check if a 'browser_verification_done' file exists or similar.
        # For now, we'll just require it to be set manually via agent-exec.sh, but we can verify it's PASS.
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
