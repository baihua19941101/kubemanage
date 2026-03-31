const ROLE_KEY = "km_user_role";

export function getCurrentRole(): string {
  return localStorage.getItem(ROLE_KEY) || "admin";
}

export function setCurrentRole(role: string) {
  localStorage.setItem(ROLE_KEY, role);
}

export async function apiFetch(input: RequestInfo | URL, init?: RequestInit) {
  const headers = new Headers(init?.headers || {});
  headers.set("X-User-Role", getCurrentRole());
  headers.set("X-User", "demo-user");
  return fetch(input, {
    ...init,
    headers
  });
}
