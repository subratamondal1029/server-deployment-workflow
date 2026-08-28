#!/usr/bin/env bash
set -Eeuo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <version>" >&2
  exit 1
fi

V="$1"
PORT=8000
TEMP_PORT=8001
HEALTH_URL="http://localhost:${TEMP_PORT}/health"

docker build -t "server:$V" .

# start new container on a temp port, don't touch old one yet
docker rm -f server_new 2>/dev/null || true
docker run -d --name server_new -p "${TEMP_PORT}:8000" "server:$V"

# health-check with retries
for i in {1..10}; do
  if curl -fs "$HEALTH_URL" >/dev/null; then
    echo "New container healthy"
    break
  fi
  if [[ $i -eq 10 ]]; then
    echo "New container failed health check — rolling back" >&2
    docker logs server_new >&2 || true
    docker rm -f server_new
    exit 1
  fi
  sleep 1
done

# cutover: only now touch the live container
docker rm -f server_new
docker stop server || true
docker rm server || true
docker run -d --name server --restart unless-stopped -p "${PORT}:8000" "server:$V"

sleep 2
docker ps --filter "name=server" --filter "status=running" -q | grep -q . || {
  echo "Live container failed to start after cutover" >&2
  exit 1
}
