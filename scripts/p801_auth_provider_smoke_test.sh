#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18092}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18092}"
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
  echo "[p801-provider-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p801-provider-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
if isinstance(v, list):
    print(",".join(str(x) for x in v))
else:
    print(v)
PY
}

json_provider_id() {
  local file="$1"
  local name="$2"
  python3 - "$file" "$name" <<'PY'
import json,sys
with open(sys.argv[1], "r", encoding="utf-8") as f:
    data = json.load(f)
for item in data.get("items", []):
    if str(item.get("name","")) == sys.argv[2]:
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
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p801-provider-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..40}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p801-provider-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p801-provider-smoke.log"
  exit 1
}

main() {
  ensure_backend

  local code body admin_access provider_name provider_id
  provider_name="ldap-$(date +%s)"

  body="${TMP_DIR}/public-providers.json"
  code="$(request_json GET /api/v1/auth/providers/public "" "" no "${body}")"
  expect_status "${code}" "200" "list public providers"

  body="${TMP_DIR}/admin-login.json"
  code="$(request_json POST /api/v1/auth/login '{"username":"admin","password":"123456","provider":"local"}' "" no "${body}")"
  expect_status "${code}" "200" "admin login local"
  admin_access="$(json_field "${body}" "accessToken")"

  body="${TMP_DIR}/create-provider.json"
  code="$(request_json POST /api/v1/auth/providers "{\"name\":\"${provider_name}\",\"type\":\"ldap\",\"config\":{\"url\":\"ldap://127.0.0.1:389\",\"baseDN\":\"dc=example,dc=com\"}}" "${admin_access}" yes "${body}")"
  expect_status "${code}" "201" "create ldap provider"

  body="${TMP_DIR}/list-providers.json"
  code="$(request_json GET /api/v1/auth/providers "" "${admin_access}" no "${body}")"
  expect_status "${code}" "200" "list auth providers"
  provider_id="$(json_provider_id "${body}" "${provider_name}")"
  if [[ -z "${provider_id}" ]]; then
    echo "[p801-provider-smoke][FAIL] created provider not found"
    exit 1
  fi

  body="${TMP_DIR}/set-default-ldap.json"
  code="$(request_json POST "/api/v1/auth/providers/${provider_id}/default" "" "${admin_access}" yes "${body}")"
  expect_status "${code}" "204" "set ldap provider as default"

  body="${TMP_DIR}/login-default-ldap.json"
  code="$(request_json POST /api/v1/auth/login '{"username":"admin","password":"123456"}' "" no "${body}")"
  expect_status "${code}" "502" "default ldap login should return unavailable"

  body="${TMP_DIR}/login-local-again.json"
  code="$(request_json POST /api/v1/auth/login '{"username":"admin","password":"123456","provider":"local"}' "" no "${body}")"
  expect_status "${code}" "200" "login with local provider should still work"

  log "all p801 auth provider smoke checks passed"
}

main "$@"
