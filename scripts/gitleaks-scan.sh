#!/bin/bash
# scripts/gitleaks-scan.sh
set -e
source "$(dirname "$0")/utils.sh"
setup_path

# Prioritize workspace-local gitleaks binary from .tester/tmp/go/bin
LOCAL_GITLEAKS="./.tester/tmp/go/bin/gitleaks"
if [ -f "$LOCAL_GITLEAKS" ]; then
    GITLEAKS_CMD="$LOCAL_GITLEAKS"
elif command -v gitleaks >/dev/null 2>&1; then
    GITLEAKS_CMD="gitleaks"
else
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}Installing gitleaks...${NC}"; fi
    # ... expensive fallback install logic ...
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew >/dev/null 2>&1; then
            brew install gitleaks > /dev/null 2>&1
            GITLEAKS_CMD="gitleaks"
        else
            VERSION="8.18.2"
            ARCH="x64"
            if [[ "$(uname -m)" == "arm64" ]]; then ARCH="arm64"; fi
            curl -sSL "https://github.com/gitleaks/gitleaks/releases/download/v${VERSION}/gitleaks_${VERSION}_darwin_${ARCH}.tar.gz" | tar -xz -C /tmp gitleaks
            GITLEAKS_CMD="/tmp/gitleaks"
        fi
    else
        VERSION="8.18.2"
        ARCH="x64"
        if [[ "$(uname -m)" == "aarch64" || "$(uname -m)" == "arm64" ]]; then ARCH="arm64"; fi
        curl -sSL "https://github.com/gitleaks/gitleaks/releases/download/v${VERSION}/gitleaks_${VERSION}_linux_${ARCH}.tar.gz" | tar -xz -C /tmp gitleaks
        GITLEAKS_CMD="/tmp/gitleaks"
    fi
fi

if [ -z "$GITLEAKS_CMD" ] || ! "$GITLEAKS_CMD" version >/dev/null 2>&1; then
     echo "  ${RED}❌ Error: gitleaks is required but could not be initialized.${NC}"
     exit 1
fi

"$GITLEAKS_CMD" protect --staged --verbose --redact
