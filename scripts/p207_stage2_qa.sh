#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BASE_PORT="${BASE_PORT:-18090}"

log() {
  echo "[p207-qa] $*"
}

run_step() {
  local name="$1"
  shift
  log "run ${name}"
  "$@"
}

request() {
  local base_url="$1"
  local method="$2"
  local path="$3"
  local role="$4"
  local outfile="$5"

  curl -sS -X "${method}" \
    -H "X-User: stage2-qa" \
    -H "X-User-Role: ${role}" \
    -o "${outfile}" \
    -w "%{http_code}" \
    "${base_url}${path}"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p207-qa][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

check_p204() {
  local port="$1"
  local base_url="http://127.0.0.1:${port}"
  local root_dir="${ROOT_DIR}"
  local tmp_dir
  local backend_pid=""
  tmp_dir="$(mktemp -d)"

  cleanup() {
    if [[ -n "${backend_pid:-}" ]]; then
      kill "${backend_pid}" >/dev/null 2>&1 || true
    fi
    if [[ -n "${tmp_dir:-}" ]]; then
      rm -rf "${tmp_dir}"
    fi
  }
  trap cleanup RETURN

  (
    cd "${root_dir}/backend"
    KM_LISTEN_ADDR=":${port}" GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p207-p204.log 2>&1
  ) &
  backend_pid="$!"

  for _ in {1..30}; do
    if curl -sS "${base_url}/api/v1/healthz" >/dev/null 2>&1; then
      break
    fi
    sleep 1
  done

  local body code
  body="${tmp_dir}/pvs.json"
  code="$(request "${base_url}" GET /api/v1/pvs viewer "${body}")"
  expect_status "${code}" "200" "list pvs"

  body="${tmp_dir}/pvcs.json"
  code="$(request "${base_url}" GET /api/v1/pvcs viewer "${body}")"
  expect_status "${code}" "200" "list pvcs"

  body="${tmp_dir}/storageclasses.json"
  code="$(request "${base_url}" GET /api/v1/storageclasses viewer "${body}")"
  expect_status "${code}" "200" "list storageclasses"

  body="${tmp_dir}/configmaps.json"
  code="$(request "${base_url}" GET /api/v1/configmaps viewer "${body}")"
  expect_status "${code}" "200" "list configmaps"

  body="${tmp_dir}/secrets.json"
  code="$(request "${base_url}" GET /api/v1/secrets viewer "${body}")"
  expect_status "${code}" "200" "list secrets"

  if ! grep -Eq '\*{4,}' "${body}"; then
    echo "[p207-qa][FAIL] secret list should contain masked secret data"
    exit 1
  fi

  log "p204 checks passed"
}

main() {
  run_step "backend tests" bash -lc "cd '${ROOT_DIR}/backend' && GOCACHE=/tmp/go-build-cache go test ./..."
  run_step "frontend build" bash -lc "cd '${ROOT_DIR}/frontend' && npm run build"

  run_step "p202 smoke" bash -lc "cd '${ROOT_DIR}' && BASE_URL=http://127.0.0.1:$((BASE_PORT + 2)) KM_LISTEN_ADDR=:$((BASE_PORT + 2)) ./scripts/p202_smoke_test.sh"
  run_step "p203 smoke" bash -lc "cd '${ROOT_DIR}' && BASE_URL=http://127.0.0.1:$((BASE_PORT + 3)) KM_LISTEN_ADDR=:$((BASE_PORT + 3)) ./scripts/p203_smoke_test.sh"
  run_step "p204 checks" check_p204 "$((BASE_PORT + 4))"
  run_step "p205 smoke" bash -lc "cd '${ROOT_DIR}' && BASE_URL=http://127.0.0.1:$((BASE_PORT + 5)) KM_LISTEN_ADDR=:$((BASE_PORT + 5)) ./scripts/p205_smoke_test.sh"
  run_step "p206 smoke" bash -lc "cd '${ROOT_DIR}' && BASE_URL=http://127.0.0.1:$((BASE_PORT + 6)) KM_LISTEN_ADDR=:$((BASE_PORT + 6)) ./scripts/p206_smoke_test.sh"

  log "stage2 qa passed"
}

main "$@"
