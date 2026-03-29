#!/bin/bash
# agent-exec.sh: workflow manager for agbalumo.
# Ported to V2 Go Harness. This script is a compatibility wrapper.

set -e

# Backwards compatibility: trim "workflow" if it's the first argument
if [ "$1" == "workflow" ]; then
    shift
fi

ISOLATE_GO="${ISOLATE_GO:-true}"
if [ "$ISOLATE_GO" == "true" ]; then
    export GOPATH="${PWD}/.tester/tmp/go"
    export GOCACHE="${PWD}/.tester/tmp/gocache"
fi

go run cmd/harness/main.go "$@"
