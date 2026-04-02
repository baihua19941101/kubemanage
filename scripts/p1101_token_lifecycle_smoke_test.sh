#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18095}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18095}"
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
  echo "[p1101-token-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p1101-token-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
v = data.get(sys.argv[2], "")
print(v if not isinstance(v, list) else ",".join(str(x) for x in v))
PY
}

json_latest_token_id_for_user() {
  local file="$1"
  local username="$2"
  python3 - "$file" "$username" <<'PY'
import json,sys
items = json.load(open(sys.argv[1], "r", encoding="utf-8")).get("items", [])
for item in items:
    if str(item.get("username","")) == sys.argv[2]:
        print(item.get("id",""))
        break
else:
    print("")
PY
}

ensure_backend() {
  log "starting backend on ${KM_LISTEN_ADDR}"
  (
    cd "${ROOT_DIR}/backend"
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p1101-token-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..45}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p1101-token-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p1101-token-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local code body admin_access username token_id
  username="p1101_user_$(date +%s)"

  body="${TMP_DIR}/admin-login.json"
  code="$(request_json POST /api/v1/auth/login '{"username":"admin","password":"123456","provider":"local"}' "" no "${body}")"
  expect_status "${code}" "200" "admin login"
  admin_access="$(json_field "${body}" "accessToken")"
  if [[ -z "${admin_access}" ]]; then
    echo "[p1101-token-smoke][FAIL] missing admin access token"
    exit 1
  fi

  body="${TMP_DIR}/create-user.json"
  code="$(request_json POST /api/v1/auth/users "{\"username\":\"${username}\",\"password\":\"123456\",\"role\":\"readonly\",\"allowedNamespaces\":[\"dev\"]}" "${admin_access}" yes "${body}")"
  expect_status "${code}" "201" "create user"

  body="${TMP_DIR}/user-login.json"
  code="$(request_json POST /api/v1/auth/login "{\"username\":\"${username}\",\"password\":\"123456\",\"provider\":\"local\"}" "" no "${body}")"
  expect_status "${code}" "200" "user login"

  body="${TMP_DIR}/list-tokens.json"
  code="$(request_json GET "/api/v1/auth/tokens?username=${username}&activeOnly=true&limit=20" "" "${admin_access}" no "${body}")"
  expect_status "${code}" "200" "list tokens"
  token_id="$(json_latest_token_id_for_user "${body}" "${username}")"
  if [[ -z "${token_id}" ]]; then
    echo "[p1101-token-smoke][FAIL] token session id not found for user ${username}"
    exit 1
  fi

  body="${TMP_DIR}/revoke-token.json"
  code="$(request_json POST "/api/v1/auth/tokens/${token_id}/revoke" "" "${admin_access}" yes "${body}")"
  expect_status "${code}" "204" "revoke token by id"

  body="${TMP_DIR}/revoke-all.json"
  code="$(request_json POST /api/v1/auth/tokens/revoke-all "{\"username\":\"${username}\"}" "${admin_access}" yes "${body}")"
  expect_status "${code}" "200" "revoke all by username"

  log "all p1101 token lifecycle smoke checks passed"
}

main "$@"

