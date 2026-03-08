#!/bin/bash

# agent-exec.sh: Multi-agent persona execution helper for agbalumo.
# Usage: ./scripts/agent-exec.sh role <persona_name>

set -e

PERSONA_DIR=".agent/personas"
GLOBAL_RULES="$PERSONA_DIR/Global.md"

function show_usage() {
    echo "Usage: $0 role <persona_name>"
    echo "Available personas:"
    ls "$PERSONA_DIR" | grep -v "Global.md" | sed 's/\.md//'
    exit 1
}

if [ "$1" != "role" ] || [ -z "$2" ]; then
    show_usage
fi

ROLE=$2
PERSONA_FILE="$PERSONA_DIR/$ROLE.md"

if [ ! -f "$PERSONA_FILE" ]; then
    echo "Error: Persona '$ROLE' not found at $PERSONA_FILE"
    show_usage
fi

echo "<activated_persona name=\"$ROLE\">"
cat "$GLOBAL_RULES"
echo ""
cat "$PERSONA_FILE"
echo "</activated_persona>"
