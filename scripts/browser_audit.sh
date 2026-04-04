#!/bin/bash
# scripts/browser_audit.sh: Orchestator for ChiefCritic browser-based user journey audits.

set -e

# Identify Project Root
PROJECT_ROOT=$(git rev-parse --show-toplevel)
cd "$PROJECT_ROOT"

# 1. Environment Preparation
export MOCK_AUTH=true
export ADMIN_CODE="agbalumo2024"
export PORT=8443
export AGBALUMO_ENV="development"
export BASE_URL="https://localhost:8443"

echo "🧹 Cleaning up old server instances..."
pkill server || true

# 2. Rebuild
echo "🔨 Rebuilding server and assets..."
SKIP_PRE_COMMIT=true ./scripts/verify_restart.sh

# 3. Launch
echo "🚀 Starting server at $BASE_URL with MOCK_AUTH=true..."
./bin/server serve > /tmp/agbalumo-audit.log 2>&1 &
SERVER_PID=$!

# 4. Wait for Readiness
echo "⏳ Waiting for healthz check..."
MAX_ATTEMPTS=15
ATTEMPT=0
while ! curl -k -s https://localhost:8443/healthz > /dev/null; do
    sleep 1
    ATTEMPT=$((ATTEMPT+1))
    if [ $ATTEMPT -ge $MAX_ATTEMPTS ]; then
        echo "❌ Server failed to start. Logs:"
        cat /tmp/agbalumo-audit.log
        kill $SERVER_PID || true
        exit 1
    fi
done

echo "✅ AUDIT SERVER READY (PID: $SERVER_PID)"
echo "--------------------------------------------------"
echo "AGENT INSTRUCTION:"
echo "1. Read journeys from .agents/user_journeys.yaml"
echo "2. Use browser_automation to play through journeys on $BASE_URL"
echo "3. Verify UX metrics (Layout Shift < 5px, Interaction < 100ms)"
echo "4. When done, run: kill $SERVER_PID"
echo "--------------------------------------------------"
