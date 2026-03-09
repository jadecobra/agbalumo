#!/bin/bash
# scripts/ci-local.sh
# Helper script to run GitHub Actions locally using act

set -e

# Ensure we are in the root directory
cd "$(dirname "$0")/.."

# Check if act is installed
if ! command -v act &> /dev/null; then
    echo "❌ act is not installed. Please install it with 'brew install act'."
    exit 1
fi

# Detect Apple M-series (arm64) and apply architecture flag
ARCH_FLAG=""
if [[ $(uname -m) == "arm64" ]]; then
    echo "Detected Apple M-series chip. Using --container-architecture linux/amd64"
    ARCH_FLAG="--container-architecture linux/amd64"
fi

# Run act with provided arguments
echo "🚀 Running local CI with act..."
act $ARCH_FLAG "$@"
