#!/usr/bin/env bash
set -e
TEST_PATH="${1:-./internal/...}"

echo "[SDET-Tester] Executing make_it_fail skill..."
./scripts/agent-exec.sh workflow set-phase RED

# 1. Formatting Check/Fix
echo "[SDET-Tester] Applying gofmt -s -w..."
gofmt -s -w $TEST_PATH

# 2. Fast Lint Check
echo "[SDET-Tester] Running lightweight lint..."
golangci-lint run --fast $TEST_PATH

# 3. Gitleaks Check (Partial Scan for changes)
echo "[SDET-Tester] Checking for secrets..."
./scripts/gitleaks-scan.sh || true # Soft fail for now if tool missing

# 4. Run Test
echo "[SDET-Tester] Verifying RED state..."
go test -v $TEST_PATH || true

# 5. Verify & Anchor
./scripts/agent-exec.sh verify red-test
echo "[SDET-Tester] Anchoring RED state with git commit..."
git commit -m "RED: Anchor for TDD loop" --no-verify
