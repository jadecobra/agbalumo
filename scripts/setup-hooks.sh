#!/bin/sh
# scripts/setup-hooks.sh

HOOK_DIR=".git/hooks"
PRE_COMMIT_HOOK="$HOOK_DIR/pre-commit"

echo "Setting up git hooks..."

if [ ! -d "$HOOK_DIR" ]; then
    echo "Error: .git directory not found. Are you in the project root?"
    exit 1
fi

# Create pre-commit hook (symlink or copy)
# Using copy/script to ensure correct path execution
cat > "$PRE_COMMIT_HOOK" <<EOF
#!/bin/sh
./scripts/pre-commit.sh
EOF

chmod +x "$PRE_COMMIT_HOOK"
echo "✅ Git pre-commit hook installed!"
