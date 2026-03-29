#!/usr/bin/env bash
set -e
FEATURE_NAME="${1:-new-feature}"

echo "[ProductOwner/SystemsArchitect] Executing design_architecture skill..."
./scripts/agent-exec.sh workflow init "$FEATURE_NAME"
swagger-cli validate docs/openapi.yaml || echo "Warning: Spec validation failed"
./scripts/agent-exec.sh verify api-spec
