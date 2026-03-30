#!/usr/bin/env bash

set -e

# Initialize Sandboxed Workspace Directories
PROJECT_ROOT=$(pwd)
SANDBOX_DIR="$PROJECT_ROOT/.tester/tmp"

echo "🛠️  Initializing Sandboxed Workspace in $SANDBOX_DIR..."

# Create necessary subdirectories
mkdir -p "$SANDBOX_DIR/go/pkg/mod"
mkdir -p "$SANDBOX_DIR/go/cache"
mkdir -p "$SANDBOX_DIR/gh"

# Verify .gitignore (already checked, but safety first)
if ! grep -q ".tester/tmp/" "$PROJECT_ROOT/.gitignore"; then
    echo "Adding .tester/tmp/ to .gitignore..."
    echo ".tester/tmp/" >> "$PROJECT_ROOT/.gitignore"
fi

echo "✅ Sandboxed directories created successfully."
echo "👉 To activate the sandbox in your terminal, run:"
echo "   source scripts/sandbox.env"
