import { create } from "zustand";
import { apiFetch } from "../lib/api";

type Cluster = {
  name: string;
  version: string;
  status: string;
  nodes: number;
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

type LiveCluster = {
  name: string;
  version: string;
  status: string;
  nodes: number;
  apiServer: string;
  source: string;
};

type LiveNamespace = {
  name: string;
  status: string;
  age: string;
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
  liveCluster: LiveCluster | null;
  liveNamespaces: LiveNamespace[];
  loading: boolean;
  switching: string;
  error: string;
  load: () => Promise<void>;
  loadConnections: () => Promise<void>;
  loadLiveData: () => Promise<void>;
  switchCluster: (name: string) => Promise<void>;
  importKubeconfig: (payload: { name: string; kubeconfigContent: string }) => Promise<boolean>;
  importToken: (payload: { name: string; apiServer: string; bearerToken: string; caCert: string; skipTlsVerify: boolean }) => Promise<boolean>;
  testConnection: (payload: { mode: string; apiServer?: string; kubeconfigContent?: string; bearerToken?: string; caCert?: string; skipTlsVerify?: boolean }) => Promise<ConnectionTestResult>;
  activateConnection: (id: number) => Promise<boolean>;
};

export const useClusterStore = create<ClusterState>((set, get) => ({
  clusters: [],
  current: "",
  connections: [],
  liveCluster: null,
  liveNamespaces: [],
  loading: false,
  switching: "",
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
  loadLiveData: async () => {
    try {
      const [clusterResp, nsResp] = await Promise.all([
        apiFetch("/api/v1/clusters/live"),
        apiFetch("/api/v1/namespaces/live")
      ]);
      if (!clusterResp.ok || !nsResp.ok) {
        throw new Error("当前尚未激活真实集群连接");
      }
      const cluster = (await clusterResp.json()) as LiveCluster;
      const namespaces = (await nsResp.json()) as { items: LiveNamespace[] };
      set({ liveCluster: cluster, liveNamespaces: namespaces.items, error: "" });
    } catch (err) {
      set({ liveCluster: null, liveNamespaces: [], error: err instanceof Error ? err.message : "获取真实集群数据失败" });
    }
  },
  switchCluster: async (name: string) => {
    set({ switching: name, error: "" });
    try {
      const resp = await apiFetch("/api/v1/clusters/switch", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name })
      });
      if (!resp.ok) {
        throw new Error("切换示例集群失败");
      }
      set({ current: name });
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "切换示例集群失败" });
    } finally {
      set({ switching: "" });
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
    await Promise.all([get().loadConnections(), get().loadLiveData()]);
    return true;
  }
}));
