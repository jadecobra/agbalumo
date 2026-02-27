#!/bin/sh
# scripts/setup-hooks.sh

# scripts/setup-hooks.sh

HOOK_DIR=".git/hooks"
PRE_COMMIT_HOOK="$HOOK_DIR/pre-commit"

echo "Setting up git hooks..."

if [ ! -d "$HOOK_DIR" ]; then
    echo "❌ Error: .git directory not found. Are you in the project root?"
    exit 1
fi

# Check for dependencies
check_dep() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "⚠️  Warning: $1 is not installed. Some hooks may fail."
  fi
}

check_dep "go"
check_dep "npm"
check_dep "git"
check_dep "bc"
check_dep "lsof"

# Backup existing hook if it's not ours
if [ -f "$PRE_COMMIT_HOOK" ]; then
    if ! grep -q "security-check.sh" "$PRE_COMMIT_HOOK"; then
        echo "📦 Backing up existing pre-commit hook to $PRE_COMMIT_HOOK.bak"
        cp "$PRE_COMMIT_HOOK" "$PRE_COMMIT_HOOK.bak"
    fi
fi

# Create pre-commit hook
cat > "$PRE_COMMIT_HOOK" <<EOF
#!/bin/sh
# agbalumo 10x Engineer Pre-commit Hook
# Runs security checks and quality checks before commit

./scripts/security-check.sh && ./scripts/pre-commit.sh
EOF

chmod +x "$PRE_COMMIT_HOOK"
echo "✅ Git pre-commit hook installed successfully!"

