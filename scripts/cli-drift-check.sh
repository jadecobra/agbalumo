#!/bin/bash
set -e

# Colors for output
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
BLUE=$(printf '\033[0;34m')
NC=$(printf '\033[0m')

printf "${BLUE}Running CLI Drift Check...${NC}\n"

# 1. Extract commands from cmd/*.go
# We look for Use: "..." lines and take the first word.
IMPLEMENTED_COMMANDS=$(grep -rh 'Use:[[:space:]]*"[^"]*"' cmd/ | \
    sed -E 's/.*Use:[[:space:]]*"([^ "]+).*/\1/' | \
    grep -v "^$" | \
    sort -u)

# 2. Extract command headers from docs/cli.md
# We look for ### (main commands) or ##### (subcommands)
# We exclude "subcommands", "commands", "flags", "example", etc.
CLI_MD_COMMANDS=$(grep -E '^###+ ' docs/cli.md | \
    sed -E 's/^###+ //g' | \
    tr '[:upper:]' '[:lower:]' | \
    grep -vE '^(subcommands|flags|example|quick reference|environment variables)$' | \
    sort -u)

# Function to find differences
check_diff() {
    local source_name=$1
    local target_name=$2
    local source_content="$3"
    local target_content="$4"
    
    local tmp_source=$(mktemp)
    local tmp_target=$(mktemp)
    echo "$source_content" > "$tmp_source"
    echo "$target_content" > "$tmp_target"

    local error_found=0
    local missing=$(comm -23 "$tmp_source" "$tmp_target")
    
    if [ ! -z "$missing" ]; then
        local err_file=$(mktemp)
        echo "$missing" | while read -r line; do
            if [ ! -z "$line" ] && [[ "$line" != "agbalumo" ]]; then
                printf "${RED}❌ Missing in %s: %s (found in %s)${NC}\n" "$target_name" "$line" "$source_name"
                echo "1" > "$err_file"
            fi
        done
        if [ -s "$err_file" ]; then error_found=1; fi
        rm "$err_file"
    fi
    
    rm "$tmp_source" "$tmp_target"
    return $error_found
}

DRIFT_DETECTED=0

printf "\n${BLUE}Comparing Code vs CLI Documentation...${NC}\n"
check_diff "Code (cmd/*.go)" "CLI Docs (docs/cli.md)" "$IMPLEMENTED_COMMANDS" "$CLI_MD_COMMANDS" || DRIFT_DETECTED=1

if [ $DRIFT_DETECTED -eq 1 ]; then
    printf "\n${RED}❌ CLI Drift Detected! Please update documentation or code to match.${NC}\n"
    exit 1
else
    printf "\n${GREEN}✅ CLI Documentation is in sync with implementation.${NC}\n"
    exit 0
fi
