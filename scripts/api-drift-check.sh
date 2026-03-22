#!/bin/bash
set -e

# Colors for output
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
BLUE=$(printf '\033[0;34m')
NC=$(printf '\033[0m')

# Check for format flag
FMT="json"
if [ "$1" = "--text" ]; then FMT="text"; fi

source "$(dirname "$0")/utils.sh"

if [ "$FMT" = "text" ]; then
    printf "${BLUE}Running API Drift Check...${NC}\n"
fi

# Helper for normalization
normalize() {
    sed 's/:id/{id}/g' | \
    sed -E 's|//+|/|g' | \
    sed -E 's|/$||' | \
    awk '{ 
        if ($2 == "") $2 = "/";
        print $1 " " $2 
    }' | \
    sort -u
}

# 1. Extract routes from cmd/server.go
ROUTES_TMP=$(mktemp)

# Direct Echo routes
grep -E 'e\.(GET|POST|PUT|DELETE|PATCH)\("' cmd/server.go | \
    sed -E 's/.*e\.([A-Z]+)\("([^"]*)".*/\1 \2/' >> "$ROUTES_TMP"

# Groups in cmd/server.go
adminGroup_Path=$(grep "adminGroup := e.Group" cmd/server.go | sed -E 's/.*e\.Group\("([^"]+)".*/\1/')
adminLoginGroup_Path=$(grep "adminLoginGroup := adminGroup.Group" cmd/server.go | sed -E 's/.*adminGroup\.Group\("([^"]+)".*/\1/')

# Extract adminGroup routes
grep -E 'adminGroup\.(GET|POST|PUT|DELETE|PATCH)\("' cmd/server.go | \
    sed -E "s@.*adminGroup\.([A-Z]+)\(\"([^\"]*)\".*@\1 ${adminGroup_Path}\2@" >> "$ROUTES_TMP"

# Extract adminLoginGroup routes
grep -E 'adminLoginGroup\.(GET|POST|PUT|DELETE|PATCH)\("' cmd/server.go | \
    sed -E "s@.*adminLoginGroup\.([A-Z]+)\(\"([^\"]*)\".*@\1 ${adminGroup_Path}${adminLoginGroup_Path}\2@" >> "$ROUTES_TMP"

IMPLEMENTED_ROUTES=$(normalize < "$ROUTES_TMP")
rm "$ROUTES_TMP"

# 2. Extract endpoints from docs/openapi.yaml
OPENAPI_ENDPOINTS=$(npx swagger-cli bundle docs/openapi.yaml -r -t yaml | awk '
    /^  \x27?\/[^:]+\x27?:$/ { 
        current_path = $1
        sub(/:$/, "", current_path)
        sub(/^\x27/, "", current_path)
        sub(/\x27$/, "", current_path)
    }
    /^    (get|post|put|delete|patch):$/ { 
        method = toupper(substr($1, 1, length($1)-1))
        print method " " current_path
    }
' | normalize)

# 3. Extract endpoints from docs/api.md
API_MD_ENDPOINTS=$(grep -E '^\| (GET|POST|PUT|DELETE|PATCH) \| `?/[^`| ]*`? \|' docs/api.md | \
    sed -E 's/\| ([A-Z]+) \| `?([^`| ]*)`?.*/\1 \2/' | normalize)

DRIFT_WARNINGS="[]"
COLLECTED_WARNINGS=()

check_diff() {
    local source_name=$1
    local target_name=$2
    local source_content="$3"
    local target_content="$4"
    
    local tmp_source=$(mktemp)
    local tmp_target=$(mktemp)
    echo "$source_content" | grep -v "^$" > "$tmp_source"
    echo "$target_content" | grep -v "^$" > "$tmp_target"

    local error_found=0
    local missing=$(comm -23 "$tmp_source" "$tmp_target")
    
    if [ ! -z "$missing" ]; then
        while IFS= read -r line; do
            if [ ! -z "$line" ]; then
                local warn_msg="Missing in $target_name: $line (found in $source_name)"
                COLLECTED_WARNINGS+=("$warn_msg")
                if [ "$FMT" = "text" ]; then
                    printf "${RED}❌ %s${NC}\n" "$warn_msg"
                fi
            fi
        done <<< "$missing"
        error_found=1
    fi
    
    rm "$tmp_source" "$tmp_target"
    return $error_found
}

DRIFT_DETECTED=0

if [ "$FMT" = "text" ]; then printf "\n${BLUE}Comparing Code vs OpenAPI Spec...${NC}\n"; fi
check_diff "Code (cmd/server.go)" "OpenAPI (docs/openapi.yaml)" "$IMPLEMENTED_ROUTES" "$OPENAPI_ENDPOINTS" || DRIFT_DETECTED=1

if [ "$FMT" = "text" ]; then printf "\n${BLUE}Comparing Code vs API Markdown...${NC}\n"; fi
check_diff "Code (cmd/server.go)" "API Docs (docs/api.md)" "$IMPLEMENTED_ROUTES" "$API_MD_ENDPOINTS" || DRIFT_DETECTED=1

if [ "$FMT" = "text" ]; then printf "\n${BLUE}Comparing OpenAPI vs API Markdown...${NC}\n"; fi
check_diff "OpenAPI (docs/openapi.yaml)" "API Docs (docs/api.md)" "$OPENAPI_ENDPOINTS" "$API_MD_ENDPOINTS" || DRIFT_DETECTED=1

# Convert bash array to JSON array string for our envelope
if [ ${#COLLECTED_WARNINGS[@]} -gt 0 ]; then
    DRIFT_WARNINGS=$(printf '%s\n' "${COLLECTED_WARNINGS[@]}" | jq -R . | jq -s .)
fi

if [ $DRIFT_DETECTED -eq 1 ]; then
    if [ "$FMT" = "text" ]; then
        printf "\n${RED}❌ API Drift Detected! Please update documentation or code to match.${NC}\n"
    else
        output_json_envelope false "api-drift-check.sh" "API Drift Detected! Please update documentation or code to match." "$DRIFT_WARNINGS"
    fi
    exit 1
else
    if [ "$FMT" = "text" ]; then
        printf "\n${GREEN}✅ All APIs are in sync across implementation and documentation.${NC}\n"
    else
        output_json_envelope true "api-drift-check.sh" "All APIs are in sync across implementation and documentation." "[]"
    fi
    exit 0
fi
