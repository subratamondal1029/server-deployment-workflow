#!/usr/bin/env bash
set -Eeuo pipefail

if [[ $# -ne 1 ]]; then
  echo "Usage: $0 <version>" >&2
  exit 1
fi

V="$1"

docker stop server || true
docker rm server || true

docker run -d --name server --restart unless-stopped -p 8000:8000 server:$V
