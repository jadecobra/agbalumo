#!/bin/bash
set -e

# Verify golangci-lint configuration is valid
# This catches config errors locally before CI

echo "Verifying golangci-lint configuration..."

if ! command -v golangci-lint >/dev/null 2>&1; then
    echo "Error: golangci-lint is not installed"
    exit 1
fi

GOLANGCI_VERSION=$(golangci-lint --version | head -1)
echo "Using: $GOLANGCI_VERSION"

CONFIG_FILE="scripts/.golangci.yml"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Config file $CONFIG_FILE not found"
    exit 1
fi

echo "Validating config file: $CONFIG_FILE"

if golangci-lint config verify --config="$CONFIG_FILE"; then
    echo "✅ Config is valid"
    exit 0
else
    echo "❌ Config validation failed"
    exit 1
fi
