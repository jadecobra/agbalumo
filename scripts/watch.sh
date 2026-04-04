#!/bin/bash
# scripts/watch.sh
# 10x Engineer Feedback Loop (Watch Mode)

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting agbalumo 10x Feedback Loop...${NC}"
echo -e "${BLUE}👀 Watching for changes (Pure Go engine)...${NC}"

# Run the native Go watcher
exec go run cmd/verify/main.go watch go run main.go serve
