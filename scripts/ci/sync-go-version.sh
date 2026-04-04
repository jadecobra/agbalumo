#!/bin/bash
# scripts/ci/sync-go-version.sh
# Extracts the Go version from go.mod for use in CI workflows.

set -e

GO_MOD_PATH="${1:-go.mod}"

if [ ! -f "$GO_MOD_PATH" ]; then
    echo "Error: $GO_MOD_PATH not found" >&2
    exit 1
fi

# Extract version (e.g., "go 1.25.0" -> "1.25.0")
VERSION=$(grep "^go " "$GO_MOD_PATH" | awk '{print $2}')

if [ -z "$VERSION" ]; then
    echo "Error: Could not find go version in $GO_MOD_PATH" >&2
    exit 1
fi

echo "$VERSION"
