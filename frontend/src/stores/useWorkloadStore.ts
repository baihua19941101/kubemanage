import { create } from "zustand";
import { apiFetch } from "../lib/api";

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

type StatefulSet = {
  name: string;
  namespace: string;
  replicas: number;
  ready: number;
  service: string;
  image: string;
  age: string;
};

type DaemonSet = {
  name: string;
  namespace: string;
  desired: number;
  current: number;
  image: string;
  age: string;
};

type Job = {
  name: string;
  namespace: string;
  completions: number;
  failed: number;
  status: string;
  age: string;
};

type CronJob = {
  name: string;
  namespace: string;
  schedule: string;
  suspend: boolean;
  lastRun: string;
  age: string;
};

type WorkloadState = {
  deployments: Deployment[];
  pods: Pod[];
  statefulSets: StatefulSet[];
  daemonSets: DaemonSet[];
  jobs: Job[];
  cronJobs: CronJob[];
  loading: boolean;
  error: string;
  load: () => Promise<void>;
  getDeploymentYAML: (name: string) => Promise<string>;
  saveDeploymentYAML: (name: string, yaml: string) => Promise<boolean>;
  getPodYAML: (name: string) => Promise<string>;
  savePodYAML: (name: string, yaml: string) => Promise<boolean>;
  getPodLogs: (name: string) => Promise<string>;
  getStatefulSetYAML: (name: string) => Promise<string>;
  saveStatefulSetYAML: (name: string, yaml: string) => Promise<boolean>;
  getDaemonSetYAML: (name: string) => Promise<string>;
  saveDaemonSetYAML: (name: string, yaml: string) => Promise<boolean>;
  getJobYAML: (name: string) => Promise<string>;
  saveJobYAML: (name: string, yaml: string) => Promise<boolean>;
  getCronJobYAML: (name: string) => Promise<string>;
  saveCronJobYAML: (name: string, yaml: string) => Promise<boolean>;
};

export const useWorkloadStore = create<WorkloadState>((set) => ({
  deployments: [],
  pods: [],
  statefulSets: [],
  daemonSets: [],
  jobs: [],
  cronJobs: [],
  loading: false,
  error: "",
  load: async () => {
    set({ loading: true, error: "" });
    try {
      const [dResp, pResp, sResp, dsResp, jResp, cjResp] = await Promise.all([
        apiFetch("/api/v1/deployments"),
        apiFetch("/api/v1/pods"),
        apiFetch("/api/v1/statefulsets"),
        apiFetch("/api/v1/daemonsets"),
        apiFetch("/api/v1/jobs"),
        apiFetch("/api/v1/cronjobs")
      ]);
      if (!dResp.ok || !pResp.ok || !sResp.ok || !dsResp.ok || !jResp.ok || !cjResp.ok) {
        throw new Error("加载工作负载失败");
      }
      const dData = (await dResp.json()) as { items: Deployment[] };
      const pData = (await pResp.json()) as { items: Pod[] };
      const sData = (await sResp.json()) as { items: StatefulSet[] };
      const dsData = (await dsResp.json()) as { items: DaemonSet[] };
      const jData = (await jResp.json()) as { items: Job[] };
      const cjData = (await cjResp.json()) as { items: CronJob[] };
      set({
        deployments: dData.items,
        pods: pData.items,
        statefulSets: sData.items,
        daemonSets: dsData.items,
        jobs: jData.items,
        cronJobs: cjData.items
      });
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "加载工作负载失败" });
    } finally {
      set({ loading: false });
    }
  },
  getDeploymentYAML: async (name: string) => {
    const resp = await apiFetch(`/api/v1/deployments/${name}/yaml`);
    if (!resp.ok) {
      throw new Error("获取 Deployment YAML 失败");
    }
    return resp.text();
  },
  saveDeploymentYAML: async (name: string, yaml: string) => {
    const resp = await apiFetch(`/api/v1/deployments/${name}/yaml`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ yaml })
    });
    return resp.ok;
  },
  getPodYAML: async (name: string) => {
    const resp = await apiFetch(`/api/v1/pods/${name}/yaml`);
    if (!resp.ok) {
      throw new Error("获取 Pod YAML 失败");
    }
    return resp.text();
  },
  savePodYAML: async (name: string, yaml: string) => {
    const resp = await apiFetch(`/api/v1/pods/${name}/yaml`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ yaml })
    });
    return resp.ok;
  },
  getPodLogs: async (name: string) => {
    const resp = await apiFetch(`/api/v1/pods/${name}/logs`);
    if (!resp.ok) {
      throw new Error("获取 Pod 日志失败");
    }
    return resp.text();
  },
  getStatefulSetYAML: async (name: string) => {
    const resp = await apiFetch(`/api/v1/statefulsets/${name}/yaml`);
    if (!resp.ok) {
      throw new Error("获取 StatefulSet YAML 失败");
    }
    return resp.text();
  },
  saveStatefulSetYAML: async (name: string, yaml: string) => {
    const resp = await apiFetch(`/api/v1/statefulsets/${name}/yaml`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ yaml })
    });
    return resp.ok;
  },
  getDaemonSetYAML: async (name: string) => {
    const resp = await apiFetch(`/api/v1/daemonsets/${name}/yaml`);
    if (!resp.ok) {
      throw new Error("获取 DaemonSet YAML 失败");
    }
    return resp.text();
  },
  saveDaemonSetYAML: async (name: string, yaml: string) => {
    const resp = await apiFetch(`/api/v1/daemonsets/${name}/yaml`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ yaml })
    });
    return resp.ok;
  },
  getJobYAML: async (name: string) => {
    const resp = await apiFetch(`/api/v1/jobs/${name}/yaml`);
    if (!resp.ok) {
      throw new Error("获取 Job YAML 失败");
    }
    return resp.text();
  },
  saveJobYAML: async (name: string, yaml: string) => {
    const resp = await apiFetch(`/api/v1/jobs/${name}/yaml`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ yaml })
    });
    return resp.ok;
  },
  getCronJobYAML: async (name: string) => {
    const resp = await apiFetch(`/api/v1/cronjobs/${name}/yaml`);
    if (!resp.ok) {
      throw new Error("获取 CronJob YAML 失败");
    }
    return resp.text();
  },
  saveCronJobYAML: async (name: string, yaml: string) => {
    const resp = await apiFetch(`/api/v1/cronjobs/${name}/yaml`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ yaml })
    });
    return resp.ok;
  }
}));
