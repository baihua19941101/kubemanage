import {
  Box,
  Button,
  Chip,
  Stack,
  TextField,
  Typography
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import YamlDialog from "../components/framework/YamlDialog";
import { useClusterStore } from "../stores/useClusterStore";
import { useAuthStore } from "../stores/useAuthStore";

type ClusterRow = {
  name: string;
  version: string;
  status: string;
  nodes: number;
};

export default function ClusterPage() {
  const clusters = useClusterStore((s) => s.clusters);
  const current = useClusterStore((s) => s.current);
  const loading = useClusterStore((s) => s.loading);
  const switching = useClusterStore((s) => s.switching);
  const load = useClusterStore((s) => s.load);
  const switchCluster = useClusterStore((s) => s.switchCluster);
  const canWorkloadWrite = useAuthStore((s) => s.canWorkloadWrite);
  const [keyword, setKeyword] = useState("");
  const [selected, setSelected] = useState<ClusterRow | null>(null);
  const [yamlOpen, setYamlOpen] = useState(false);

  useEffect(() => {
    void load();
  }, [load]);

  const filtered = useMemo(
    () =>
      clusters.filter((c) =>
        c.name.toLowerCase().includes(keyword.toLowerCase().trim())
      ),
    [clusters, keyword]
  );

  const columns = [
    {
      key: "name",
      header: "名称",
      render: (row: ClusterRow) => row.name
    },
    {
      key: "version",
      header: "版本",
      render: (row: ClusterRow) => row.version
    },
    {
      key: "status",
      header: "状态",
      render: (row: ClusterRow) => (
        <Chip
          size="small"
          color={row.status === "ready" ? "success" : "default"}
          label={row.status}
        />
      )
    },
    {
      key: "nodes",
      header: "节点数",
      render: (row: ClusterRow) => row.nodes
    },
    {
      key: "current",
      header: "当前",
      render: (row: ClusterRow) =>
        current === row.name ? <Chip size="small" color="primary" label="当前集群" /> : "-"
    }
  ];

  const yamlText = selected
    ? `apiVersion: v1\nkind: Cluster\nmetadata:\n  name: ${selected.name}\nspec:\n  version: ${selected.version}\nstatus:\n  phase: ${selected.status}\n  nodes: ${selected.nodes}\n`
    : "";

  return (
    <>
      <PageScaffold
        title="集群管理"
        description="统一展示集群状态并执行集群切换操作"
        toolbar={
          <Stack direction={{ xs: "column", sm: "row" }} spacing={1.5}>
            <TextField
              size="small"
              label="按名称筛选"
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              sx={{ width: 260 }}
            />
            <Box sx={{ color: "text.secondary", alignSelf: "center" }}>
              共 {filtered.length} 条
            </Box>
          </Stack>
        }
      >
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
        title={selected ? `集群详情 - ${selected.name}` : "集群详情"}
        onClose={() => setSelected(null)}
        actions={
          selected ? (
            <Stack direction="row" spacing={1}>
              <Button
                size="small"
                variant="contained"
                disabled={
                  !canWorkloadWrite() ||
                  current === selected.name ||
                  switching === selected.name
                }
                onClick={() => {
                  if (selected) {
                    void switchCluster(selected.name);
                  }
                }}
              >
                {switching === selected.name ? "切换中..." : "切换到该集群"}
              </Button>
              <Button size="small" onClick={() => setYamlOpen(true)}>
                查看 YAML
              </Button>
            </Stack>
          ) : null
        }
      >
        {selected && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selected.name}</Typography>
            <Typography variant="body2">版本：{selected.version}</Typography>
            <Typography variant="body2">状态：{selected.status}</Typography>
            <Typography variant="body2">节点数：{selected.nodes}</Typography>
          </Stack>
        )}
      </DetailDrawer>

      <YamlDialog
        open={yamlOpen}
        title={selected ? `集群 YAML - ${selected.name}` : "集群 YAML"}
        yaml={yamlText}
        onClose={() => setYamlOpen(false)}
      />
    </>
  );
}
