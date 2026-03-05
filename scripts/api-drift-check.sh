#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

printf "${BLUE}Running API Drift Check...${NC}\n"

# Helper for normalization
# 1. method path
# 2. replace :id with {id}
# 3. ensure leading slash if missing (should not happen here)
# 4. remove trailing slashes
# 5. if path is empty, make it /
# 6. deduplicate slashes
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

# Direct Echo routes: e.GET("/path", ...)
grep -E 'e\.(GET|POST|PUT|DELETE|PATCH)\("' cmd/server.go | \
    sed -E 's/.*e\.([A-Z]+)\("([^"]*)".*/\1 \2/' >> "$ROUTES_TMP"

# Groups in cmd/server.go
adminGroup_Path=$(grep "adminGroup := e.Group" cmd/server.go | sed -E 's/.*e\.Group\("([^"]+)".*/\1/')
adminLoginGroup_Path=$(grep "adminLoginGroup := adminGroup.Group" cmd/server.go | sed -E 's/.*adminGroup\.Group\("([^"]+)".*/\1/')

# Extract adminGroup routes: adminGroup.GET("/path", ...)
grep -E 'adminGroup\.(GET|POST|PUT|DELETE|PATCH)\("' cmd/server.go | \
    sed -E "s@.*adminGroup\.([A-Z]+)\(\"([^\"]*)\".*@\1 ${adminGroup_Path}\2@" >> "$ROUTES_TMP"

# Extract adminLoginGroup routes: adminLoginGroup.POST("", ...)
grep -E 'adminLoginGroup\.(GET|POST|PUT|DELETE|PATCH)\("' cmd/server.go | \
    sed -E "s@.*adminLoginGroup\.([A-Z]+)\(\"([^\"]*)\".*@\1 ${adminGroup_Path}${adminLoginGroup_Path}\2@" >> "$ROUTES_TMP"

IMPLEMENTED_ROUTES=$(normalize < "$ROUTES_TMP")
rm "$ROUTES_TMP"

# 2. Extract endpoints from docs/openapi.yaml
OPENAPI_ENDPOINTS=$(awk '
    /^  \/[^:]+:$/ { 
        current_path = substr($1, 1, length($1)-1)
    }
    /^    (get|post|put|delete|patch):$/ { 
        method = toupper(substr($1, 1, length($1)-1))
        print method " " current_path
    }
' docs/openapi.yaml | normalize)

# 3. Extract endpoints from docs/api.md
API_MD_ENDPOINTS=$(grep -E '^\| (GET|POST|PUT|DELETE|PATCH) \| `?/[^`| ]*`? \|' docs/api.md | \
    sed -E 's/\| ([A-Z]+) \| `?([^`| ]*)`?.*/\1 \2/' | normalize)

# Function to find differences
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
    # Endpoints in source but missing in target
    local missing=$(comm -23 "$tmp_source" "$tmp_target")
    
    if [ ! -z "$missing" ]; then
        local err_file=$(mktemp)
        echo "$missing" | while read -r line; do
            if [ ! -z "$line" ]; then
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

printf "\n${BLUE}Comparing Code vs OpenAPI Spec...${NC}\n"
check_diff "Code (cmd/server.go)" "OpenAPI (docs/openapi.yaml)" "$IMPLEMENTED_ROUTES" "$OPENAPI_ENDPOINTS" || DRIFT_DETECTED=1

printf "\n${BLUE}Comparing Code vs API Markdown...${NC}\n"
check_diff "Code (cmd/server.go)" "API Docs (docs/api.md)" "$IMPLEMENTED_ROUTES" "$API_MD_ENDPOINTS" || DRIFT_DETECTED=1

printf "\n${BLUE}Comparing OpenAPI vs API Markdown...${NC}\n"
check_diff "OpenAPI (docs/openapi.yaml)" "API Docs (docs/api.md)" "$OPENAPI_ENDPOINTS" "$API_MD_ENDPOINTS" || DRIFT_DETECTED=1

if [ $DRIFT_DETECTED -eq 1 ]; then
    printf "\n${RED}❌ API Drift Detected! Please update documentation or code to match.${NC}\n"
    exit 1
else
    printf "\n${GREEN}✅ All APIs are in sync across implementation and documentation.${NC}\n"
    exit 0
fi
