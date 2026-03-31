import { create } from "zustand";

type Deployment = {
  name: string;
  namespace: string;
  image: string;
  replicas: number;
  ready: number;
  age: string;
};

type Pod = {
  name: string;
  namespace: string;
  node: string;
  status: string;
  restarts: number;
  ip: string;
  image: string;
  age: string;
};

type WorkloadState = {
  deployments: Deployment[];
  pods: Pod[];
  loading: boolean;
  error: string;
  load: () => Promise<void>;
  getDeploymentYAML: (name: string) => Promise<string>;
  saveDeploymentYAML: (name: string, yaml: string) => Promise<boolean>;
  getPodYAML: (name: string) => Promise<string>;
  savePodYAML: (name: string, yaml: string) => Promise<boolean>;
  getPodLogs: (name: string) => Promise<string>;
};

export const useWorkloadStore = create<WorkloadState>((set) => ({
  deployments: [],
  pods: [],
  loading: false,
  error: "",
  load: async () => {
    set({ loading: true, error: "" });
    try {
      const [dResp, pResp] = await Promise.all([
        fetch("/api/v1/deployments"),
        fetch("/api/v1/pods")
      ]);
      if (!dResp.ok || !pResp.ok) {
        throw new Error("加载工作负载失败");
      }
      const dData = (await dResp.json()) as { items: Deployment[] };
      const pData = (await pResp.json()) as { items: Pod[] };
      set({ deployments: dData.items, pods: pData.items });
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "加载工作负载失败" });
    } finally {
      set({ loading: false });
    }
  },
  getDeploymentYAML: async (name: string) => {
    const resp = await fetch(`/api/v1/deployments/${name}/yaml`);
    if (!resp.ok) {
      throw new Error("获取 Deployment YAML 失败");
    }
    return resp.text();
  },
  saveDeploymentYAML: async (name: string, yaml: string) => {
    const resp = await fetch(`/api/v1/deployments/${name}/yaml`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ yaml })
    });
    return resp.ok;
  },
  getPodYAML: async (name: string) => {
    const resp = await fetch(`/api/v1/pods/${name}/yaml`);
    if (!resp.ok) {
      throw new Error("获取 Pod YAML 失败");
    }
    return resp.text();
  },
  savePodYAML: async (name: string, yaml: string) => {
    const resp = await fetch(`/api/v1/pods/${name}/yaml`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ yaml })
    });
    return resp.ok;
  },
  getPodLogs: async (name: string) => {
    const resp = await fetch(`/api/v1/pods/${name}/logs`);
    if (!resp.ok) {
      throw new Error("获取 Pod 日志失败");
    }
    return resp.text();
  }
}));
