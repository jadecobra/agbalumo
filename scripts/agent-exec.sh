#!/bin/bash
# agent-exec.sh: workflow manager for agbalumo.
# Ported to V2 Go Harness. This script is a compatibility wrapper.

set -e

# Backwards compatibility: trim "workflow" if it's the first argument
if [ "$1" == "workflow" ]; then
    shift
fi

go run cmd/harness/main.go "$@"
