#!/usr/bin/env bash
# ChiefCritic "Anti-Programmer-Art" Gate
set -e

PLAN_FILE="${1:-implementation_plan.md}"

echo "[ChiefCritic] Auditing implementation plan: $PLAN_FILE..."

if [ ! -f "$PLAN_FILE" ]; then
    echo "❌ [ERROR] Implementation plan file not found: $PLAN_FILE"
    exit 1
fi

MANDATORY_HEADERS=(
    "Target User Avatar"
    "Pain Point Mapping"
    "Strategic Critique"
    "Technical Contract"
    "Security STRIDE"
)

for header in "${MANDATORY_HEADERS[@]}"; do
    if ! grep -qi "$header" "$PLAN_FILE"; then
        echo "❌ [ERROR] Missing mandatory section: '$header'"
        exit 1
    fi
done

# Depth Analysis: Strategic Critique should be more than a one-liner
CRITIQUE_LINES=$(sed -n '/Strategic Critique/,/##/p' "$PLAN_FILE" | grep -v "Strategic Critique" | grep -v "##" | sed '/^[[:space:]]*$/d' | wc -l)

if [ "$CRITIQUE_LINES" -lt 3 ]; then
    echo "❌ [REJECTED] Strategic Critique is insufficient. Push back harder!"
    exit 1
fi

# Depth Analysis: Technical Contract should contain code-ish bits (interfaces/schemas)
if ! grep -q "\`" "$PLAN_FILE"; then
    echo "❌ [REJECTED] Technical Contract lacks specific code references or schemas. Use backticks for contracts."
    exit 1
fi

echo "✅ [ChiefCritic] Plan approved. Proceeding to execution phase."
