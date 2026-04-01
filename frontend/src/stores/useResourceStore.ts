import { create } from "zustand";
import { apiFetch } from "../lib/api";

type ServiceItem = {
  name: string;
  namespace: string;
  type: string;
  clusterIP: string;
  ports: string;
  pods: number;
  age: string;
};

type ConfigMapItem = {
  name: string;
  namespace: string;
  dataCount: number;
  age: string;
};

type SecretItem = {
 name: string;
 namespace: string;
 type: string;
 data: Record<string, string>;
 age: string;
};

type IngressItem = {
  name: string;
  namespace: string;
  className: string;
  hosts: string[];
  address: string;
  tls: boolean;
  age: string;
};

type HPAItem = {
  name: string;
  namespace: string;
  targetKind: string;
  targetName: string;
  minReplicas: number;
  maxReplicas: number;
  currentReplicas: number;
  targetCPUPercent: number;
  currentCPUPercent: number;
  age: string;
};

type HPATarget = {
  kind: string;
  name: string;
  namespace: string;
  currentReplicas: number;
  desiredReplicas: number;
};

type ResourceState = {
  services: ServiceItem[];
  configMaps: ConfigMapItem[];
  secrets: SecretItem[];
  ingresses: IngressItem[];
  hpas: HPAItem[];
  loading: boolean;
  error: string;
  load: () => Promise<void>;
  getIngressServices: (name: string) => Promise<ServiceItem[]>;
  getHPATarget: (name: string) => Promise<HPATarget>;
};

export const useResourceStore = create<ResourceState>((set) => ({
  services: [],
  configMaps: [],
  secrets: [],
  ingresses: [],
  hpas: [],
  loading: false,
  error: "",
  load: async () => {
    set({ loading: true, error: "" });
    try {
      const [sResp, cResp, secResp, ingResp, hpaResp] = await Promise.all([
        apiFetch("/api/v1/services"),
        apiFetch("/api/v1/configmaps"),
        apiFetch("/api/v1/secrets"),
        apiFetch("/api/v1/ingresses"),
        apiFetch("/api/v1/hpas")
      ]);
      if (!sResp.ok || !cResp.ok || !secResp.ok || !ingResp.ok || !hpaResp.ok) {
        throw new Error("加载服务与配置资源失败");
      }
      const sData = (await sResp.json()) as { items: ServiceItem[] };
      const cData = (await cResp.json()) as { items: ConfigMapItem[] };
      const secData = (await secResp.json()) as { items: SecretItem[] };
      const ingData = (await ingResp.json()) as { items: IngressItem[] };
      const hpaData = (await hpaResp.json()) as { items: HPAItem[] };
      set({
        services: sData.items,
        configMaps: cData.items,
        secrets: secData.items,
        ingresses: ingData.items,
        hpas: hpaData.items
      });
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "加载失败" });
    } finally {
      set({ loading: false });
    }
  },
  getIngressServices: async (name: string) => {
    const resp = await apiFetch(`/api/v1/ingresses/${name}/services`);
    if (!resp.ok) {
      throw new Error("获取 Ingress 关联服务失败");
    }
    const data = (await resp.json()) as { items: ServiceItem[] };
    return data.items;
  },
  getHPATarget: async (name: string) => {
    const resp = await apiFetch(`/api/v1/hpas/${name}/target`);
    if (!resp.ok) {
      throw new Error("获取 HPA 目标失败");
    }
    return (await resp.json()) as HPATarget;
  }
}));
