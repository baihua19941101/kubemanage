#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:18093}"
KM_LISTEN_ADDR="${KM_LISTEN_ADDR:-:18093}"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d)"
BACKEND_PID=""
LDAP_CONTAINER=""
LDAP_PORT="${LDAP_PORT:-31389}"

# 默认启动本地测试 LDAP（rroemhild/test-openldap，监听 10389）；如需外部 LDAP，可显式传入 LDAP_URL。
LDAP_URL="${LDAP_URL:-}"
LDAP_BASE_DN="${LDAP_BASE_DN:-dc=planetexpress,dc=com}"
LDAP_LOGIN_ATTR="${LDAP_LOGIN_ATTR:-uid}"
LDAP_USER_FILTER="${LDAP_USER_FILTER:-}"
if [[ -z "${LDAP_USER_FILTER}" ]]; then
  LDAP_USER_FILTER='({{loginAttr}}={{username}})'
fi
LDAP_BIND_DN="${LDAP_BIND_DN:-cn=admin,dc=planetexpress,dc=com}"
LDAP_BIND_PASSWORD="${LDAP_BIND_PASSWORD:-GoodNewsEveryone}"
LDAP_TEST_USER="${LDAP_TEST_USER:-fry}"
LDAP_TEST_PASSWORD="${LDAP_TEST_PASSWORD:-fry}"

cleanup() {
  if [[ -n "${BACKEND_PID}" ]]; then
    kill "${BACKEND_PID}" >/dev/null 2>&1 || true
    wait "${BACKEND_PID}" 2>/dev/null || true
  fi
  if [[ -n "${LDAP_CONTAINER}" ]]; then
    docker rm -f "${LDAP_CONTAINER}" >/dev/null 2>&1 || true
  fi
  rm -rf "${TMP_DIR}"
}
trap cleanup EXIT

log() {
  echo "[p901-ldap-smoke] $*"
}

expect_status() {
  local actual="$1"
  local expected="$2"
  local msg="$3"
  if [[ "${actual}" != "${expected}" ]]; then
    echo "[p901-ldap-smoke][FAIL] ${msg} expected=${expected} actual=${actual}"
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
    KM_LISTEN_ADDR="${KM_LISTEN_ADDR}" KM_K8S_ADAPTER_MODE=mock GOCACHE=/tmp/go-build-cache go run ./cmd/server >/tmp/kubemanage-backend-p901-ldap-smoke.log 2>&1
  ) &
  BACKEND_PID="$!"

  for _ in {1..45}; do
    if curl -sS "${BASE_URL}/api/v1/healthz" >/dev/null 2>&1; then
      log "backend startup ready"
      return
    fi
    sleep 1
  done

  echo "[p901-ldap-smoke][FAIL] backend failed to start, see /tmp/kubemanage-backend-p901-ldap-smoke.log"
  exit 1
}

ensure_ldap() {
  if [[ -n "${LDAP_URL}" ]]; then
    log "using external LDAP URL: ${LDAP_URL}"
    return
  fi

  LDAP_CONTAINER="p901-openldap-$(date +%s)"
  log "starting local ldap container ${LDAP_CONTAINER} on port ${LDAP_PORT}"
  docker run -d --rm --name "${LDAP_CONTAINER}" -p "${LDAP_PORT}:10389" rroemhild/test-openldap:latest >/dev/null
  LDAP_URL="ldap://127.0.0.1:${LDAP_PORT}"

  for _ in {1..40}; do
    if docker exec "${LDAP_CONTAINER}" ldapsearch -x \
      -H "ldap://127.0.0.1:10389" \
      -D "${LDAP_BIND_DN}" \
      -w "${LDAP_BIND_PASSWORD}" \
      -b "${LDAP_BASE_DN}" \
      "(objectClass=*)" dn >/dev/null 2>&1; then
      log "local ldap ready: ${LDAP_URL}"
      return
    fi
    sleep 1
  done

  echo "[p901-ldap-smoke][FAIL] local ldap failed to start"
  exit 1
}

main() {
  ensure_ldap
  ensure_backend

  local code body admin_access provider_name provider_id ldap_access
  provider_name="ldap-live-$(date +%s)"

  body="${TMP_DIR}/admin-login-local.json"
  code="$(request_json POST /api/v1/auth/login '{"username":"admin","password":"123456","provider":"local"}' "" no "${body}")"
  expect_status "${code}" "200" "admin login local"
  admin_access="$(json_field "${body}" "accessToken")"

  # 映射本地账号（LDAP 校验通过后，按同名本地用户签发 JWT）
  body="${TMP_DIR}/create-local-mapped-user.json"
  code="$(request_json POST /api/v1/auth/users "{\"username\":\"${LDAP_TEST_USER}\",\"password\":\"123456\",\"role\":\"readonly\",\"allowedNamespaces\":[\"dev\"]}" "${admin_access}" yes "${body}")"
  if [[ "${code}" != "201" && "${code}" != "409" ]]; then
    echo "[p901-ldap-smoke][FAIL] create mapped local user expected 201/409 actual=${code}"
    exit 1
  fi

  body="${TMP_DIR}/create-ldap-provider.json"
  code="$(request_json POST /api/v1/auth/providers "{\"name\":\"${provider_name}\",\"type\":\"ldap\",\"config\":{\"url\":\"${LDAP_URL}\",\"baseDN\":\"${LDAP_BASE_DN}\",\"bindDN\":\"${LDAP_BIND_DN}\",\"bindPassword\":\"${LDAP_BIND_PASSWORD}\",\"loginAttr\":\"${LDAP_LOGIN_ATTR}\",\"userFilter\":\"${LDAP_USER_FILTER}\",\"timeoutSeconds\":\"8\"}}" "${admin_access}" yes "${body}")"
  expect_status "${code}" "201" "create ldap provider"

  body="${TMP_DIR}/list-providers.json"
  code="$(request_json GET /api/v1/auth/providers "" "${admin_access}" no "${body}")"
  expect_status "${code}" "200" "list providers"
  provider_id="$(json_provider_id "${body}" "${provider_name}")"
  if [[ -z "${provider_id}" ]]; then
    echo "[p901-ldap-smoke][FAIL] ldap provider id not found"
    exit 1
  fi

  body="${TMP_DIR}/set-default-ldap.json"
  code="$(request_json POST "/api/v1/auth/providers/${provider_id}/default" "" "${admin_access}" yes "${body}")"
  expect_status "${code}" "204" "set default ldap provider"

  body="${TMP_DIR}/ldap-login-success.json"
  code="$(request_json POST /api/v1/auth/login "{\"username\":\"${LDAP_TEST_USER}\",\"password\":\"${LDAP_TEST_PASSWORD}\"}" "" no "${body}")"
  expect_status "${code}" "200" "ldap login success"
  ldap_access="$(json_field "${body}" "accessToken")"
  if [[ -z "${ldap_access}" ]]; then
    echo "[p901-ldap-smoke][FAIL] ldap login missing access token"
    exit 1
  fi

  body="${TMP_DIR}/ldap-login-bad-password.json"
  code="$(request_json POST /api/v1/auth/login "{\"username\":\"${LDAP_TEST_USER}\",\"password\":\"wrong-password\"}" "" no "${body}")"
  expect_status "${code}" "401" "ldap login with wrong password should fail"

  log "all p901 ldap bind smoke checks passed"
}

main "$@"
