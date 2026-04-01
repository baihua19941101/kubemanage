const ROLE_KEY = "km_user_role";
const USER_KEY = "km_user_name";
const ACCESS_TOKEN_KEY = "km_access_token";
const REFRESH_TOKEN_KEY = "km_refresh_token";

export function getCurrentRole(): string {
  return localStorage.getItem(ROLE_KEY) || "admin";
}

export function setCurrentRole(role: string) {
  localStorage.setItem(ROLE_KEY, role);
}

export function getCurrentUser(): string {
  return localStorage.getItem(USER_KEY) || "demo-user";
}

export function setCurrentUser(user: string) {
  localStorage.setItem(USER_KEY, user);
}

export function getAccessToken(): string {
  return localStorage.getItem(ACCESS_TOKEN_KEY) || "";
}

export function getRefreshToken(): string {
  return localStorage.getItem(REFRESH_TOKEN_KEY) || "";
}

export function setAuthTokens(accessToken: string, refreshToken: string) {
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
}

export function clearAuthTokens() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}

export class ApiRequestError extends Error {
  status: number;
  requestId?: string;
  code?: string;
  hint?: string;

  constructor(message: string, status: number, requestId?: string, code?: string, hint?: string) {
    super(message);
    this.name = "ApiRequestError";
    this.status = status;
    this.requestId = requestId;
    this.code = code;
    this.hint = hint;
  }
}

export async function apiFetch(input: RequestInfo | URL, init?: RequestInit) {
  const headers = new Headers(init?.headers || {});
  const accessToken = getAccessToken();
  if (accessToken) {
    headers.set("Authorization", `Bearer ${accessToken}`);
  } else {
    headers.set("X-User-Role", getCurrentRole());
    headers.set("X-User", getCurrentUser());
  }
  return fetch(input, {
    ...init,
    headers
  });
}

export async function parseApiError(resp: Response, fallbackMessage: string): Promise<ApiRequestError> {
  const requestIdFromHeader = resp.headers.get("X-Request-Id") || undefined;
  let message = fallbackMessage;
  let code: string | undefined;
  let hint: string | undefined;
  let requestId = requestIdFromHeader;

  const contentType = resp.headers.get("Content-Type") || "";
  if (contentType.includes("application/json")) {
    const body = (await resp.json().catch(() => null)) as { error?: string; code?: string; hint?: string; requestId?: string } | null;
    if (body?.error) {
      message = body.error;
    }
    code = body?.code;
    hint = body?.hint;
    requestId = body?.requestId || requestId;
  } else {
    const text = (await resp.text().catch(() => "")).trim();
    if (text) {
      message = text;
    }
  }

  if (requestId) {
    message = `${message}（requestId: ${requestId}）`;
  }
  return new ApiRequestError(message, resp.status, requestId, code, hint);
}
