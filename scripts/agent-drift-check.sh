#!/bin/bash
# agent-drift-check.sh: validates squad configuration and enforces double-commit rule.
set -e

# Robust PATH discovery
source "$(dirname "$0")/utils.sh"
setup_path

printf "${BLUE}Running Modular Agent Drift Check (10x Standard)...${NC}\n"

# Run the Go-based validation script
GOPATH="${PWD}/.tester/tmp/go" GOCACHE="${PWD}/.tester/tmp/gocache" go run scripts/verify-persona.go

printf "${GREEN}✅ Squad configuration is consistent and in sync with documentation.${NC}\n"
exit 0
