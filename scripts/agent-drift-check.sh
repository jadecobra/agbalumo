#!/bin/bash
set -e

# Colors for output
RED=$(printf '\033[0;31m')
GREEN=$(printf '\033[0;32m')
YELLOW=$(printf '\033[1;33m')
BLUE=$(printf '\033[1;34m')
NC=$(printf '\033[0m')

printf "${BLUE}Running Agent Drift Check (10x Standard)...${NC}\n"

AGENT_YAML=".agents/agent.yaml"
CODING_STANDARDS="docs/CODING_STANDARDS.md"

# 1. Enforce "Double-Commit" rule if staged
YAML_STAGED=$(git diff --cached --name-only | grep -E "^$(basename $(dirname $AGENT_YAML))/$(basename $AGENT_YAML)$" || true)
MD_STAGED=$(git diff --cached --name-only | grep -E "^$CODING_STANDARDS$" || true)

if [ -n "$YAML_STAGED" ] && [ -z "$MD_STAGED" ]; then
    printf "${RED}❌ Double-Commit Rule Violated!${NC}\n"
    printf "Changes to ${YELLOW}$AGENT_YAML${NC} must be mirrored in ${YELLOW}$CODING_STANDARDS${NC} (Section 5).\n"
    printf "Please stage both files together.\n"
    exit 1
fi

# 2. Verify persona sync (regardless of staging, for local checks)
# Extract agent names from YAML
YAMl_AGENTS=$(grep "  - name: " "$AGENT_YAML" | awk '{print $NF}' | sort)

# Extract agent names from Markdown (Section 5)
MD_AGENTS=$(sed -n '/## 5. Agent Protocol/,/## 6/p' "$CODING_STANDARDS" | grep -E '^\*   \*\*[^:]+\*\*:' | sed -E 's/\*   \*\*([^:*]+)\*\*:.*/\1/' | sort)

if [ "$YAMl_AGENTS" != "$MD_AGENTS" ]; then
    printf "${RED}❌ Agent Drift Detected!${NC}\n"
    printf "Personas in ${YELLOW}$AGENT_YAML${NC} do not match ${YELLOW}$CODING_STANDARDS${NC} (Section 5).\n\n"
    
    printf "${BLUE}YAML Agents:${NC}\n$YAMl_AGENTS\n\n"
    printf "${BLUE}Markdown Agents:${NC}\n$MD_AGENTS\n\n"
    
    printf "Please ensure the lists match exactly.\n"
    exit 1
fi

printf "${GREEN}✅ Agent personas are in sync and Double-Commit rule honored.${NC}\n"
exit 0
