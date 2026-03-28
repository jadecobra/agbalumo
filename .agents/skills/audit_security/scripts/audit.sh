#!/usr/bin/env bash
set -e

echo "[SecurityEngineer] Executing audit_security skill..."
./scripts/agent-exec.sh verify security-static
