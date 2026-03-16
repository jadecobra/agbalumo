#!/bin/bash
# agent-gate.sh: Forwarded to harness verify command

set -e

# Robust PATH discovery
source "$(dirname "$0")/utils.sh"
setup_path

go run cmd/harness/main.go verify "$@"
