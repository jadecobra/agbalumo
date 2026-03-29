#!/bin/bash
# scripts/watch.sh
# 10x Engineer Feedback Loop (Watch Mode)

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting agbalumo 10x Feedback Loop...${NC}"
echo -e "${BLUE}👀 Watching for changes in .go and Taskfile.yml${NC}"

# Run task watch which uses go-task's native watch feature
# but wrap it with some helpful context
exec task watch
