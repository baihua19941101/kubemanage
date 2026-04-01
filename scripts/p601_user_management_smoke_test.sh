#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18090}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18090}"
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
  echo "[p601-user-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p601-user-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
with open(sys.argv[1], 'r', encoding='utf-8') as f:
    data = json.load(f)
v = data.get(sys.argv[2], "")
if isinstance(v, list):
    print(",".join(str(x) for x in v))
else:
    print(v)
PY
}

json_has_username() {
  local file="$1"
  local username="$2"
  python3 - "$file" "$username" <<'PY'
import json,sys
with open(sys.argv[1], 'r', encoding='utf-8') as f:
    data = json.load(f)
items = data.get("items", [])
print("yes" if any(str(x.get("username", "")) == sys.argv[2] for x in items) else "no")
PY
}

ensure_backend() {
  log "starting backend on ${KM_LISTEN_ADDR}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p601-user-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..40}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p601-user-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p601-user-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local code body access username new_password readonly_access exists
  username="p601_readonly_$(date +%s)"
  new_password="654321"

  body="${TMP_DIR}/login-admin.json"
  code="$(request_json POST /api/v1/auth/login '{"username":"admin","password":"123456"}' "" no "${body}")"
  expect_status "${code}" "200" "admin login"
  access="$(json_field "${body}" "accessToken")"
  if [[ -z "${access}" ]]; then
    echo "[p601-user-smoke][FAIL] missing admin access token"
    exit 1
  fi

  body="${TMP_DIR}/create-user.json"
  code="$(request_json POST /api/v1/auth/users "{\"username\":\"${username}\",\"password\":\"123456\",\"role\":\"readonly\",\"allowedNamespaces\":[\"dev\"]}" "${access}" yes "${body}")"
  expect_status "${code}" "201" "create readonly user"

  body="${TMP_DIR}/list-users-admin.json"
  code="$(request_json GET /api/v1/auth/users "" "${access}" no "${body}")"
  expect_status "${code}" "200" "admin list users"
  exists="$(json_has_username "${body}" "${username}")"
  if [[ "${exists}" != "yes" ]]; then
    echo "[p601-user-smoke][FAIL] created user not found in list"
    exit 1
  fi

  body="${TMP_DIR}/disable-user.json"
  code="$(request_json PATCH "/api/v1/auth/users/${username}/status" '{"isActive":false}' "${access}" yes "${body}")"
  expect_status "${code}" "204" "disable user"

  body="${TMP_DIR}/login-disabled.json"
  code="$(request_json POST /api/v1/auth/login "{\"username\":\"${username}\",\"password\":\"123456\"}" "" no "${body}")"
  expect_status "${code}" "403" "disabled user login"

  body="${TMP_DIR}/reset-password.json"
  code="$(request_json POST "/api/v1/auth/users/${username}/reset-password" "{\"password\":\"${new_password}\"}" "${access}" yes "${body}")"
  expect_status "${code}" "204" "reset password"

  body="${TMP_DIR}/enable-user.json"
  code="$(request_json PATCH "/api/v1/auth/users/${username}/status" '{"isActive":true}' "${access}" yes "${body}")"
  expect_status "${code}" "204" "enable user"

  body="${TMP_DIR}/login-reset-password.json"
  code="$(request_json POST /api/v1/auth/login "{\"username\":\"${username}\",\"password\":\"${new_password}\"}" "" no "${body}")"
  expect_status "${code}" "200" "login with reset password"
  readonly_access="$(json_field "${body}" "accessToken")"

  body="${TMP_DIR}/list-users-readonly.json"
  code="$(request_json GET /api/v1/auth/users "" "${readonly_access}" no "${body}")"
  expect_status "${code}" "403" "readonly list users should be forbidden"

  log "all p601 user management smoke checks passed"
}

main "$@"
