import { create } from "zustand";
import {
  ApiRequestError,
  apiFetch,
  clearAuthTokens,
  getAccessToken,
  getCurrentRole,
  getCurrentUser,
  getRefreshToken,
  parseApiError,
  setAuthTokens,
  setCurrentRole,
  setCurrentUser
} from "../lib/api";

type LoginResult = { ok: boolean; error?: string };

type AuthState = {
  role: string;
  user: string;
  accessToken: string;
  refreshToken: string;
  namespaces: string[];
  authenticated: boolean;
  bootstrap: () => Promise<void>;
  login: (username: string, password: string, provider?: string) => Promise<LoginResult>;
  logout: () => Promise<void>;
  setRole: (role: string) => void;
  canClusterManage: () => boolean;
  canNamespaceWrite: () => boolean;
  canWorkloadWrite: () => boolean;
  canAuditRead: () => boolean;
  canUserManage: () => boolean;
  canWriteNamespace: (namespace: string) => boolean;
  allowedNamespaces: () => string[];
};

type LoginResponse = {
  accessToken: string;
  refreshToken: string;
  tokenType: string;
  expiresIn: number;
  user: string;
  role: string;
  allowedNamespaces?: string[];
};

type RefreshResponse = LoginResponse;

type MeResponse = {
  user: string;
  role: string;
  permissions: string;
  allowedNamespaces: string[];
};

function normalizeRole(role: string) {
  if (role === "admin") return "admin";
  if (role === "standard-user" || role === "operator") return "standard-user";
  return "readonly";
}

function normalizeAllowedNamespaces(role: string, namespaces?: string[]) {
  const normalizedRole = normalizeRole(role);
  if (normalizedRole === "admin") {
    return ["*"];
  }
  const values = (namespaces || [])
    .map((item) => item.trim())
    .filter(Boolean);
  if (values.length > 0) {
    return Array.from(new Set(values));
  }
  if (normalizedRole === "standard-user") {
    return ["dev"];
  }
  return [];
}

export const useAuthStore = create<AuthState>((set, get) => ({
  role: normalizeRole(getCurrentRole()),
  user: getCurrentUser(),
  accessToken: getAccessToken(),
  refreshToken: getRefreshToken(),
  namespaces: normalizeRole(getCurrentRole()) === "admin" ? ["*"] : normalizeRole(getCurrentRole()) === "standard-user" ? ["dev"] : [],
  authenticated: Boolean(getAccessToken()),
  bootstrap: async () => {
    const accessToken = getAccessToken();
    const refreshToken = getRefreshToken();
    if (!accessToken) {
      set({
        role: normalizeRole(getCurrentRole()),
        user: getCurrentUser(),
        accessToken: "",
        refreshToken,
        namespaces: normalizeAllowedNamespaces(getCurrentRole()),
        authenticated: false
      });
      return;
    }
    const meResp = await apiFetch("/api/v1/auth/me");
    if (meResp.ok) {
      const me = (await meResp.json()) as MeResponse;
      const role = normalizeRole(me.role);
      const namespaces = normalizeAllowedNamespaces(role, me.allowedNamespaces);
      setCurrentRole(role);
      setCurrentUser(me.user);
      set({ role, user: me.user, accessToken, refreshToken, namespaces, authenticated: true });
      return;
    }

    if (!refreshToken) {
      clearAuthTokens();
      setCurrentRole("readonly");
      setCurrentUser("demo-user");
      set({ role: "readonly", user: "demo-user", accessToken: "", refreshToken: "", namespaces: [], authenticated: false });
      return;
    }

    const refreshResp = await apiFetch("/api/v1/auth/refresh", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refreshToken })
    });
    if (!refreshResp.ok) {
      clearAuthTokens();
      setCurrentRole("readonly");
      setCurrentUser("demo-user");
      set({ role: "readonly", user: "demo-user", accessToken: "", refreshToken: "", authenticated: false });
      return;
    }
    const refreshed = (await refreshResp.json()) as RefreshResponse;
    setAuthTokens(refreshed.accessToken, refreshed.refreshToken);
    const role = normalizeRole(refreshed.role);
    const namespaces = normalizeAllowedNamespaces(role, refreshed.allowedNamespaces);
    setCurrentRole(role);
    setCurrentUser(refreshed.user);
    set({
      role,
      user: refreshed.user,
      accessToken: refreshed.accessToken,
      refreshToken: refreshed.refreshToken,
      namespaces,
      authenticated: true
    });
  },
  login: async (username: string, password: string, provider = "local") => {
    const resp = await apiFetch("/api/v1/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password, provider })
    });
    if (!resp.ok) {
      const err = await parseApiError(resp, "登录失败");
      return { ok: false, error: err.message };
    }
    const payload = (await resp.json()) as LoginResponse;
    const role = normalizeRole(payload.role);
    const namespaces = normalizeAllowedNamespaces(role, payload.allowedNamespaces);
    setCurrentRole(role);
    setCurrentUser(payload.user);
    setAuthTokens(payload.accessToken, payload.refreshToken);
    set({
      role,
      user: payload.user,
      accessToken: payload.accessToken,
      refreshToken: payload.refreshToken,
      namespaces,
      authenticated: true
    });
    return { ok: true };
  },
  logout: async () => {
    const refreshToken = get().refreshToken || getRefreshToken();
    if (refreshToken) {
      const resp = await apiFetch("/api/v1/auth/logout", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refreshToken })
      });
      if (!resp.ok && resp.status !== 401) {
        const err = await parseApiError(resp, "退出登录失败");
        if (err instanceof ApiRequestError) {
          throw err;
        }
      }
    }
    clearAuthTokens();
    setCurrentRole("readonly");
    setCurrentUser("demo-user");
    set({ role: "readonly", user: "demo-user", accessToken: "", refreshToken: "", namespaces: [], authenticated: false });
  },
  setRole: (role: string) => {
    const normalized = normalizeRole(role);
    setCurrentRole(normalized);
    set({ role: normalized });
  },
  canClusterManage: () => get().role === "admin",
  canNamespaceWrite: () => {
    const role = get().role;
    return role === "admin" || role === "standard-user";
  },
  canWorkloadWrite: () => {
    const role = get().role;
    return role === "admin" || role === "standard-user";
  },
  canAuditRead: () => get().role === "admin",
  canUserManage: () => get().role === "admin",
  canWriteNamespace: (namespace: string) => {
    const role = get().role;
    if (role === "admin") {
      return true;
    }
    if (role !== "standard-user") {
      return false;
    }
    const target = namespace.trim();
    if (!target) {
      return false;
    }
    const scopes = get().namespaces;
    if (scopes.includes("*")) {
      return true;
    }
    return scopes.includes(target);
  },
  allowedNamespaces: () => {
    return get().namespaces;
  }
}));
