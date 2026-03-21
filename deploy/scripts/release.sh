#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TAG="${1:-$(date +%Y%m%d%H%M%S)}"

echo "[release] root=${ROOT_DIR} tag=${TAG}"

cd "${ROOT_DIR}"
make lint
make test-backend
make test-frontend

echo "[release] build api image"
docker build -f deploy/docker/api.Dockerfile -t summon-king-api:"${TAG}" .

echo "[release] build worker image"
docker build -f deploy/docker/worker.Dockerfile -t summon-king-worker:"${TAG}" .

echo "[release] completed tag=${TAG}"
