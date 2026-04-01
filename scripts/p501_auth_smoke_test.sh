#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18089}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18089}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d)"
BACKEND_PID=""

cleanup() {
  if [[ -n "${BACKEND_PID}" ]]; then
    kill "${BACKEND_PID}" >/dev/null 2>&1 || true
    wait "${BACKEND_PID}" 2>/dev/null || true
  fi
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

log() {
  echo "[p501-auth-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p501-auth-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
    exit 1
  fi
}

request_json() {
  local method="$1"
  local path="$2"
  local body="${3:-}"
  local auth="${4:-}"
  local confirm="${5:-no}"
  local outfile="$6"

  local auth_header=()
  if [[ -n "${auth}" ]]; then
    auth_header=(-H "Authorization: Bearer ${auth}")
  fi
  local confirm_header=()
  if [[ "${confirm}" == "yes" ]]; then
    confirm_header=(-H "X-Action-Confirm: CONFIRM")
  fi

  if [[ -n "${body}" ]]; then
    curl -sS -X "${method}" \
      -H "Content-Type: application/json" \
      "${auth_header[@]}" \
      "${confirm_header[@]}" \
      -d "${body}" \
      -o "${outfile}" \
      -w "%{http_code}" \
      "${BASE_URL}${path}"
  else
    curl -sS -X "${method}" \
      -H "Content-Type: application/json" \
      "${auth_header[@]}" \
      "${confirm_header[@]}" \
      -o "${outfile}" \
      -w "%{http_code}" \
      "${BASE_URL}${path}"
  fi
}

json_field() {
  local file="$1"
  local field="$2"
  python3 - "$file" "$field" <<'PY'
import json,sys
with open(sys.argv[1],'r',encoding='utf-8') as f:
    data=json.load(f)
v=data.get(sys.argv[2],"")
if isinstance(v,list):
    print(",".join(str(x) for x in v))
else:
    print(v)
PY
}

ensure_backend() {
  log "starting backend on ${KM_LISTEN_ADDR}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p501-auth-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..40}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done
  echo "[p501-auth-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p501-auth-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local code body access refresh readonly_access readonly_refresh

  body="${TMP_DIR}/login-admin.json"
  code="$(request_json POST /api/v1/auth/login '{"username":"admin","password":"123456"}' "" no "${body}")"
  expect_status "${code}" "200" "admin login"
  access="$(json_field "${body}" "accessToken")"
  refresh="$(json_field "${body}" "refreshToken")"
  if [[ -z "${access}" || -z "${refresh}" ]]; then
    echo "[p501-auth-smoke][FAIL] admin login missing tokens"
    exit 1
  fi

  body="${TMP_DIR}/create-readonly.json"
  code="$(request_json POST /api/v1/auth/users '{"username":"readonly1","password":"123456","role":"readonly","allowedNamespaces":["dev"]}' "${access}" yes "${body}")"
  if [[ "${code}" != "201" && "${code}" != "409" ]]; then
    echo "[p501-auth-smoke][FAIL] create readonly user expected 201 or 409 actual=${code}"
    exit 1
  fi

  body="${TMP_DIR}/login-readonly.json"
  code="$(request_json POST /api/v1/auth/login '{"username":"readonly1","password":"123456"}' "" no "${body}")"
  expect_status "${code}" "200" "readonly login"
  readonly_access="$(json_field "${body}" "accessToken")"
  readonly_refresh="$(json_field "${body}" "refreshToken")"

  body="${TMP_DIR}/readonly-create-user-deny.json"
  code="$(request_json POST /api/v1/auth/users '{"username":"tmpuser","password":"123456","role":"readonly","allowedNamespaces":["dev"]}' "${readonly_access}" yes "${body}")"
  expect_status "${code}" "403" "readonly create user should be forbidden"

  body="${TMP_DIR}/readonly-me.json"
  code="$(request_json GET /api/v1/auth/me "" "${readonly_access}" no "${body}")"
  expect_status "${code}" "200" "readonly get me"
  role="$(json_field "${body}" "role")"
  if [[ "${role}" != "readonly" ]]; then
    echo "[p501-auth-smoke][FAIL] readonly role mismatch: ${role}"
    exit 1
  fi

  body="${TMP_DIR}/refresh-readonly.json"
  code="$(request_json POST /api/v1/auth/refresh "{\"refreshToken\":\"${readonly_refresh}\"}" "" no "${body}")"
  expect_status "${code}" "200" "readonly refresh"

  body="${TMP_DIR}/logout-readonly.json"
  code="$(request_json POST /api/v1/auth/logout "{\"refreshToken\":\"${readonly_refresh}\"}" "" no "${body}")"
  expect_status "${code}" "204" "readonly logout"

  body="${TMP_DIR}/refresh-revoked.json"
  code="$(request_json POST /api/v1/auth/refresh "{\"refreshToken\":\"${readonly_refresh}\"}" "" no "${body}")"
  expect_status "${code}" "401" "refresh revoked token"

  log "all auth smoke checks passed"
}

main "$@"

