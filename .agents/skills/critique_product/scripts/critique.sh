#!/usr/bin/env bash
set -e

echo "[ChiefCritic] Executing critique_product skill..."
./scripts/verify_restart.sh

echo "Note: Before verifying the browser gate, ensure browser_subagent has tested the specific feature AND all existing user journeys in the UI."
./scripts/agent-exec.sh verify browser-verification

./scripts/agent-exec.sh workflow set-phase IDLE
./scripts/agent-exec.sh workflow init none
