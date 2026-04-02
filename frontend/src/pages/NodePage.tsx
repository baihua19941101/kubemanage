import { Alert, Button, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import YamlDialog from "../components/framework/YamlDialog";
import { apiFetch, parseApiError } from "../lib/api";

type NodeItem = {
  name: string;
  roles: string;
  version: string;
  internalIP: string;
  status: string;
  osImage: string;
  cpu: string;
  memory: string;
  podCount: number;
  labelsCount: number;
  taintsCount: number;
  age: string;
};

export default function NodePage() {
  const [items, setItems] = useState<NodeItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [keyword, setKeyword] = useState("");
  const [selected, setSelected] = useState<NodeItem | null>(null);
  const [yamlOpen, setYamlOpen] = useState(false);
  const [yamlTitle, setYamlTitle] = useState("");
  const [yamlText, setYamlText] = useState("");

  const columns = useMemo(
    () => [
      { key: "name", header: "名称", render: (r: NodeItem) => r.name },
      { key: "roles", header: "角色", render: (r: NodeItem) => r.roles || "<none>" },
      { key: "version", header: "Kubelet", render: (r: NodeItem) => r.version || "-" },
      { key: "ip", header: "Internal IP", render: (r: NodeItem) => r.internalIP || "-" },
      { key: "status", header: "状态", render: (r: NodeItem) => r.status || "-" },
      { key: "os", header: "OS", render: (r: NodeItem) => r.osImage || "-" },
      { key: "cpu", header: "CPU", render: (r: NodeItem) => r.cpu || "-" },
      { key: "memory", header: "Memory", render: (r: NodeItem) => r.memory || "-" },
      { key: "pods", header: "Pods", render: (r: NodeItem) => r.podCount },
      { key: "age", header: "Age", render: (r: NodeItem) => r.age || "-" }
    ],
    []
  );

  const filtered = useMemo(() => {
    const q = keyword.toLowerCase().trim();
    if (!q) return items;
    return items.filter(
      (item) =>
        item.name.toLowerCase().includes(q) ||
        item.roles.toLowerCase().includes(q) ||
        item.internalIP.toLowerCase().includes(q) ||
        item.status.toLowerCase().includes(q)
    );
  }, [items, keyword]);

  async function load() {
    setLoading(true);
    setError("");
    try {
      const resp = await apiFetch("/api/v1/nodes");
      if (!resp.ok) {
        throw await parseApiError(resp, "加载节点列表失败");
      }
      const data = (await resp.json()) as { items: NodeItem[] };
      setItems(data.items);
    } catch (err) {
      setError(err instanceof Error ? err.message : "加载节点列表失败");
    } finally {
      setLoading(false);
    }
  }

  async function openYaml(nodeName: string) {
    setError("");
    try {
      const resp = await apiFetch(`/api/v1/nodes/${encodeURIComponent(nodeName)}/yaml`);
      if (!resp.ok) {
        throw await parseApiError(resp, "获取节点 YAML 失败");
      }
      const text = await resp.text();
      setYamlTitle(`Node YAML - ${nodeName}`);
      setYamlText(text);
      setYamlOpen(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : "获取节点 YAML 失败");
    }
  }

  useEffect(() => {
    void load();
  }, []);

  return (
    <>
      <PageScaffold
        title="节点管理"
        description="支持节点列表、详情与 YAML 查看下载"
        actions={
          <Button variant="outlined" onClick={() => void load()}>
            刷新
          </Button>
        }
        toolbar={
          <TextField
            size="small"
            label="按名称/角色/IP/状态筛选"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            sx={{ width: 320 }}
          />
        }
      >
        {error && (
          <Alert severity="error" sx={{ m: 1.5 }}>
            {error}
          </Alert>
        )}
        <ResourceTable loading={loading} rows={filtered} rowKey={(r) => r.name} columns={columns} onRowClick={(r) => setSelected(r)} />
      </PageScaffold>

      <DetailDrawer
        open={Boolean(selected)}
        title={selected ? `节点详情 - ${selected.name}` : "节点详情"}
        onClose={() => setSelected(null)}
        actions={
          selected ? (
            <Stack direction="row" spacing={1}>
              <Button size="small" onClick={() => void openYaml(selected.name)}>
                查看 YAML
              </Button>
              <Button size="small" component="a" href={selected ? `/api/v1/nodes/${encodeURIComponent(selected.name)}/yaml/download` : "#"}>
                下载 YAML
              </Button>
            </Stack>
          ) : null
        }
      >
        {selected && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selected.name}</Typography>
            <Typography variant="body2">角色：{selected.roles || "<none>"}</Typography>
            <Typography variant="body2">状态：{selected.status}</Typography>
            <Typography variant="body2">Kubelet：{selected.version}</Typography>
            <Typography variant="body2">Internal IP：{selected.internalIP}</Typography>
            <Typography variant="body2">OS：{selected.osImage}</Typography>
            <Typography variant="body2">CPU：{selected.cpu}</Typography>
            <Typography variant="body2">Memory：{selected.memory}</Typography>
            <Typography variant="body2">Pods：{selected.podCount}</Typography>
            <Typography variant="body2">Labels：{selected.labelsCount}</Typography>
            <Typography variant="body2">Taints：{selected.taintsCount}</Typography>
            <Typography variant="body2">Age：{selected.age}</Typography>
          </Stack>
        )}
      </DetailDrawer>

      <YamlDialog open={yamlOpen} title={yamlTitle} yaml={yamlText} onClose={() => setYamlOpen(false)} />
    </>
  );
}

