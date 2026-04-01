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
  authenticated: boolean;
  bootstrap: () => Promise<void>;
  login: (username: string, password: string) => Promise<LoginResult>;
  logout: () => Promise<void>;
  setRole: (role: string) => void;
  canClusterManage: () => boolean;
  canNamespaceWrite: () => boolean;
  canWorkloadWrite: () => boolean;
  canAuditRead: () => boolean;
  canUserManage: () => boolean;
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

export const useAuthStore = create<AuthState>((set, get) => ({
  role: normalizeRole(getCurrentRole()),
  user: getCurrentUser(),
  accessToken: getAccessToken(),
  refreshToken: getRefreshToken(),
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
        authenticated: false
      });
      return;
    }
    const meResp = await apiFetch("/api/v1/auth/me");
    if (meResp.ok) {
      const me = (await meResp.json()) as MeResponse;
      const role = normalizeRole(me.role);
      setCurrentRole(role);
      setCurrentUser(me.user);
      set({ role, user: me.user, accessToken, refreshToken, authenticated: true });
      return;
    }

    if (!refreshToken) {
      clearAuthTokens();
      setCurrentRole("readonly");
      setCurrentUser("demo-user");
      set({ role: "readonly", user: "demo-user", accessToken: "", refreshToken: "", authenticated: false });
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
    setCurrentRole(role);
    setCurrentUser(refreshed.user);
    set({
      role,
      user: refreshed.user,
      accessToken: refreshed.accessToken,
      refreshToken: refreshed.refreshToken,
      authenticated: true
    });
  },
  login: async (username: string, password: string) => {
    const resp = await apiFetch("/api/v1/auth/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password })
    });
    if (!resp.ok) {
      const err = await parseApiError(resp, "登录失败");
      return { ok: false, error: err.message };
    }
    const payload = (await resp.json()) as LoginResponse;
    const role = normalizeRole(payload.role);
    setCurrentRole(role);
    setCurrentUser(payload.user);
    setAuthTokens(payload.accessToken, payload.refreshToken);
    set({
      role,
      user: payload.user,
      accessToken: payload.accessToken,
      refreshToken: payload.refreshToken,
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
    set({ role: "readonly", user: "demo-user", accessToken: "", refreshToken: "", authenticated: false });
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
  allowedNamespaces: () => {
    const role = get().role;
    if (role === "admin") return ["*"];
    if (role === "standard-user") return ["dev"];
    return [];
  }
}));
