import { Alert, Button, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import YamlDialog from "../components/framework/YamlDialog";
import { useAuthStore } from "../stores/useAuthStore";
import { useNamespaceStore } from "../stores/useNamespaceStore";

type NamespaceRow = {
  name: string;
  status: string;
  age: string;
};

export default function NamespacePage() {
  const items = useNamespaceStore((s) => s.items);
  const loading = useNamespaceStore((s) => s.loading);
  const error = useNamespaceStore((s) => s.error);
  const load = useNamespaceStore((s) => s.load);
  const create = useNamespaceStore((s) => s.create);
  const remove = useNamespaceStore((s) => s.remove);
  const fetchYaml = useNamespaceStore((s) => s.fetchYaml);
  const canNamespaceWrite = useAuthStore((s) => s.canNamespaceWrite);
  const canWriteNamespace = useAuthStore((s) => s.canWriteNamespace);
  const allowedNamespaces = useAuthStore((s) => s.allowedNamespaces);

  const [keyword, setKeyword] = useState("");
  const [createName, setCreateName] = useState("");
  const [selected, setSelected] = useState<NamespaceRow | null>(null);
  const [yamlOpen, setYamlOpen] = useState(false);
  const [yamlText, setYamlText] = useState("");

  useEffect(() => {
    void load();
  }, [load]);

  const filtered = useMemo(
    () =>
      items.filter((n) =>
        n.name.toLowerCase().includes(keyword.toLowerCase().trim())
      ),
    [items, keyword]
  );

  const columns = [
    { key: "name", header: "名称", render: (row: NamespaceRow) => row.name },
    { key: "status", header: "状态", render: (row: NamespaceRow) => row.status },
    { key: "age", header: "Age", render: (row: NamespaceRow) => row.age }
  ];

  const normalizedCreateName = createName.trim();
  const canCreateByScope = normalizedCreateName.length > 0 && canWriteNamespace(normalizedCreateName);
  const createDeniedByScope = normalizedCreateName.length > 0 && !canWriteNamespace(normalizedCreateName);
  const scopeHint =
    allowedNamespaces().length === 0
      ? "当前账号无可写命名空间"
      : allowedNamespaces()[0] === "*"
        ? "当前账号可写全部命名空间"
        : `可写命名空间：${allowedNamespaces().join(", ")}`;

  return (
    <>
      <PageScaffold
        title="名称空间管理"
        description="支持名称空间创建、删除与 YAML 查看下载"
        actions={
          <Stack direction="row" spacing={1}>
            <TextField
              size="small"
              placeholder="新建名称空间"
              value={createName}
              onChange={(e) => setCreateName(e.target.value)}
              error={createDeniedByScope}
              helperText={createDeniedByScope ? `无权限创建：${normalizedCreateName}` : scopeHint}
              sx={{ width: 220 }}
            />
            <Button
              variant="contained"
              disabled={!canNamespaceWrite() || !canCreateByScope}
              onClick={async () => {
                await create(createName);
                setCreateName("");
              }}
            >
              创建
            </Button>
          </Stack>
        }
        toolbar={
          <Stack direction="row" spacing={1.5}>
            <TextField
              size="small"
              label="按名称筛选"
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              sx={{ width: 260 }}
            />
          </Stack>
        }
      >
        {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}
        <ResourceTable
          loading={loading}
          rows={filtered}
          rowKey={(r) => r.name}
          columns={columns}
          onRowClick={(row) => setSelected(row)}
        />
      </PageScaffold>

      <DetailDrawer
        open={Boolean(selected)}
        title={selected ? `名称空间详情 - ${selected.name}` : "名称空间详情"}
        onClose={() => setSelected(null)}
        actions={
          selected ? (
            <Stack direction="row" spacing={1}>
              <Button
                size="small"
                onClick={async () => {
                  if (!selected) return;
                  const text = await fetchYaml(selected.name);
                  setYamlText(text);
                  setYamlOpen(true);
                }}
              >
                查看 YAML
              </Button>
              <Button
                size="small"
                component="a"
                href={selected ? `/api/v1/namespaces/${selected.name}/yaml/download` : "#"}
              >
                下载 YAML
              </Button>
              <Button
                size="small"
                color="error"
                disabled={!canNamespaceWrite() || !canWriteNamespace(selected.name)}
                onClick={async () => {
                  if (!selected) return;
                  if (window.confirm(`确认删除名称空间 ${selected.name} ?`)) {
                    await remove(selected.name);
                    setSelected(null);
                  }
                }}
              >
                删除
              </Button>
            </Stack>
          ) : null
        }
      >
        {selected && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selected.name}</Typography>
            <Typography variant="body2">状态：{selected.status}</Typography>
            <Typography variant="body2">Age：{selected.age}</Typography>
          </Stack>
        )}
      </DetailDrawer>

      <YamlDialog
        open={yamlOpen}
        title={selected ? `Namespace YAML - ${selected.name}` : "Namespace YAML"}
        yaml={yamlText}
        onClose={() => setYamlOpen(false)}
      />
    </>
  );
}
