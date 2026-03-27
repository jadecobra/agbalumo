#!/bin/bash
# scripts/utils.sh - Common shell utilities

# Colors
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
YELLOW=$(printf '\033[1;33m')
BLUE=$(printf '\033[1;34m')
CYAN=$(printf '\033[0;36m')
BOLD=$(printf '\033[1m')
NC=$(printf '\033[0m')

# Path Discovery
setup_path() {
    for dir in /usr/local/bin /opt/homebrew/bin /usr/bin /bin; do
        case ":$PATH:" in
            *":$dir:"*) ;;
            *) export PATH="$PATH:$dir" ;;
        esac
    done
}

# Messaging
pass() { if [ "$FMT" != "json" ]; then echo "${GREEN}  ✅ PASS:${NC} $1"; fi }
warn() { if [ "$FMT" != "json" ]; then echo "${YELLOW}  ⚠️  WARN:${NC} $1"; fi }
fail() { if [ "$FMT" != "json" ]; then echo "${RED}  ❌ FAIL:${NC} $1"; fi }
info() { if [ "$FMT" != "json" ]; then echo "${CYAN}  ℹ️  INFO:${NC} $1"; fi }

# output_json_envelope <success_bool> <command_string> <output_string_or_json> [warnings_json_array]
output_json_envelope() {
    local success=$1
    local cmd=$2
    local out=$3
    local warnings=${4:-"[]"}

    # If $out is valid JSON, insert it directly, otherwise as array of strings
    local out_json
    if echo "$out" | jq -e . >/dev/null 2>&1; then
        out_json="$out"
    else
        out_json=$(jq -Rn --arg str "$out" '$str')
    fi

    # If $warnings is valid JSON, insert it directly, otherwise wrap in an array
    local warnings_json
    if echo "$warnings" | jq -e . >/dev/null 2>&1; then
        warnings_json="$warnings"
    else
        warnings_json=$(jq -Rn --arg str "$warnings" '[$str]')
    fi

    # Convert success to explicit boolean 
    local bool_success="true"
    if [ "$success" = "false" ] || [ "$success" = "0" ]; then
        bool_success="false"
    fi

    jq -n \
        --argjson success "$bool_success" \
        --arg cmd "$cmd" \
        --argjson output "$out_json" \
        --argjson warnings "$warnings_json" \
        '{success: $success, command: $cmd, output: $output, warnings: $warnings}'
}
