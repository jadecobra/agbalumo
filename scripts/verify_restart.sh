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
if [ "$SKIP_PRE_COMMIT" != "true" ]; then
    echo "🔍 Running Quality Checks..."
    go run cmd/verify/main.go precommit
else
    echo "🔍 Skipping pre-commit checks..."
fi

# 2. Build Assets & Server
echo "🎨 Building CSS..."
npm run build:css

echo "🔨 Building Server..."
mkdir -p /tmp/.tester/servers
mkdir -p /tmp/.tester/data

# Use workspace-local Go paths to avoid permission issues
ISOLATE_GO="${ISOLATE_GO:-true}"
if [ "$ISOLATE_GO" == "true" ]; then
    export GOPATH="$(pwd)/.tester/tmp/go"
    export GOCACHE="$(pwd)/.tester/tmp/gocache"
fi
go build -o /tmp/.tester/servers/agbalumo main.go

# 3. Restart the Server
echo "🔄 Restarting Server..."

# Function for graceful shutdown (portable: works on macOS bash 3.2+)
shutdown_port() {
  local port=$1
  local pid

  # lsof -t returns one PID per line; iterate each one individually
  while IFS= read -r pid; do
    [ -z "$pid" ] && continue
    echo "Attempting graceful shutdown of port $port (PID: $pid)..."
    kill "$pid" 2>/dev/null || true  # SIGTERM

    # Wait up to 5 seconds for graceful exit
    local i
    for i in 1 2 3 4 5; do
      if ! ps -p "$pid" > /dev/null 2>&1; then
        echo "✅ Process $pid on port $port exited gracefully."
        break
      fi
      sleep 1
    done

    # Force-kill if still alive
    if ps -p "$pid" > /dev/null 2>&1; then
      echo "⚠️  Process $pid on port $port did not exit, forcing (SIGKILL)..."
      kill -9 "$pid" 2>/dev/null || true
    fi
  done < <(lsof -ti:"$port" 2>/dev/null || true)
}

shutdown_port 8443
shutdown_port 8080

# Start the new server in the background
echo "🚀 Starting new server instance..."
export AGBALUMO_ENV=development
nohup /tmp/.tester/servers/agbalumo serve > /tmp/.tester/servers/server.log 2>&1 &
NEW_PID=$!
echo "Server started with PID: $NEW_PID"
echo "Logs are being written to /tmp/.tester/servers/server.log"

# Wait a moment to ensure it doesn't crash immediately
sleep 2
if ps -p $NEW_PID > /dev/null; then
   echo "✅ Server is running!"
else
   echo "❌ Server failed to start. Check /tmp/.tester/servers/server.log:"
   cat /tmp/.tester/servers/server.log
   exit 1
fi

