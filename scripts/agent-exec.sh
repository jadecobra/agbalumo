#!/bin/bash
# agent-exec.sh: workflow manager for agbalumo.
# Usage: 
#   ./scripts/agent-exec.sh workflow <subcommand> [args]

set -e

STATE_FILE=".agent/state.json"

function show_usage() {
    echo "Usage:"
    echo "  $0 workflow init <feature_name> [workflow_type]"
    echo "  $0 workflow set-phase <IDLE|RED|GREEN|REFACTOR>"
    echo "  $0 workflow gate <gate_id> <PENDING|PASS|FAIL>"
    echo "  $0 workflow verify <gate_id>"
    echo "  $0 workflow status"
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
            local workflow_type=${2:-feature}
            if [ -z "$feature" ]; then echo "Error: feature name required"; exit 1; fi
            if [[ ! "$workflow_type" =~ ^(feature|bugfix|refactor)$ ]]; then echo "Error: invalid workflow type '$workflow_type'"; exit 1; fi
            
            jq -n --arg f "$feature" --arg wt "$workflow_type" --arg t "$(date -u +"%Y-%m-%dT%H:%M:%SZ")" \
                '{feature: $f, workflow_type: $wt, phase: "IDLE", gates: { "red-test": "PENDING", "api-spec": "PENDING", "implementation": "PENDING", "lint": "PENDING", "coverage": "PENDING", "browser-verification": "PENDING"}, updated_at: $t}' \
                > "$STATE_FILE"
            echo "Workflow initialized for $workflow_type: $feature"
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

if [ "$1" == "workflow" ]; then
    handle_workflow "$2" "$3" "$4"
else
    show_usage
fi

