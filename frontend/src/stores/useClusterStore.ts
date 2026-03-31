import { create } from "zustand";

type Cluster = {
  name: string;
  version: string;
  status: string;
  nodes: number;
};

type ClusterState = {
  clusters: Cluster[];
  current: string;
  loading: boolean;
  switching: string;
  load: () => Promise<void>;
  switchCluster: (name: string) => Promise<void>;
};

export const useClusterStore = create<ClusterState>((set) => ({
  clusters: [],
  current: "",
  loading: false,
  switching: "",
  load: async () => {
    set({ loading: true });
    try {
      const resp = await fetch("/api/v1/clusters");
      if (!resp.ok) {
        throw new Error("fetch clusters failed");
      }
      const data = (await resp.json()) as { items: Cluster[]; current: string };
      set({ clusters: data.items, current: data.current });
    } finally {
      set({ loading: false });
    }
  },
  switchCluster: async (name: string) => {
    set({ switching: name });
    try {
      const resp = await fetch("/api/v1/clusters/switch", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({ name })
      });
      if (!resp.ok) {
        throw new Error("switch cluster failed");
      }
      set({ current: name });
    } finally {
      set({ switching: "" });
    }
  }
}));
