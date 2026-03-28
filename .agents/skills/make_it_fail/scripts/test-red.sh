#!/usr/bin/env bash
set -e
TEST_PATH="${1:-./internal/...}"

echo "[SDET-Tester] Executing make_it_fail skill..."
./scripts/agent-exec.sh workflow set-phase RED
go test -v $TEST_PATH || true
./scripts/agent-exec.sh verify red-test
