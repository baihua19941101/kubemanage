import { create } from "zustand";
import { getCurrentRole, setCurrentRole } from "../lib/api";

type AuthState = {
  role: string;
  setRole: (role: string) => void;
  canNamespaceWrite: () => boolean;
  canWorkloadWrite: () => boolean;
  canAuditRead: () => boolean;
};

export const useAuthStore = create<AuthState>((set, get) => ({
  role: getCurrentRole(),
  setRole: (role: string) => {
    setCurrentRole(role);
    set({ role });
  },
  canNamespaceWrite: () => {
    const role = get().role;
    return role === "admin" || role === "operator";
  },
  canWorkloadWrite: () => {
    const role = get().role;
    return role === "admin" || role === "operator";
  },
  canAuditRead: () => get().role === "admin"
}));
