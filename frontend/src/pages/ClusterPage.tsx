import {
  Alert,
  Box,
  Button,
  Chip,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  MenuItem,
  Select,
  Stack,
  Switch,
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

type ConnectionRow = {
  id: number;
  name: string;
  mode: string;
  apiServer: string;
  status: string;
  isDefault: boolean;
  skipTLSVerify: boolean;
  hasKubeconfig: boolean;
  hasBearerToken: boolean;
  hasCaCert: boolean;
  lastCheckedAt?: string;
  lastError?: string;
};

type LiveNamespaceRow = {
  name: string;
  status: string;
  age: string;
};

type ImportMode = "kubeconfig" | "token";

export default function ClusterPage() {
  const clusters = useClusterStore((s) => s.clusters);
  const current = useClusterStore((s) => s.current);
  const loading = useClusterStore((s) => s.loading);
  const switching = useClusterStore((s) => s.switching);
  const connections = useClusterStore((s) => s.connections);
  const liveCluster = useClusterStore((s) => s.liveCluster);
  const liveNamespaces = useClusterStore((s) => s.liveNamespaces);
  const error = useClusterStore((s) => s.error);
  const load = useClusterStore((s) => s.load);
  const loadConnections = useClusterStore((s) => s.loadConnections);
  const loadLiveData = useClusterStore((s) => s.loadLiveData);
  const switchCluster = useClusterStore((s) => s.switchCluster);
  const importKubeconfig = useClusterStore((s) => s.importKubeconfig);
  const importToken = useClusterStore((s) => s.importToken);
  const testConnection = useClusterStore((s) => s.testConnection);
  const activateConnection = useClusterStore((s) => s.activateConnection);
  const canWorkloadWrite = useAuthStore((s) => s.canWorkloadWrite);
  const canClusterManage = useAuthStore((s) => s.canClusterManage);
  const role = useAuthStore((s) => s.role);

  const [keyword, setKeyword] = useState("");
  const [selected, setSelected] = useState<ClusterRow | null>(null);
  const [yamlOpen, setYamlOpen] = useState(false);
  const [importOpen, setImportOpen] = useState(false);
  const [importMode, setImportMode] = useState<ImportMode>("kubeconfig");
  const [connectionName, setConnectionName] = useState("");
  const [kubeconfigContent, setKubeconfigContent] = useState("");
  const [apiServer, setApiServer] = useState("");
  const [bearerToken, setBearerToken] = useState("");
  const [caCert, setCaCert] = useState("");
  const [skipTlsVerify, setSkipTlsVerify] = useState(false);
  const [submitLoading, setSubmitLoading] = useState(false);
  const [testLoading, setTestLoading] = useState(false);
  const [testResult, setTestResult] = useState<string>("");
  const [connectionFilter, setConnectionFilter] = useState("");

  useEffect(() => {
    void Promise.all([load(), loadConnections(), loadLiveData()]);
  }, [load, loadConnections, loadLiveData]);

  const filteredClusters = useMemo(
    () => clusters.filter((c) => c.name.toLowerCase().includes(keyword.toLowerCase().trim())),
    [clusters, keyword]
  );

  const filteredConnections = useMemo(
    () => connections.filter((c) => c.name.toLowerCase().includes(connectionFilter.toLowerCase().trim())),
    [connections, connectionFilter]
  );

  const clusterColumns = [
    { key: "name", header: "名称", render: (row: ClusterRow) => row.name },
    { key: "version", header: "版本", render: (row: ClusterRow) => row.version },
    {
      key: "status",
      header: "状态",
      render: (row: ClusterRow) => <Chip size="small" color={row.status === "ready" ? "success" : "default"} label={row.status} />
    },
    { key: "nodes", header: "节点数", render: (row: ClusterRow) => row.nodes },
    {
      key: "current",
      header: "当前",
      render: (row: ClusterRow) => (current === row.name ? <Chip size="small" color="primary" label="当前集群" /> : "-")
    }
  ];

  const connectionColumns = [
    { key: "name", header: "连接名称", render: (row: ConnectionRow) => row.name },
    { key: "mode", header: "模式", render: (row: ConnectionRow) => row.mode },
    { key: "server", header: "API Server", render: (row: ConnectionRow) => row.apiServer || "-" },
    {
      key: "status",
      header: "状态",
      render: (row: ConnectionRow) => <Chip size="small" color={row.status === "connected" ? "success" : row.status === "failed" ? "error" : "default"} label={row.status} />
    },
    {
      key: "active",
      header: "已激活",
      render: (row: ConnectionRow) => (row.isDefault ? <Chip size="small" color="primary" label="live" /> : "-")
    }
  ];

  const liveNamespaceColumns = [
    { key: "name", header: "名称", render: (row: LiveNamespaceRow) => row.name },
    { key: "status", header: "状态", render: (row: LiveNamespaceRow) => row.status },
    { key: "age", header: "Age", render: (row: LiveNamespaceRow) => row.age }
  ];

  const yamlText = selected
    ? `apiVersion: v1\nkind: Cluster\nmetadata:\n  name: ${selected.name}\nspec:\n  version: ${selected.version}\nstatus:\n  phase: ${selected.status}\n  nodes: ${selected.nodes}\n`
    : "";

  async function handleImport() {
    setSubmitLoading(true);
    setTestResult("");
    try {
      const ok =
        importMode === "kubeconfig"
          ? await importKubeconfig({ name: connectionName, kubeconfigContent })
          : await importToken({
              name: connectionName,
              apiServer,
              bearerToken,
              caCert,
              skipTlsVerify
            });
      if (ok) {
        setImportOpen(false);
        resetImportForm();
      }
    } finally {
      setSubmitLoading(false);
    }
  }

  async function handleTestConnection() {
    setTestLoading(true);
    setTestResult("");
    try {
      const result = await testConnection({
        mode: importMode,
        apiServer,
        kubeconfigContent,
        bearerToken,
        caCert,
        skipTlsVerify
      });
      setTestResult(`连接成功：${result.server || "unknown"}，版本 ${result.version || "unknown"}，节点 ${result.nodeCount ?? 0}`);
    } catch (err) {
      setTestResult(err instanceof Error ? err.message : "连接测试失败");
    } finally {
      setTestLoading(false);
    }
  }

  function resetImportForm() {
    setImportMode("kubeconfig");
    setConnectionName("");
    setKubeconfigContent("");
    setApiServer("");
    setBearerToken("");
    setCaCert("");
    setSkipTlsVerify(false);
    setTestResult("");
  }

  return (
    <>
      <PageScaffold
        title="集群管理"
        description="保留示例集群切换，同时支持真实集群导入、连接测试、激活与 live 数据预览"
        actions={
          <Stack direction="row" spacing={1}>
            <Button variant="outlined" onClick={() => void Promise.all([loadConnections(), loadLiveData()])}>刷新 live 数据</Button>
            <Button variant="contained" onClick={() => setImportOpen(true)} disabled={!canClusterManage()}>
              导入真实集群
            </Button>
          </Stack>
        }
        toolbar={
          <Stack direction={{ xs: "column", sm: "row" }} spacing={1.5} useFlexGap flexWrap="wrap">
            <TextField size="small" label="按示例集群名称筛选" value={keyword} onChange={(e) => setKeyword(e.target.value)} sx={{ width: 240 }} />
            <TextField size="small" label="按真实连接名称筛选" value={connectionFilter} onChange={(e) => setConnectionFilter(e.target.value)} sx={{ width: 240 }} />
            <Box sx={{ color: "text.secondary", alignSelf: "center" }}>当前角色：{role}</Box>
          </Stack>
        }
      >
        {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}

        <Box sx={{ p: 1.5, borderBottom: "1px solid #d7e1ef", bgcolor: "#f8fbff" }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>示例集群</Typography>
          <Typography variant="body2" color="text.secondary">当前仍保留示例集群切换能力，便于本地演示与回退。</Typography>
        </Box>
        <ResourceTable loading={loading} rows={filteredClusters} rowKey={(r) => r.name} columns={clusterColumns} onRowClick={(row) => setSelected(row)} />

        <Box sx={{ p: 1.5, borderTop: "1px solid #d7e1ef", borderBottom: "1px solid #d7e1ef", bgcolor: "#f8fbff" }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>真实连接管理</Typography>
          <Typography variant="body2" color="text.secondary">支持 kubeconfig 与 API Server + Token + CA 导入，敏感字段默认不回显。</Typography>
        </Box>
        <ResourceTable
          loading={false}
          rows={filteredConnections}
          rowKey={(r) => String(r.id)}
          columns={connectionColumns}
          onRowClick={(row) => {
            if (!canClusterManage()) return;
            void activateConnection(row.id);
          }}
        />

        <Box sx={{ p: 1.5, borderTop: "1px solid #d7e1ef", borderBottom: "1px solid #d7e1ef", bgcolor: "#f8fbff" }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>Live 数据预览</Typography>
          <Typography variant="body2" color="text.secondary">当真实连接激活成功后，这里展示 live cluster 和 live namespaces。</Typography>
        </Box>
        <Box sx={{ p: 2, borderBottom: "1px solid #d7e1ef" }}>
          {liveCluster ? (
            <Stack spacing={1}>
              <Typography variant="body2">名称：{liveCluster.name}</Typography>
              <Typography variant="body2">版本：{liveCluster.version}</Typography>
              <Typography variant="body2">API Server：{liveCluster.apiServer}</Typography>
              <Typography variant="body2">节点数：{liveCluster.nodes}</Typography>
              <Typography variant="body2">数据源：{liveCluster.source}</Typography>
            </Stack>
          ) : (
            <Typography color="text.secondary">尚未激活真实集群连接</Typography>
          )}
        </Box>
        <ResourceTable loading={false} rows={liveNamespaces} rowKey={(r) => r.name} columns={liveNamespaceColumns} />
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
                disabled={!canWorkloadWrite() || current === selected.name || switching === selected.name}
                onClick={() => {
                  if (selected) {
                    void switchCluster(selected.name);
                  }
                }}
              >
                {switching === selected.name ? "切换中..." : "切换到该集群"}
              </Button>
              <Button size="small" onClick={() => setYamlOpen(true)}>查看 YAML</Button>
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

      <YamlDialog open={yamlOpen} title={selected ? `集群 YAML - ${selected.name}` : "集群 YAML"} yaml={yamlText} onClose={() => setYamlOpen(false)} />

      <Dialog open={importOpen} onClose={() => setImportOpen(false)} fullWidth maxWidth="md">
        <DialogTitle>导入真实集群</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 1 }}>
            <Alert severity="info">首批支持两种方式：kubeconfig / API Server + Token + CA。导入与测试仅限 admin。</Alert>
            {testResult && <Alert severity={testResult.startsWith("连接成功") ? "success" : "warning"}>{testResult}</Alert>}
            <TextField label="连接名称" value={connectionName} onChange={(e) => setConnectionName(e.target.value)} fullWidth />
            <Select value={importMode} onChange={(e) => setImportMode(e.target.value as ImportMode)}>
              <MenuItem value="kubeconfig">kubeconfig</MenuItem>
              <MenuItem value="token">API Server + Token + CA</MenuItem>
            </Select>
            {importMode === "kubeconfig" ? (
              <TextField label="kubeconfig 内容" multiline minRows={12} value={kubeconfigContent} onChange={(e) => setKubeconfigContent(e.target.value)} fullWidth />
            ) : (
              <Stack spacing={2}>
                <TextField label="API Server" value={apiServer} onChange={(e) => setApiServer(e.target.value)} fullWidth />
                <TextField label="Bearer Token" value={bearerToken} onChange={(e) => setBearerToken(e.target.value)} fullWidth multiline minRows={3} />
                <TextField label="CA Cert" value={caCert} onChange={(e) => setCaCert(e.target.value)} fullWidth multiline minRows={6} />
                <FormControlLabel control={<Switch checked={skipTlsVerify} onChange={(e) => setSkipTlsVerify(e.target.checked)} />} label="跳过 TLS 校验（仅测试环境）" />
              </Stack>
            )}
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => { resetImportForm(); setImportOpen(false); }}>取消</Button>
          <Button onClick={() => void handleTestConnection()} disabled={testLoading || !canClusterManage()}>
            {testLoading ? "测试中..." : "测试连接"}
          </Button>
          <Button variant="contained" onClick={() => void handleImport()} disabled={submitLoading || !canClusterManage() || !connectionName.trim()}>
            {submitLoading ? "导入中..." : "确认导入"}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
