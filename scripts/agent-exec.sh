#!/bin/bash
# agent-exec.sh: Multi-agent persona execution helper and workflow manager for agbalumo.
# Usage: 
#   ./scripts/agent-exec.sh role <persona_name>
#   ./scripts/agent-exec.sh workflow <subcommand> [args]

set -e

PERSONA_DIR=".agent/personas"
GLOBAL_RULES="$PERSONA_DIR/Global.md"
STATE_FILE=".agent/state.json"

function show_usage() {
    echo "Usage:"
    echo "  $0 role <persona_name>"
    echo "  $0 workflow init <feature_name>"
    echo "  $0 workflow set-persona <persona_name>"
    echo "  $0 workflow set-phase <IDLE|RED|GREEN|REFACTOR>"
    echo "  $0 workflow gate <gate_id> <PENDING|PASS|FAIL>"
    echo "  $0 workflow verify <gate_id>"
    echo "  $0 workflow status"
    echo ""
    echo "Available personas:"
    ls "$PERSONA_DIR" | grep -v "Global.md" | sed 's/\.md//'
    exit 1
}

function handle_workflow() {
    local cmd=$1
    shift

    if [ ! -f "$STATE_FILE" ]; then
        echo "{}" > "$STATE_FILE"
    fi

    case "$cmd" in
        init)
            local feature=$1
            if [ -z "$feature" ]; then echo "Error: feature name required"; exit 1; fi
            jq -n --arg f "$feature" --arg t "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
                '{feature: $f, persona: "none", phase: "IDLE", gates: { "red-test": "PENDING", "api-spec": "PENDING", "implementation": "PENDING", "lint": "PENDING", "coverage": "PENDING", "browser-verification": "PENDING"}, updated_at: $t}' \
                > "$STATE_FILE"
            echo "Workflow initialized for feature: $feature"
            ;;
        set-persona)
            local persona=$1
            if [ -z "$persona" ]; then echo "Error: persona name required"; exit 1; fi
            jq --arg p "$persona" --arg t "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
                '.persona = $p | .updated_at = $t' "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
            echo "Persona set to: $persona"
            ;;
        set-phase)
            local phase=$1
            if [[ ! "$phase" =~ ^(IDLE|RED|GREEN|REFACTOR)$ ]]; then echo "Error: invalid phase '$phase'"; exit 1; fi
            jq --arg p "$phase" --arg t "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
                '.phase = $p | .updated_at = $t' "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
            echo "Phase set to: $phase"
            ;;
        gate)
            local gate=$1
            local status=$2
            if [ -z "$gate" ] || [ -z "$status" ]; then echo "Usage: workflow gate <gate_id> <status>"; exit 1; fi
            if [[ ! "$status" =~ ^(PENDING|PASS|FAIL)$ ]]; then echo "Error: invalid status '$status'"; exit 1; fi
            jq --arg g "$gate" --arg s "$status" --arg t "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
                '.gates[$g] = $s | .updated_at = $t' "$STATE_FILE" > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
            echo "Gate '$gate' set to: $status"
            ;;
        verify)
            local gate=$1
            if [ -z "$gate" ]; then echo "Usage: workflow verify <gate_id>"; exit 1; fi
            bash scripts/agent-gate.sh "$gate"
            ;;
        status)
            jq . "$STATE_FILE"
            ;;
        *)
            show_usage
            ;;
    esac
}

if [ "$1" == "role" ]; then
    if [ -z "$2" ]; then show_usage; fi
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
elif [ "$1" == "workflow" ]; then
    handle_workflow "$2" "$3" "$4"
else
    show_usage
fi
