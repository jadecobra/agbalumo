#!/bin/bash
# scripts/gitleaks-scan.sh
set -e
source "$(dirname "$0")/utils.sh"
setup_path

if ! command -v gitleaks >/dev/null 2>&1; then
    if [ "$FMT" != "json" ]; then echo "  ${YELLOW}Installing gitleaks...${NC}"; fi
    if [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew >/dev/null 2>&1; then
            brew install gitleaks > /dev/null 2>&1
        else
            VERSION="8.18.2"
            ARCH="x64"
            if [[ "$(uname -m)" == "arm64" ]]; then ARCH="arm64"; fi
            curl -sSL "https://github.com/gitleaks/gitleaks/releases/download/v${VERSION}/gitleaks_${VERSION}_darwin_${ARCH}.tar.gz" | tar -xz -C /tmp gitleaks
            export PATH="$PATH:/tmp"
        fi
    else
        VERSION="8.18.2"
        ARCH="x64"
        if [[ "$(uname -m)" == "aarch64" || "$(uname -m)" == "arm64" ]]; then ARCH="arm64"; fi
        curl -sSL "https://github.com/gitleaks/gitleaks/releases/download/v${VERSION}/gitleaks_${VERSION}_linux_${ARCH}.tar.gz" | tar -xz -C /tmp gitleaks
        export PATH="$PATH:/tmp"
    fi
fi

if ! command -v gitleaks >/dev/null 2>&1; then
     echo "  ${RED}❌ Error: gitleaks is required but could not be installed.${NC}"
     exit 1
fi

gitleaks protect --staged --verbose --redact
