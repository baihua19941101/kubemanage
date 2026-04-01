import { create } from "zustand";
import { apiFetch } from "../lib/api";

type Cluster = {
  state: string;
  name: string;
  provider: string;
  distro: string;
  kubernetesVersion: string;
  architecture: string;
  cpu: string;
  memory: string;
  pods: number;
};

type ClusterConnection = {
  id: number;
  name: string;
  mode: string;
  apiServer: string;
  skipTLSVerify: boolean;
  isDefault: boolean;
  status: string;
  lastCheckedAt?: string;
  lastError?: string;
  hasKubeconfig: boolean;
  hasBearerToken: boolean;
  hasCaCert: boolean;
};

type ConnectionTestResult = {
  success: boolean;
  version?: string;
  server?: string;
  nodeCount?: number;
  namespaceCount?: number;
  message: string;
};

type ClusterState = {
  clusters: Cluster[];
  current: string;
  connections: ClusterConnection[];
  loading: boolean;
  error: string;
  load: () => Promise<void>;
  loadConnections: () => Promise<void>;
  importKubeconfig: (payload: { name: string; kubeconfigContent: string }) => Promise<boolean>;
  importToken: (payload: { name: string; apiServer: string; bearerToken: string; caCert: string; skipTlsVerify: boolean }) => Promise<boolean>;
  testConnection: (payload: { mode: string; apiServer?: string; kubeconfigContent?: string; bearerToken?: string; caCert?: string; skipTlsVerify?: boolean }) => Promise<ConnectionTestResult>;
  activateConnection: (id: number) => Promise<boolean>;
};

export const useClusterStore = create<ClusterState>((set, get) => ({
  clusters: [],
  current: "",
  connections: [],
  loading: false,
  error: "",
  load: async () => {
    set({ loading: true, error: "" });
    try {
      const resp = await apiFetch("/api/v1/clusters");
      if (!resp.ok) {
        throw new Error("获取集群列表失败");
      }
      const data = (await resp.json()) as { items: Cluster[]; current: string };
      set({ clusters: data.items, current: data.current });
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "获取集群列表失败" });
    } finally {
      set({ loading: false });
    }
  },
  loadConnections: async () => {
    try {
      const resp = await apiFetch("/api/v1/clusters/connections");
      if (!resp.ok) {
        throw new Error("获取集群连接失败");
      }
      const data = (await resp.json()) as { items: ClusterConnection[] };
      set({ connections: data.items });
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "获取集群连接失败" });
    }
  },
  importKubeconfig: async (payload) => {
    set({ error: "" });
    const resp = await apiFetch("/api/v1/clusters/connections/import/kubeconfig", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (!resp.ok) {
      const body = await resp.json().catch(() => null) as { error?: string } | null;
      set({ error: body?.error || "导入 kubeconfig 失败" });
      return false;
    }
    await get().loadConnections();
    return true;
  },
  importToken: async (payload) => {
    set({ error: "" });
    const resp = await apiFetch("/api/v1/clusters/connections/import/token", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (!resp.ok) {
      const body = await resp.json().catch(() => null) as { error?: string } | null;
      set({ error: body?.error || "导入 token 集群失败" });
      return false;
    }
    await get().loadConnections();
    return true;
  },
  testConnection: async (payload) => {
    const resp = await apiFetch("/api/v1/clusters/connections/test", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (!resp.ok) {
      const body = await resp.json().catch(() => null) as { error?: string } | null;
      throw new Error(body?.error || "测试连接失败");
    }
    return resp.json() as Promise<ConnectionTestResult>;
  },
  activateConnection: async (id: number) => {
    set({ error: "" });
    const resp = await apiFetch(`/api/v1/clusters/connections/${id}/activate`, {
      method: "POST"
    });
    if (!resp.ok) {
      const body = await resp.json().catch(() => null) as { error?: string } | null;
      set({ error: body?.error || "激活真实集群失败" });
      return false;
    }
    await Promise.all([get().loadConnections(), get().load()]);
    return true;
  }
}));
