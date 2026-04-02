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

type PVItem = {
  name: string;
  capacity: string;
  accessModes: string;
  reclaimPolicy: string;
  status: string;
  claimRef: string;
  storageClass: string;
  age: string;
};

type PVCItem = {
  name: string;
  namespace: string;
  status: string;
  volume: string;
  capacity: string;
  accessModes: string;
  storageClass: string;
  age: string;
};

type StorageClassItem = {
  name: string;
  provisioner: string;
  reclaimPolicy: string;
  volumeBindingMode: string;
  allowVolumeExpansion: boolean;
  age: string;
};

type LimitRangeItem = {
  name: string;
  namespace: string;
  limitsCount: number;
  defaultCpu: string;
  defaultMemory: string;
  age: string;
};

type ResourceQuotaItem = {
  name: string;
  namespace: string;
  hardPods: string;
  usedPods: string;
  hardCpu: string;
  usedCpu: string;
  hardMemory: string;
  usedMemory: string;
  hardPvcs: string;
  usedPvcs: string;
  age: string;
};

type NetworkPolicyItem = {
  name: string;
  namespace: string;
  podSelector: string;
  policyTypes: string;
  ingressRules: number;
  egressRules: number;
  age: string;
};

type ResourceState = {
  services: ServiceItem[];
  configMaps: ConfigMapItem[];
  secrets: SecretItem[];
  ingresses: IngressItem[];
  hpas: HPAItem[];
  pvs: PVItem[];
  pvcs: PVCItem[];
  storageClasses: StorageClassItem[];
  limitRanges: LimitRangeItem[];
  resourceQuotas: ResourceQuotaItem[];
  networkPolicies: NetworkPolicyItem[];
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
  pvs: [],
  pvcs: [],
  storageClasses: [],
  limitRanges: [],
  resourceQuotas: [],
  networkPolicies: [],
  loading: false,
  error: "",
  load: async () => {
    set({ loading: true, error: "" });
    try {
      const [sResp, cResp, secResp, ingResp, hpaResp, pvResp, pvcResp, scResp, lrResp, rqResp, npResp] = await Promise.all([
        apiFetch("/api/v1/services"),
        apiFetch("/api/v1/configmaps"),
        apiFetch("/api/v1/secrets"),
        apiFetch("/api/v1/ingresses"),
        apiFetch("/api/v1/hpas"),
        apiFetch("/api/v1/pvs"),
        apiFetch("/api/v1/pvcs"),
        apiFetch("/api/v1/storageclasses"),
        apiFetch("/api/v1/limitranges"),
        apiFetch("/api/v1/resourcequotas"),
        apiFetch("/api/v1/networkpolicies")
      ]);
      if (
        !sResp.ok ||
        !cResp.ok ||
        !secResp.ok ||
        !ingResp.ok ||
        !hpaResp.ok ||
        !pvResp.ok ||
        !pvcResp.ok ||
        !scResp.ok ||
        !lrResp.ok ||
        !rqResp.ok ||
        !npResp.ok
      ) {
        throw new Error("加载服务与配置资源失败");
      }
      const sData = (await sResp.json()) as { items: ServiceItem[] };
      const cData = (await cResp.json()) as { items: ConfigMapItem[] };
      const secData = (await secResp.json()) as { items: SecretItem[] };
      const ingData = (await ingResp.json()) as { items: IngressItem[] };
      const hpaData = (await hpaResp.json()) as { items: HPAItem[] };
      const pvData = (await pvResp.json()) as { items: PVItem[] };
      const pvcData = (await pvcResp.json()) as { items: PVCItem[] };
      const scData = (await scResp.json()) as { items: StorageClassItem[] };
      const lrData = (await lrResp.json()) as { items: LimitRangeItem[] };
      const rqData = (await rqResp.json()) as { items: ResourceQuotaItem[] };
      const npData = (await npResp.json()) as { items: NetworkPolicyItem[] };
      set({
        services: sData.items,
        configMaps: cData.items,
        secrets: secData.items,
        ingresses: ingData.items,
        hpas: hpaData.items,
        pvs: pvData.items,
        pvcs: pvcData.items,
        storageClasses: scData.items,
        limitRanges: lrData.items,
        resourceQuotas: rqData.items,
        networkPolicies: npData.items
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
