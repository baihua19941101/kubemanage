import { create } from "zustand";
import { apiFetch, parseApiError } from "../lib/api";

type NamespaceItem = {
  name: string;
  status: string;
  age: string;
};

type NamespaceState = {
  items: NamespaceItem[];
  loading: boolean;
  error: string;
  load: () => Promise<void>;
  create: (name: string) => Promise<void>;
  remove: (name: string) => Promise<void>;
  fetchYaml: (name: string) => Promise<string>;
};

export const useNamespaceStore = create<NamespaceState>((set, get) => ({
  items: [],
  loading: false,
  error: "",
  load: async () => {
    set({ loading: true, error: "" });
    try {
      const resp = await apiFetch("/api/v1/namespaces");
      if (!resp.ok) {
        throw await parseApiError(resp, "获取名称空间列表失败");
      }
      const data = (await resp.json()) as { items: NamespaceItem[] };
      set({ items: data.items });
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "未知错误" });
    } finally {
      set({ loading: false });
    }
  },
  create: async (name: string) => {
    const trimmed = name.trim();
    if (!trimmed) {
      set({ error: "名称空间名称不能为空" });
      return;
    }
    set({ error: "" });
    try {
      const resp = await apiFetch("/api/v1/namespaces", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Action-Confirm": "CONFIRM"
        },
        body: JSON.stringify({ name: trimmed })
      });
      if (!resp.ok) {
        const err = await parseApiError(resp, "创建名称空间失败");
        set({ error: err.message });
        return;
      }
      await get().load();
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "创建名称空间失败" });
    }
  },
  remove: async (name: string) => {
    try {
      set({ error: "" });
      const resp = await apiFetch(`/api/v1/namespaces/${name}`, {
        method: "DELETE",
        headers: {
          "X-Action-Confirm": "CONFIRM"
        }
      });
      if (!resp.ok) {
        const err = await parseApiError(resp, "删除名称空间失败");
        set({ error: err.message });
        return;
      }
      await get().load();
    } catch (err) {
      set({ error: err instanceof Error ? err.message : "删除名称空间失败" });
    }
  },
  fetchYaml: async (name: string) => {
    const resp = await apiFetch(`/api/v1/namespaces/${name}/yaml`);
    if (!resp.ok) {
      throw await parseApiError(resp, "获取 YAML 失败");
    }
    return resp.text();
  }
}));
