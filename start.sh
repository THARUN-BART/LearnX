#!/bin/sh
set -e

echo "=== Starting LearnX Single-Server Deployment on Render ==="

# 1. Start Pedagogy Python Microservice in background on 127.0.0.1:8001
echo "-> Starting Pedagogy Microservice on 127.0.0.1:8001..."
cd /app/pedagogy-service
python3 -m uvicorn person_c.api:app --host 127.0.0.1 --port 8001 &
PEDAGOGY_PID=$!

# Wait for Pedagogy service to report healthy
echo "-> Awaiting Pedagogy Service health..."
for i in $(seq 1 30); do
  if curl -s http://127.0.0.1:8001/health > /dev/null 2>&1; then
    echo "-> Pedagogy Microservice is healthy and ready!"
    break
  fi
  sleep 0.5
done

# 2. Configure Go Backend to bind to 0.0.0.0:${PORT:-8080}
cd /app
export SERVER_HOST="0.0.0.0"
export SERVER_PORT="${PORT:-8080}"
export PEDAGOGY_SERVICE_URL="http://127.0.0.1:8001"

echo "-> Starting Go Backend on ${SERVER_HOST}:${SERVER_PORT}..."

# Handle termination gracefully
trap "kill -TERM $PEDAGOGY_PID 2>/dev/null; exit 0" TERM INT

# Exec Go backend as foreground process
exec /app/learnx-server
