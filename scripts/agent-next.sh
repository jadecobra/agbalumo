#!/bin/bash
# scripts/agent-next.sh: Unified wrapper for agent workflow automation.
# Automatically infers the current phase and runs relevant gates.
#
# Usage:
#   ./scripts/agent-next.sh [pattern]  - Run next pending gate(s) for current phase

set -e

# Robust PATH discovery and utilities
source "$(dirname "$0")/utils.sh"
setup_path

STATE_FILE=".agent/state.json"

if [ ! -f "$STATE_FILE" ]; then
    echo "Error: State file not found. Initialize with ./scripts/agent-exec.sh workflow init <feature>"
    exit 1
fi

FEATURE=$(jq -r .feature "$STATE_FILE")
PHASE=$(jq -r .phase "$STATE_FILE")
WORKFLOW_TYPE=$(jq -r '.workflow_type // "feature"' "$STATE_FILE")

if [ "$FEATURE" == "null" ] || [ "$FEATURE" == "none" ] || [ -z "$FEATURE" ]; then
    echo "No active feature. Use './scripts/agent-exec.sh workflow init <feature>' to start."
    exit 0
fi

info "Current Feature: $FEATURE [$WORKFLOW_TYPE] ($PHASE)"

case "$PHASE" in
    IDLE)
        warn "Phase is IDLE. Transitioning to RED..."
        ./scripts/agent-exec.sh workflow set-phase RED
        # Reload state and run next
        exec "$0" "$@"
        ;;
    RED)
        RED_TEST=$(jq -r '.gates["red-test"]' "$STATE_FILE")
        API_SPEC=$(jq -r '.gates["api-spec"]' "$STATE_FILE")
        
        if [ "$RED_TEST" != "PASS" ]; then
            info "Gate 'red-test' is pending..."
            ./scripts/agent-gate.sh red-test "$@"
        fi
        
        # Reload state to check if phase changed or red-test passed
        RED_TEST=$(jq -r '.gates["red-test"]' "$STATE_FILE")
        if [ "$API_SPEC" != "PASS" ] && [ "$RED_TEST" == "PASS" ]; then
            info "Gate 'api-spec' is pending..."
            ./scripts/agent-gate.sh api-spec
        fi
        ;;
    GREEN)
        IMPL=$(jq -r '.gates["implementation"]' "$STATE_FILE")
        if [ "$IMPL" != "PASS" ]; then
            info "Gate 'implementation' is pending..."
            ./scripts/agent-gate.sh implementation
        fi
        ;;
    REFACTOR)
        LINT=$(jq -r '.gates["lint"]' "$STATE_FILE")
        COVERAGE=$(jq -r '.gates["coverage"]' "$STATE_FILE")
        BROWSER=$(jq -r '.gates["browser-verification"]' "$STATE_FILE")
        
        if [ "$LINT" != "PASS" ]; then
            info "Gate 'lint' is pending..."
            ./scripts/agent-gate.sh lint
        fi
        
        if [ "$COVERAGE" != "PASS" ]; then
            info "Gate 'coverage' is pending..."
            ./scripts/agent-gate.sh coverage
        fi
        
        # Final status check for browser
        LINT=$(jq -r '.gates["lint"]' "$STATE_FILE")
        COVERAGE=$(jq -r '.gates["coverage"]' "$STATE_FILE")
        if [ "$LINT" == "PASS" ] && [ "$COVERAGE" == "PASS" ]; then
            if [ "$BROWSER" != "PASS" ]; then
                warn "All automated REFACTOR gates passed."
                info "Remaining: browser-verification (PENDING)"
                info "Please run browser subagent or manually pass it: ./scripts/agent-exec.sh workflow gate browser-verification PASS"
            else
                pass "All gates for feature '$FEATURE' are PASS!"
                info "Suggested: Run ./scripts/agent-exec.sh workflow init none to reset for next task."
            fi
        fi
        ;;
    *)
        echo "Error: Unknown phase '$PHASE'"
        exit 1
        ;;
esac
