#!/bin/sh
# scripts/setup-hooks.sh

HOOK_DIR=".git/hooks"
PRE_COMMIT_HOOK="$HOOK_DIR/pre-commit"
PRE_PUSH_HOOK="$HOOK_DIR/pre-push"

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
check_dep "task"

# Backup existing hook if it's not ours
if [ -f "$PRE_COMMIT_HOOK" ]; then
    if ! grep -q "task pre-commit" "$PRE_COMMIT_HOOK"; then
        echo "📦 Backing up existing pre-commit hook to $PRE_COMMIT_HOOK.bak"
        cp "$PRE_COMMIT_HOOK" "$PRE_COMMIT_HOOK.bak"
    fi
fi

# Creates pre-commit hook with sequential fast-checks
cat > "$PRE_COMMIT_HOOK" <<EOF
#!/bin/sh
# agbalumo 10x Engineer Pre-commit Hook
# 1. Fast Checks (Fmt/Lint/Build) -> 2. Heavy Checks (Tests/Audit)

task pre-commit
EOF

chmod +x "$PRE_COMMIT_HOOK"
echo "✅ Git pre-commit hook installed successfully!"

# Creates pre-push hook with full CI validation
cat > "$PRE_PUSH_HOOK" <<EOF
#!/bin/sh
# agbalumo 10x Engineer Pre-push Hook
# Runs the full production CI suite before pushing to remote.
# Bypass using: git push --no-verify
echo "🚀 Running comprehensive CI verification before push..."
task ci
if [ \$? -ne 0 ]; then
  echo "❌ CI validation failed! Fix errors or use 'git push --no-verify' to ignore."
  exit 1
fi
EOF
chmod +x "$PRE_PUSH_HOOK"
echo "✅ Git pre-push hook installed successfully!"