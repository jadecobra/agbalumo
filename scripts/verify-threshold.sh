#!/bin/bash
# scripts/verify-threshold.sh - Verify coverage threshold anti-degradation
THRESHOLD_FILE=".agents/coverage-threshold"
if [ -f "$THRESHOLD_FILE" ] && git ls-files --error-unmatch "$THRESHOLD_FILE" >/dev/null 2>&1; then
    if git diff --cached --name-only | grep -q "^$THRESHOLD_FILE$"; then
        OLD_THRESHOLD=$(git show HEAD:$THRESHOLD_FILE 2>/dev/null || echo "0.0")
        NEW_THRESHOLD=$(cat "$THRESHOLD_FILE")
        if awk "BEGIN {exit !($NEW_THRESHOLD < $OLD_THRESHOLD)}"; then
            echo "Error: Coverage threshold cannot be lowered ($NEW_THRESHOLD < $OLD_THRESHOLD)"
            exit 1
        fi
    fi
fi
