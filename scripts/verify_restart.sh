#!/bin/bash
set -e
export PATH=$PATH:/opt/homebrew/bin


# 1. Run Quality Checks (Tests + Coverage)
echo "🔍 Running Quality Checks..."
./scripts/pre-commit.sh

# 2. Build Assets & Server
echo "🎨 Building CSS..."
export PATH=$PATH:/opt/homebrew/bin
npm run build:css

echo "🔨 Building Server..."
mkdir -p bin
go build -o bin/agbalumo main.go

# 3. Restart the Server
echo "🔄 Restarting Server..."

# Find and kill existing process on port 8443
PID_HTTPS=$(lsof -ti:8443 || true)
if [ -n "$PID_HTTPS" ]; then
  echo "Killing existing HTTPS server (PID: $PID_HTTPS)..."
  kill -9 $PID_HTTPS
fi

# Find and kill existing process on port 8080
PID_HTTP=$(lsof -ti:8080 || true)
if [ -n "$PID_HTTP" ]; then
  echo "Killing existing HTTP server (PID: $PID_HTTP)..."
  kill -9 $PID_HTTP
fi

# Start the new server in the background
echo "🚀 Starting new server instance..."
nohup ./bin/agbalumo serve > server.log 2>&1 &
NEW_PID=$!
echo "Server started with PID: $NEW_PID"
echo "Logs are being written to server.log"

# Wait a moment to ensure it doesn't crash immediately
sleep 2
if ps -p $NEW_PID > /dev/null; then
   echo "✅ Server is running!"
else
   echo "❌ Server failed to start. Check server.log:"
   cat server.log
   exit 1
fi
