#!/bin/bash
set -e
# Robust PATH discovery for macOS and Linux
for dir in /usr/local/bin /opt/homebrew/bin /usr/bin /bin; do
    case ":$PATH:" in
        *":$dir:"*) ;;
        *) export PATH="$PATH:$dir" ;;
    esac
done


# 1. Run Quality Checks (Tests + Coverage)
echo "🔍 Running Quality Checks..."
./scripts/pre-commit.sh

# 2. Build Assets & Server
echo "🎨 Building CSS..."
npm run build:css

echo "🔨 Building Server..."
mkdir -p bin
go build -o bin/agbalumo main.go

# 3. Restart the Server
echo "🔄 Restarting Server..."

# Function for graceful shutdown
shutdown_port() {
  local port=$1
  local pid=$(lsof -ti:$port || true)
  if [ -n "$pid" ]; then
    echo "Attempting graceful shutdown of port $port (PID: $pid)..."
    kill "$pid" # SIGTERM
    
    # Wait up to 5 seconds for it to exit
    for i in {1..5}; do
      if ! ps -p "$pid" > /dev/null; then
        echo "✅ Process on port $port exited gracefully."
        return 0
      fi
      sleep 1
    done
    
    echo "⚠️  Process on port $port did not exit, forcing (SIGKILL)..."
    kill -9 "$pid"
  fi
}

shutdown_port 8443
shutdown_port 8080

# Start the new server in the background
echo "🚀 Starting new server instance..."
export AGBALUMO_ENV=development
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

