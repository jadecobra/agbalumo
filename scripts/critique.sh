#!/bin/bash
# scripts/critique.sh
# ChiefCritic Robustness Audit Script
# Analyzes Go code for technical debt: Cognitive Complexity, Constants, and Struct Alignment.

set -e

# Configuration
THRESHOLD_COGNIT=10
PROJECT_ROOT=$(git rev-parse --show-toplevel)
GOBIN="${PROJECT_ROOT}/.tester/tmp/go/bin"
PATH="${GOBIN}:${PATH}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

TARGETS=${1:-"internal cmd"}

# Derived package targets for tools that expect packages/directories
PKG_TARGETS=""
for target in $TARGETS; do
    if [ -f "$target" ]; then
        PKG_TARGETS="$(dirname "$target") $PKG_TARGETS"
    else
        PKG_TARGETS="$target $PKG_TARGETS"
    fi
done
PKG_TARGETS=$(echo "$PKG_TARGETS" | tr ' ' '\n' | sort -u | tr '\n' ' ')

echo -e "${BLUE}${BOLD}--- ChiefCritic Robustness Audit ---${NC}"
echo -e "Target Paths: ${YELLOW}${TARGETS}${NC}"
echo -e "Package Targets: ${YELLOW}${PKG_TARGETS}${NC}"
echo -e "Cognitive Complexity Threshold: ${YELLOW}< ${THRESHOLD_COGNIT}${NC}"

REPORT_FILE="${PROJECT_ROOT}/critique_report.md"
echo "# ChiefCritic Technical Debt Report" > "$REPORT_FILE"
echo "Generated: $(date)" >> "$REPORT_FILE"
echo "" >> "$REPORT_FILE"

# 1. gocognit check
echo -e "\n${BLUE}[1/3] Checking Cognitive Complexity (gocognit)...${NC}"
echo "## 1. Cognitive Complexity (Threshold < $THRESHOLD_COGNIT)" >> "$REPORT_FILE"
echo '```' >> "$REPORT_FILE"

# gocognit works on files and directories
COGNIT_OUT=$(gocognit -over $((THRESHOLD_COGNIT-1)) $TARGETS 2>&1 || true)
if [ -n "$COGNIT_OUT" ]; then
    echo -e "${RED}Found functions exceeding complexity threshold:${NC}"
    echo "$COGNIT_OUT"
    echo "$COGNIT_OUT" >> "$REPORT_FILE"
else
    echo -e "${GREEN}✅ All functions are within complexity limits.${NC}"
    echo "No complexity issues found." >> "$REPORT_FILE"
fi
echo '```' >> "$REPORT_FILE"

# 2. goconst check
echo -e "\n${BLUE}[2/3] Checking for repeated strings (goconst)...${NC}"
echo "## 2. Repeated Constants" >> "$REPORT_FILE"
echo '```' >> "$REPORT_FILE"

# goconst requires directories/packages
CONST_OUT=$(goconst $PKG_TARGETS 2>&1 || true)
if [[ -n "$CONST_OUT" && "$CONST_OUT" != *"not a directory"* ]]; then
    echo -e "${YELLOW}Found repeated strings that should be constants:${NC}"
    echo "$CONST_OUT"
    echo "$CONST_OUT" >> "$REPORT_FILE"
else
    echo -e "${GREEN}✅ No repeated strings found.${NC}"
    echo "No constant issues found." >> "$REPORT_FILE"
fi
echo '```' >> "$REPORT_FILE"

# 3. maligned check
echo -e "\n${BLUE}[3/3] Checking struct alignment (maligned)...${NC}"
echo "## 3. Struct Alignment (maligned)" >> "$REPORT_FILE"
echo '```' >> "$REPORT_FILE"

# maligned requires packages
MALIGN_OUT=$(maligned $PKG_TARGETS 2>&1 || true)
if [[ -n "$MALIGN_OUT" && "$MALIGN_OUT" != *"not a directory"* && "$MALIGN_OUT" != *"internal error"* ]]; then
    echo -e "${YELLOW}Found structs that could be better aligned:${NC}"
    echo "$MALIGN_OUT"
    echo "$MALIGN_OUT" >> "$REPORT_FILE"
else
    echo -e "${GREEN}✅ All structs are optimally aligned.${NC}"
    echo "No alignment issues found." >> "$REPORT_FILE"
fi
echo '```' >> "$REPORT_FILE"

echo -e "\n${BLUE}${BOLD}--- Audit Complete ---${NC}"
echo -e "Detailed report saved to: ${YELLOW}${REPORT_FILE}${NC}"

# Exit with failure if gocognit failed (strict enforcement)
if [[ $COGNIT_OUT == *"over"* ]]; then
    exit 1
fi
