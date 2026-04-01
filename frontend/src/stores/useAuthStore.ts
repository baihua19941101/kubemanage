import { create } from "zustand";
import { getCurrentRole, setCurrentRole } from "../lib/api";

type AuthState = {
  role: string;
  setRole: (role: string) => void;
  canClusterManage: () => boolean;
  canNamespaceWrite: () => boolean;
  canWorkloadWrite: () => boolean;
  canAuditRead: () => boolean;
  allowedNamespaces: () => string[];
};

export const useAuthStore = create<AuthState>((set, get) => ({
  role: getCurrentRole(),
  setRole: (role: string) => {
    setCurrentRole(role);
    set({ role });
  },
  canClusterManage: () => get().role === "admin",
  canNamespaceWrite: () => {
    const role = get().role;
    return role === "admin" || role === "operator";
  },
  canWorkloadWrite: () => {
    const role = get().role;
    return role === "admin" || role === "operator";
  },
  canAuditRead: () => get().role === "admin",
  allowedNamespaces: () => {
    const role = get().role;
    if (role === "admin") return ["*"];
    if (role === "operator") return ["dev"];
    return [];
  }
}));
