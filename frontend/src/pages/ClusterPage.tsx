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
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import { useAuthStore } from "../stores/useAuthStore";
import { useClusterStore } from "../stores/useClusterStore";

type ClusterRow = {
  state: string;
  name: string;
  provider: string;
  distro: string;
  kubernetesVersion: string;
  architecture: string;
  cpu: string;
  memory: string;
  pods: number;
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

type ImportMode = "kubeconfig" | "token";

export default function ClusterPage() {
  const clusters = useClusterStore((s) => s.clusters);
  const loading = useClusterStore((s) => s.loading);
  const connections = useClusterStore((s) => s.connections);
  const error = useClusterStore((s) => s.error);
  const load = useClusterStore((s) => s.load);
  const loadConnections = useClusterStore((s) => s.loadConnections);
  const importKubeconfig = useClusterStore((s) => s.importKubeconfig);
  const importToken = useClusterStore((s) => s.importToken);
  const testConnection = useClusterStore((s) => s.testConnection);
  const activateConnection = useClusterStore((s) => s.activateConnection);
  const canClusterManage = useAuthStore((s) => s.canClusterManage);
  const role = useAuthStore((s) => s.role);

  const [clusterFilter, setClusterFilter] = useState("");
  const [connectionFilter, setConnectionFilter] = useState("");
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
  const [testResult, setTestResult] = useState("");

  useEffect(() => {
    void Promise.all([load(), loadConnections()]);
  }, [load, loadConnections]);

  const filteredClusters = useMemo(
    () =>
      clusters.filter((c) => {
        const keyword = clusterFilter.toLowerCase().trim();
        return c.name.toLowerCase().includes(keyword) || c.provider.toLowerCase().includes(keyword) || c.distro.toLowerCase().includes(keyword);
      }),
    [clusters, clusterFilter]
  );

  const filteredConnections = useMemo(
    () => connections.filter((c) => c.name.toLowerCase().includes(connectionFilter.toLowerCase().trim())),
    [connections, connectionFilter]
  );

  const clusterColumns = [
    {
      key: "state",
      header: "State",
      render: (row: ClusterRow) => <Chip size="small" color={row.state === "Ready" ? "success" : "warning"} label={row.state} />
    },
    { key: "name", header: "Name", render: (row: ClusterRow) => row.name },
    {
      key: "providerDistro",
      header: "Provider / Distro",
      render: (row: ClusterRow) => `${row.provider || "-"} / ${row.distro || "-"}`
    },
    {
      key: "versionArch",
      header: "Kubernetes Version / Architecture",
      render: (row: ClusterRow) => `${row.kubernetesVersion || "-"} / ${row.architecture || "-"}`
    },
    { key: "cpu", header: "CPU", render: (row: ClusterRow) => row.cpu || "-" },
    { key: "memory", header: "Memory", render: (row: ClusterRow) => row.memory || "-" },
    { key: "pods", header: "Pods", render: (row: ClusterRow) => row.pods }
  ];

  const connectionColumns = [
    { key: "name", header: "连接名称", render: (row: ConnectionRow) => row.name },
    { key: "mode", header: "模式", render: (row: ConnectionRow) => row.mode },
    { key: "server", header: "API Server", render: (row: ConnectionRow) => row.apiServer || "-" },
    {
      key: "status",
      header: "状态",
      render: (row: ConnectionRow) => (
        <Chip size="small" color={row.status === "connected" ? "success" : row.status === "failed" ? "error" : "default"} label={row.status} />
      )
    },
    {
      key: "active",
      header: "已激活",
      render: (row: ConnectionRow) => (row.isDefault ? <Chip size="small" color="primary" label="active" /> : "-")
    }
  ];

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
        await Promise.all([loadConnections(), load()]);
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
        description="仅保留真实连接与真实数据展示；不再提供 mock 集群数据和 Live 概览区"
        actions={
          <Stack direction="row" spacing={1}>
            <Button variant="outlined" onClick={() => void Promise.all([load(), loadConnections()])}>
              刷新
            </Button>
            <Button variant="contained" onClick={() => setImportOpen(true)} disabled={!canClusterManage()}>
              导入真实集群
            </Button>
          </Stack>
        }
        toolbar={
          <Stack direction={{ xs: "column", sm: "row" }} spacing={1.5} useFlexGap flexWrap="wrap">
            <TextField size="small" label="筛选集群（名称/Provider/Distro）" value={clusterFilter} onChange={(e) => setClusterFilter(e.target.value)} sx={{ width: 280 }} />
            <TextField size="small" label="筛选连接名称" value={connectionFilter} onChange={(e) => setConnectionFilter(e.target.value)} sx={{ width: 220 }} />
            <Box sx={{ color: "text.secondary", alignSelf: "center" }}>当前角色：{role}</Box>
          </Stack>
        }
      >
        {error && (
          <Alert severity="error" sx={{ m: 1.5 }}>
            {error}
          </Alert>
        )}

        <Box sx={{ p: 1.5, borderBottom: "1px solid #d7e1ef", bgcolor: "#f8fbff" }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
            真实集群列表
          </Typography>
          <Typography variant="body2" color="text.secondary">
            展示真实集群字段：State / Name / Provider / Distro / Kubernetes Version / Architecture / CPU / Memory / Pods（组合列展示）
          </Typography>
        </Box>
        <ResourceTable loading={loading} rows={filteredClusters} rowKey={(r) => r.name} columns={clusterColumns} />

        <Box sx={{ p: 1.5, borderTop: "1px solid #d7e1ef", borderBottom: "1px solid #d7e1ef", bgcolor: "#f8fbff" }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>
            真实连接管理
          </Typography>
          <Typography variant="body2" color="text.secondary">
            点击连接行可执行激活，激活后主资源页读取该真实集群数据。
          </Typography>
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
      </PageScaffold>

      <Dialog
        open={importOpen}
        onClose={() => {
          setImportOpen(false);
          resetImportForm();
        }}
        fullWidth
        maxWidth="md"
      >
        <DialogTitle>导入真实集群连接</DialogTitle>
        <DialogContent dividers>
          <Stack spacing={2}>
            <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
              <TextField label="连接名称" value={connectionName} onChange={(e) => setConnectionName(e.target.value)} fullWidth />
              <Box sx={{ minWidth: 180 }}>
                <Typography variant="caption" color="text.secondary">
                  导入模式
                </Typography>
                <Select value={importMode} onChange={(e) => setImportMode(e.target.value as ImportMode)} fullWidth size="small">
                  <MenuItem value="kubeconfig">kubeconfig</MenuItem>
                  <MenuItem value="token">API Server + Token + CA</MenuItem>
                </Select>
              </Box>
            </Stack>

            {importMode === "kubeconfig" ? (
              <TextField
                label="kubeconfig 内容"
                placeholder="粘贴 kubeconfig YAML"
                value={kubeconfigContent}
                onChange={(e) => setKubeconfigContent(e.target.value)}
                multiline
                minRows={10}
                fullWidth
              />
            ) : (
              <Stack spacing={2}>
                <TextField label="API Server" placeholder="https://127.0.0.1:6443" value={apiServer} onChange={(e) => setApiServer(e.target.value)} fullWidth />
                <TextField label="Bearer Token" value={bearerToken} onChange={(e) => setBearerToken(e.target.value)} fullWidth multiline minRows={3} />
                <TextField label="CA 证书（可选）" value={caCert} onChange={(e) => setCaCert(e.target.value)} fullWidth multiline minRows={4} />
                <FormControlLabel control={<Switch checked={skipTlsVerify} onChange={(e) => setSkipTlsVerify(e.target.checked)} />} label="跳过 TLS 校验（仅测试环境）" />
              </Stack>
            )}

            {testResult && <Alert severity={testResult.startsWith("连接成功") ? "success" : "error"}>{testResult}</Alert>}
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              setImportOpen(false);
              resetImportForm();
            }}
          >
            取消
          </Button>
          <Button onClick={() => void handleTestConnection()} disabled={testLoading || submitLoading || !canClusterManage()}>
            {testLoading ? "测试中..." : "测试连接"}
          </Button>
          <Button variant="contained" onClick={() => void handleImport()} disabled={submitLoading || !canClusterManage()}>
            {submitLoading ? "导入中..." : "导入并保存"}
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
