import { create } from "zustand";

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

type ResourceState = {
  services: ServiceItem[];
  configMaps: ConfigMapItem[];
  secrets: SecretItem[];
  loading: boolean;
  error: string;
  load: () => Promise<void>;
};

export const useResourceStore = create<ResourceState>((set) => ({
  services: [],
  configMaps: [],
  secrets: [],
  loading: false,
  error: "",
  load: async () => {
    set({ loading: true, error: "" });
    try {
      const [sResp, cResp, secResp] = await Promise.all([
        fetch("/api/v1/services"),
        fetch("/api/v1/configmaps"),
        fetch("/api/v1/secrets")
      ]);
      if (!sResp.ok || !cResp.ok || !secResp.ok) {
        throw new Error("加载服务与配置资源失败");
      }
      const sData = (await sResp.json()) as { items: ServiceItem[] };
      const cData = (await cResp.json()) as { items: ConfigMapItem[] };
      const secData = (await secResp.json()) as { items: SecretItem[] };
      set({
        services: sData.items,
        configMaps: cData.items,
        secrets: secData.items
      });
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "加载失败" });
    } finally {
      set({ loading: false });
    }
  }
}));
