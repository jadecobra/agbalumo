#!/bin/bash
# scripts/gitleaks-scan.sh
set -e
source "$(dirname "$0")/utils.sh"
setup_path

# Prioritize workspace-local gitleaks binary from .tester/tmp/go/bin
LOCAL_GITLEAKS="./.tester/tmp/go/bin/gitleaks"
if [ -f "$LOCAL_GITLEAKS" ]; then
    GITLEAKS_CMD="$LOCAL_GITLEAKS"
else
    echo "  ${RED}❌ Error: gitleaks binary not found at $LOCAL_GITLEAKS. Run 'task gitleaks:install' first.${NC}"
    exit 1
fi

"$GITLEAKS_CMD" protect --staged --verbose --redact
