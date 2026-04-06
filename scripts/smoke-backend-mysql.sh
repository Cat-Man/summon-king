#!/usr/bin/env bash
set -eEuo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
COMPOSE_FILE="${ROOT_DIR}/docker-compose.mysql.yml"
PROJECT_NAME="${COMPOSE_PROJECT_NAME:-zhzw-mysql-smoke}"
MYSQL_PORT="${MYSQL_PORT:-33306}"
MYSQL_DATABASE="${MYSQL_DATABASE:-zhzw}"
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-secret}"
MYSQL_SMOKE_DSN="${MYSQL_SMOKE_DSN:-root:${MYSQL_ROOT_PASSWORD}@tcp(127.0.0.1:${MYSQL_PORT})/${MYSQL_DATABASE}?parseTime=true&multiStatements=true}"

cleanup() {
  docker compose -f "${COMPOSE_FILE}" -p "${PROJECT_NAME}" down -v --remove-orphans >/dev/null 2>&1 || true
}

dump_logs() {
  echo "mysql smoke failed, dumping container logs..." >&2
  docker compose -f "${COMPOSE_FILE}" -p "${PROJECT_NAME}" logs --no-color mysql >&2 || true
}

trap cleanup EXIT
trap dump_logs ERR

docker compose -f "${COMPOSE_FILE}" -p "${PROJECT_NAME}" up -d --wait

(
  cd "${ROOT_DIR}/apps/backend"
  MYSQL_SMOKE_DSN="${MYSQL_SMOKE_DSN}" \
  go test ./internal/bootstrap -run '^TestNewRouterWithConfig_UsesRealMySQLStorageAndPersistsAcrossRebuilds$$' -count=1 -v
)
