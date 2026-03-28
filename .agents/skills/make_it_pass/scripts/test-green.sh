#!/usr/bin/env bash
set -e
TEST_PATH="${1:-./internal/...}"

echo "[BackendEngineer] Executing make_it_pass skill..."
go test -v $TEST_PATH
go test -v ./cmd/...
./scripts/agent-exec.sh verify implementation
