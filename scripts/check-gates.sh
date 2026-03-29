#!/bin/bash
# scripts/check-gates.sh - Workflow gate enforcement
STATE_FILE=".agents/state.json"
[ ! -f "$STATE_FILE" ] && exit 0
FEATURE=$(jq -r .feature "$STATE_FILE")
[ "$FEATURE" = "none" ] || [ "$FEATURE" = "null" ] || [ -z "$FEATURE" ] && exit 0
PHASE=$(jq -r .phase "$STATE_FILE")
WORKFLOW=$(jq -r '.workflow_type // "feature"' "$STATE_FILE")
case "$PHASE" in
    RED) G="red-test" ;;
    GREEN) G="red-test api-spec implementation" ;;
    REFACTOR) G="red-test api-spec implementation lint coverage" ;;
    IDLE) G="red-test api-spec implementation lint coverage browser-verification" ;;
    *) exit 0 ;;
esac
FAILURES=0; F_G=""
for g in $G; do
    S=$(jq -r ".gates[\"$g\"]" "$STATE_FILE")
    if [ "$S" != "PASS" ] && [ "$S" != "PASSED" ]; then
        FAILURES=$((FAILURES + 1)); F_G="$F_G $g($S)"
    fi
done
if [ "$FAILURES" -gt 0 ]; then
    echo "  ❌ Workflow gate enforcement failed for '$FEATURE' [$WORKFLOW] ($PHASE):"
    echo "  ❌ Required gates not PASS:$F_G"
    exit 1
fi
echo "  ✅ Workflow gates verified ($PHASE: $(echo $G | wc -w | tr -d ' ') gates)"
