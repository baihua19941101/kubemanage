import {
  Alert,
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  MenuItem,
  Stack,
  TextField,
  Typography
} from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import TerminalDialog from "../components/framework/TerminalDialog";
import YamlDialog from "../components/framework/YamlDialog";
import { useAuthStore } from "../stores/useAuthStore";
import { useWorkloadStore } from "../stores/useWorkloadStore";

export type WorkloadMode = "deployments" | "pods" | "statefulsets" | "daemonsets" | "jobs" | "cronjobs";

type Props = {
  initialMode?: WorkloadMode;
  showModeSwitcher?: boolean;
};

type DeploymentRow = {
  name: string;
  namespace: string;
  image: string;
  replicas: number;
  ready: number;
  age: string;
};

type PodRow = {
  name: string;
  namespace: string;
  node: string;
  status: string;
  restarts: number;
  ip: string;
  image: string;
  age: string;
};

type StatefulSetRow = {
  name: string;
  namespace: string;
  replicas: number;
  ready: number;
  service: string;
  image: string;
  age: string;
};

type DaemonSetRow = {
  name: string;
  namespace: string;
  desired: number;
  current: number;
  image: string;
  age: string;
};

type JobRow = {
  name: string;
  namespace: string;
  completions: number;
  failed: number;
  status: string;
  age: string;
};

type CronJobRow = {
  name: string;
  namespace: string;
  schedule: string;
  suspend: boolean;
  lastRun: string;
  age: string;
};

type SaveMeta = {
  lastSavedAt?: string;
  lastRequestId?: string;
  history: Array<{ at: string; requestId?: string }>;
};

export default function WorkloadPage({ initialMode = "deployments", showModeSwitcher = true }: Props) {
  const deployments = useWorkloadStore((s) => s.deployments);
  const pods = useWorkloadStore((s) => s.pods);
  const statefulSets = useWorkloadStore((s) => s.statefulSets);
  const daemonSets = useWorkloadStore((s) => s.daemonSets);
  const jobs = useWorkloadStore((s) => s.jobs);
  const cronJobs = useWorkloadStore((s) => s.cronJobs);
  const loading = useWorkloadStore((s) => s.loading);
  const error = useWorkloadStore((s) => s.error);
  const load = useWorkloadStore((s) => s.load);

  const getDeploymentYAML = useWorkloadStore((s) => s.getDeploymentYAML);
  const saveDeploymentYAML = useWorkloadStore((s) => s.saveDeploymentYAML);
  const getPodYAML = useWorkloadStore((s) => s.getPodYAML);
  const savePodYAML = useWorkloadStore((s) => s.savePodYAML);
  const getPodLogs = useWorkloadStore((s) => s.getPodLogs);
  const getTerminalCapabilities = useWorkloadStore((s) => s.getTerminalCapabilities);
  const createTerminalSession = useWorkloadStore((s) => s.createTerminalSession);
  const getStatefulSetYAML = useWorkloadStore((s) => s.getStatefulSetYAML);
  const saveStatefulSetYAML = useWorkloadStore((s) => s.saveStatefulSetYAML);
  const getDaemonSetYAML = useWorkloadStore((s) => s.getDaemonSetYAML);
  const saveDaemonSetYAML = useWorkloadStore((s) => s.saveDaemonSetYAML);
  const getJobYAML = useWorkloadStore((s) => s.getJobYAML);
  const saveJobYAML = useWorkloadStore((s) => s.saveJobYAML);
  const getCronJobYAML = useWorkloadStore((s) => s.getCronJobYAML);
  const saveCronJobYAML = useWorkloadStore((s) => s.saveCronJobYAML);

  const canWorkloadWrite = useAuthStore((s) => s.canWorkloadWrite);

  const [mode, setMode] = useState<WorkloadMode>(initialMode);
  const [keyword, setKeyword] = useState("");
  const [selectedName, setSelectedName] = useState("");
  const [yamlOpen, setYamlOpen] = useState(false);
  const [yamlText, setYamlText] = useState("");
  const [yamlError, setYamlError] = useState("");
  const [yamlNotice, setYamlNotice] = useState("");
  const [yamlLoading, setYamlLoading] = useState(false);
  const [yamlSaveMetaByResource, setYamlSaveMetaByResource] = useState<Record<string, SaveMeta>>({});
  const [logsOpen, setLogsOpen] = useState(false);
  const [rawLogsText, setRawLogsText] = useState("");
  const [logKeyword, setLogKeyword] = useState("");
  const [logContainer, setLogContainer] = useState("");
  const [logContainers, setLogContainers] = useState<string[]>([]);
  const [logCaseSensitive, setLogCaseSensitive] = useState(false);
  const [logMatchOnly, setLogMatchOnly] = useState(false);
  const [logFollow, setLogFollow] = useState(false);
  const [logsError, setLogsError] = useState("");
  const [logsLoading, setLogsLoading] = useState(false);
  const [terminalNotice, setTerminalNotice] = useState("");
  const [terminalOpen, setTerminalOpen] = useState(false);
  const [terminalWsPath, setTerminalWsPath] = useState("");
  const [logsNotice, setLogsNotice] = useState("");

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    setMode(initialMode);
    setSelectedName("");
  }, [initialMode]);

  const selectedResourceKey = selectedName ? `${mode}:${selectedName}` : "";
  const selectedSaveMeta = selectedResourceKey ? yamlSaveMetaByResource[selectedResourceKey] : undefined;

  const lowerKeyword = keyword.toLowerCase().trim();

  const filteredDeployments = useMemo(
    () => deployments.filter((d) => d.name.toLowerCase().includes(lowerKeyword)),
    [deployments, lowerKeyword]
  );
  const filteredPods = useMemo(
    () => pods.filter((p) => p.name.toLowerCase().includes(lowerKeyword)),
    [pods, lowerKeyword]
  );
  const filteredStatefulSets = useMemo(
    () => statefulSets.filter((s) => s.name.toLowerCase().includes(lowerKeyword)),
    [statefulSets, lowerKeyword]
  );
  const filteredDaemonSets = useMemo(
    () => daemonSets.filter((d) => d.name.toLowerCase().includes(lowerKeyword)),
    [daemonSets, lowerKeyword]
  );
  const filteredJobs = useMemo(
    () => jobs.filter((j) => j.name.toLowerCase().includes(lowerKeyword)),
    [jobs, lowerKeyword]
  );
  const filteredCronJobs = useMemo(
    () => cronJobs.filter((c) => c.name.toLowerCase().includes(lowerKeyword)),
    [cronJobs, lowerKeyword]
  );

  const selectedDeployment = deployments.find((x) => x.name === selectedName) ?? null;
  const selectedPod = pods.find((x) => x.name === selectedName) ?? null;
  const selectedStatefulSet = statefulSets.find((x) => x.name === selectedName) ?? null;
  const selectedDaemonSet = daemonSets.find((x) => x.name === selectedName) ?? null;
  const selectedJob = jobs.find((x) => x.name === selectedName) ?? null;
  const selectedCronJob = cronJobs.find((x) => x.name === selectedName) ?? null;

  const visibleLogLines = useMemo(
    () => rawLogsText.split("\n").filter((line) => line.length > 0),
    [rawLogsText]
  );

  const matchedLogCount = useMemo(() => {
    if (!logKeyword.trim()) {
      return visibleLogLines.length;
    }
    const needle = logCaseSensitive ? logKeyword : logKeyword.toLowerCase();
    return visibleLogLines.filter((line) => {
      const haystack = logCaseSensitive ? line : line.toLowerCase();
      return haystack.includes(needle);
    }).length;
  }, [logCaseSensitive, logKeyword, visibleLogLines]);

  const deploymentColumns = [
    { key: "name", header: "名称", render: (r: DeploymentRow) => r.name },
    { key: "ns", header: "命名空间", render: (r: DeploymentRow) => r.namespace },
    { key: "image", header: "镜像", render: (r: DeploymentRow) => r.image },
    { key: "replicas", header: "副本", render: (r: DeploymentRow) => `${r.ready}/${r.replicas}` },
    { key: "age", header: "Age", render: (r: DeploymentRow) => r.age }
  ];

  const podColumns = [
    { key: "name", header: "名称", render: (r: PodRow) => r.name },
    { key: "ns", header: "命名空间", render: (r: PodRow) => r.namespace },
    { key: "status", header: "状态", render: (r: PodRow) => r.status },
    { key: "node", header: "节点", render: (r: PodRow) => r.node },
    { key: "age", header: "Age", render: (r: PodRow) => r.age }
  ];

  const statefulColumns = [
    { key: "name", header: "名称", render: (r: StatefulSetRow) => r.name },
    { key: "ns", header: "命名空间", render: (r: StatefulSetRow) => r.namespace },
    { key: "service", header: "服务名", render: (r: StatefulSetRow) => r.service },
    { key: "replicas", header: "副本", render: (r: StatefulSetRow) => `${r.ready}/${r.replicas}` },
    { key: "age", header: "Age", render: (r: StatefulSetRow) => r.age }
  ];

  const daemonColumns = [
    { key: "name", header: "名称", render: (r: DaemonSetRow) => r.name },
    { key: "ns", header: "命名空间", render: (r: DaemonSetRow) => r.namespace },
    { key: "image", header: "镜像", render: (r: DaemonSetRow) => r.image },
    { key: "desired", header: "期望/当前", render: (r: DaemonSetRow) => `${r.current}/${r.desired}` },
    { key: "age", header: "Age", render: (r: DaemonSetRow) => r.age }
  ];

  const jobColumns = [
    { key: "name", header: "名称", render: (r: JobRow) => r.name },
    { key: "ns", header: "命名空间", render: (r: JobRow) => r.namespace },
    { key: "status", header: "状态", render: (r: JobRow) => r.status },
    { key: "comp", header: "完成/失败", render: (r: JobRow) => `${r.completions}/${r.failed}` },
    { key: "age", header: "Age", render: (r: JobRow) => r.age }
  ];

  const cronColumns = [
    { key: "name", header: "名称", render: (r: CronJobRow) => r.name },
    { key: "ns", header: "命名空间", render: (r: CronJobRow) => r.namespace },
    { key: "schedule", header: "调度", render: (r: CronJobRow) => r.schedule },
    { key: "suspend", header: "暂停", render: (r: CronJobRow) => (r.suspend ? "是" : "否") },
    { key: "age", header: "Age", render: (r: CronJobRow) => r.age }
  ];

  const drawerOpen = selectedName.length > 0;
  const currentLabel =
    mode === "deployments"
      ? "Deployment"
      : mode === "pods"
        ? "Pod"
        : mode === "statefulsets"
          ? "StatefulSet"
          : mode === "daemonsets"
            ? "DaemonSet"
            : mode === "jobs"
              ? "Job"
              : "CronJob";

  async function openYaml() {
    if (!selectedName) return;
    setYamlLoading(true);
    setYamlError("");
    setYamlNotice("");
    try {
      if (mode === "deployments") {
        setYamlText(await getDeploymentYAML(selectedName));
      } else if (mode === "pods") {
        setYamlText(await getPodYAML(selectedName));
      } else if (mode === "statefulsets") {
        setYamlText(await getStatefulSetYAML(selectedName));
      } else if (mode === "daemonsets") {
        setYamlText(await getDaemonSetYAML(selectedName));
      } else if (mode === "jobs") {
        setYamlText(await getJobYAML(selectedName));
      } else {
        setYamlText(await getCronJobYAML(selectedName));
      }
      setYamlOpen(true);
    } catch (err) {
      setYamlError(err instanceof Error ? err.message : "打开 YAML 失败");
    } finally {
      setYamlLoading(false);
    }
  }

  async function saveYaml(yaml: string) {
    if (!selectedName) return;
    const saveKey = `${mode}:${selectedName}`;
    setYamlLoading(true);
    setYamlError("");
    setYamlNotice("");
    let result: { ok: true; requestId?: string } | null = null;
    try {
      if (mode === "deployments") {
        result = await saveDeploymentYAML(selectedName, yaml);
      } else if (mode === "pods") {
        result = await savePodYAML(selectedName, yaml);
      } else if (mode === "statefulsets") {
        result = await saveStatefulSetYAML(selectedName, yaml);
      } else if (mode === "daemonsets") {
        result = await saveDaemonSetYAML(selectedName, yaml);
      } else if (mode === "jobs") {
        result = await saveJobYAML(selectedName, yaml);
      } else {
        result = await saveCronJobYAML(selectedName, yaml);
      }
      if (result && result.ok) {
        await load();
        setYamlText(yaml);
        const savedAt = new Date().toISOString();
        const savedRequestId = result.requestId;
        setYamlSaveMetaByResource((prev) => {
          const current = prev[saveKey] || { history: [] };
          return {
            ...prev,
            [saveKey]: {
              lastSavedAt: savedAt,
              lastRequestId: savedRequestId || undefined,
              history: [{ at: savedAt, requestId: savedRequestId }, ...current.history].slice(0, 10)
            }
          };
        });
        const requestIdText = savedRequestId ? `（requestId: ${savedRequestId}）` : "";
        setYamlNotice(`${currentLabel} YAML 保存成功${requestIdText}`);
      } else {
        setYamlError("保存 YAML 失败，请检查权限与请求参数");
      }
    } catch (err) {
      setYamlError(err instanceof Error ? err.message : "保存 YAML 失败");
    } finally {
      setYamlLoading(false);
    }
  }

  async function openLogs() {
    if (!selectedName) return;
    setTerminalNotice("");
    setLogsNotice("");
    setLogKeyword("");
    setLogContainer("");
    setLogContainers([]);
    setLogCaseSensitive(false);
    setLogMatchOnly(false);
    setLogFollow(false);
    setRawLogsText("");
    setLogsError("");
    setLogsOpen(true);
    setLogsLoading(true);
    try {
      const capabilities = await getTerminalCapabilities(selectedName);
      setTerminalNotice(capabilities.message);
      setLogContainers(capabilities.containers || []);
      if ((capabilities.containers || []).length > 0) {
        setLogContainer(capabilities.containers![0]);
      }
    } catch (err) {
      setTerminalNotice(err instanceof Error ? err.message : "获取终端能力失败");
    } finally {
      setLogsLoading(false);
    }
  }

  async function refreshLogs() {
    if (!selectedName) return;
    setLogsLoading(true);
    setLogsError("");
    try {
      const logs = await getPodLogs(selectedName, {
        container: logContainer || undefined,
        keyword: logKeyword,
        caseSensitive: logCaseSensitive,
        matchOnly: logMatchOnly,
        follow: logFollow
      });
      setRawLogsText(logs);
    } catch (err) {
      setLogsError(err instanceof Error ? err.message : "获取 Pod 日志失败");
    } finally {
      setLogsLoading(false);
    }
  }

  async function openTerminal() {
    if (!selectedName) return;
    try {
      const result = await createTerminalSession(selectedName, logContainer || undefined);
      if (result.wsPath) {
        setTerminalWsPath(result.wsPath);
        setTerminalOpen(true);
        const ttlHint =
          result.ttlSeconds && result.expiresAt
            ? `（TTL ${result.ttlSeconds}s，过期时间 ${result.expiresAt}）`
            : result.ttlSeconds
            ? `（TTL ${result.ttlSeconds}s）`
            : "";
        setTerminalNotice(`终端会话已创建${ttlHint}，可使用 WebSocket 地址：${result.wsPath}`);
      } else {
        setTerminalNotice(result.error || "terminal gateway not enabled");
      }
    } catch (err) {
      setTerminalNotice(err instanceof Error ? err.message : "终端能力暂不可用");
    }
  }

  async function copyLogs() {
    if (!rawLogsText) return;
    try {
      await navigator.clipboard.writeText(rawLogsText);
      setLogsNotice("日志已复制到剪贴板");
    } catch {
      setLogsNotice("复制失败，请检查浏览器剪贴板权限");
    }
  }

  function clearLogFilters() {
    setLogKeyword("");
    setLogCaseSensitive(false);
    setLogMatchOnly(false);
    setLogsNotice("已清空日志筛选条件");
  }

  function downloadLogs() {
    if (!selectedName) return;
    const content = rawLogsText;
    const blob = new Blob([content], { type: "text/plain;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `${selectedName}-${formatTimestamp(new Date())}.log`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }

  useEffect(() => {
    if (!logsOpen || !selectedName) return;
    void refreshLogs();
  }, [logsOpen, selectedName, logContainer, logKeyword, logCaseSensitive, logMatchOnly]);

  useEffect(() => {
    if (!logsOpen || !selectedName || !logFollow) return;
    const timer = window.setInterval(() => {
      void refreshLogs();
    }, 2000);
    return () => window.clearInterval(timer);
  }, [logsOpen, selectedName, logFollow, logContainer, logKeyword, logCaseSensitive, logMatchOnly]);

  return (
    <>
      <PageScaffold
        title="工作负载管理"
        description="统一管理 Deployment/Pod/StatefulSet/DaemonSet/Job/CronJob，支持 YAML 编辑与日志查看"
        actions={showModeSwitcher ? (
          <Stack direction="row" spacing={1} useFlexGap flexWrap="wrap">
            <Button variant={mode === "deployments" ? "contained" : "outlined"} onClick={() => setMode("deployments")}>Deployment</Button>
            <Button variant={mode === "pods" ? "contained" : "outlined"} onClick={() => setMode("pods")}>Pod</Button>
            <Button variant={mode === "statefulsets" ? "contained" : "outlined"} onClick={() => setMode("statefulsets")}>StatefulSet</Button>
            <Button variant={mode === "daemonsets" ? "contained" : "outlined"} onClick={() => setMode("daemonsets")}>DaemonSet</Button>
            <Button variant={mode === "jobs" ? "contained" : "outlined"} onClick={() => setMode("jobs")}>Job</Button>
            <Button variant={mode === "cronjobs" ? "contained" : "outlined"} onClick={() => setMode("cronjobs")}>CronJob</Button>
          </Stack>
        ) : null}
        toolbar={
          <TextField
            size="small"
            label="按名称筛选"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            sx={{ width: 280 }}
          />
        }
      >
        {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}
        {yamlError && <Alert severity="error" sx={{ m: 1.5 }}>{yamlError}</Alert>}
        {yamlNotice && <Alert severity="success" sx={{ m: 1.5 }}>{yamlNotice}</Alert>}

        {mode === "deployments" && (
          <ResourceTable
            loading={loading}
            rows={filteredDeployments}
            rowKey={(r) => r.name}
            columns={deploymentColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "pods" && (
          <ResourceTable
            loading={loading}
            rows={filteredPods}
            rowKey={(r) => r.name}
            columns={podColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "statefulsets" && (
          <ResourceTable
            loading={loading}
            rows={filteredStatefulSets}
            rowKey={(r) => r.name}
            columns={statefulColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "daemonsets" && (
          <ResourceTable
            loading={loading}
            rows={filteredDaemonSets}
            rowKey={(r) => r.name}
            columns={daemonColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "jobs" && (
          <ResourceTable
            loading={loading}
            rows={filteredJobs}
            rowKey={(r) => r.name}
            columns={jobColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}

        {mode === "cronjobs" && (
          <ResourceTable
            loading={loading}
            rows={filteredCronJobs}
            rowKey={(r) => r.name}
            columns={cronColumns}
            onRowClick={(r) => setSelectedName(r.name)}
          />
        )}
      </PageScaffold>

      <DetailDrawer
        open={drawerOpen}
        title={selectedName ? `${currentLabel} 详情 - ${selectedName}` : `${currentLabel} 详情`}
        onClose={() => setSelectedName("")}
        actions={
          selectedName ? (
            <Stack direction="row" spacing={1}>
              <Button size="small" onClick={openYaml} disabled={yamlLoading}>
                {yamlLoading ? "YAML 加载中..." : "查看/编辑 YAML"}
              </Button>
              {mode === "pods" && (
                <Button size="small" onClick={() => void openLogs()} disabled={logsLoading}>
                  查看日志
                </Button>
              )}
            </Stack>
          ) : null
        }
      >
        {mode === "deployments" && selectedDeployment && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedDeployment.name}</Typography>
            <Typography variant="body2">命名空间：{selectedDeployment.namespace}</Typography>
            <Typography variant="body2">镜像：{selectedDeployment.image}</Typography>
            <Typography variant="body2">副本：{selectedDeployment.ready}/{selectedDeployment.replicas}</Typography>
          </Stack>
        )}

        {mode === "pods" && selectedPod && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedPod.name}</Typography>
            <Typography variant="body2">命名空间：{selectedPod.namespace}</Typography>
            <Typography variant="body2">状态：{selectedPod.status}</Typography>
            <Typography variant="body2">节点：{selectedPod.node}</Typography>
            <Typography variant="body2">IP：{selectedPod.ip}</Typography>
          </Stack>
        )}

        {mode === "statefulsets" && selectedStatefulSet && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedStatefulSet.name}</Typography>
            <Typography variant="body2">命名空间：{selectedStatefulSet.namespace}</Typography>
            <Typography variant="body2">服务名：{selectedStatefulSet.service}</Typography>
            <Typography variant="body2">镜像：{selectedStatefulSet.image}</Typography>
            <Typography variant="body2">副本：{selectedStatefulSet.ready}/{selectedStatefulSet.replicas}</Typography>
          </Stack>
        )}

        {mode === "daemonsets" && selectedDaemonSet && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedDaemonSet.name}</Typography>
            <Typography variant="body2">命名空间：{selectedDaemonSet.namespace}</Typography>
            <Typography variant="body2">镜像：{selectedDaemonSet.image}</Typography>
            <Typography variant="body2">调度：{selectedDaemonSet.current}/{selectedDaemonSet.desired}</Typography>
          </Stack>
        )}

        {mode === "jobs" && selectedJob && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedJob.name}</Typography>
            <Typography variant="body2">命名空间：{selectedJob.namespace}</Typography>
            <Typography variant="body2">状态：{selectedJob.status}</Typography>
            <Typography variant="body2">完成/失败：{selectedJob.completions}/{selectedJob.failed}</Typography>
          </Stack>
        )}

        {mode === "cronjobs" && selectedCronJob && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedCronJob.name}</Typography>
            <Typography variant="body2">命名空间：{selectedCronJob.namespace}</Typography>
            <Typography variant="body2">调度：{selectedCronJob.schedule}</Typography>
            <Typography variant="body2">暂停：{selectedCronJob.suspend ? "是" : "否"}</Typography>
            <Typography variant="body2">最近执行：{selectedCronJob.lastRun}</Typography>
          </Stack>
        )}
      </DetailDrawer>

      <YamlDialog
        open={yamlOpen}
        title={`${currentLabel} YAML 编辑`}
        yaml={yamlText}
        onClose={() => setYamlOpen(false)}
        onSave={!canWorkloadWrite() ? undefined : saveYaml}
        saving={yamlLoading}
        saveMeta={{
          lastSavedAt: selectedSaveMeta?.lastSavedAt,
          lastRequestId: selectedSaveMeta?.lastRequestId,
          history: selectedSaveMeta?.history || []
        }}
      />

      <Dialog open={logsOpen} onClose={() => setLogsOpen(false)} fullWidth maxWidth="md">
        <DialogTitle>{selectedName ? `Pod 日志 - ${selectedName}` : "Pod 日志"}</DialogTitle>
        <DialogContent>
          <Stack spacing={2} sx={{ mt: 1 }}>
            {logsError && <Alert severity="error">{logsError}</Alert>}
            {logsNotice && <Alert severity="success">{logsNotice}</Alert>}
            {terminalNotice && <Alert severity="info">终端预留状态：{terminalNotice}</Alert>}
            <Stack direction={{ xs: "column", md: "row" }} spacing={1.5}>
              {logContainers.length > 0 && (
                <TextField
                  select
                  size="small"
                  label="容器"
                  value={logContainer}
                  onChange={(e) => setLogContainer(e.target.value)}
                  sx={{ minWidth: 180 }}
                >
                  {logContainers.map((name) => (
                    <MenuItem key={name} value={name}>
                      {name}
                    </MenuItem>
                  ))}
                </TextField>
              )}
              <TextField
                size="small"
                label="日志关键字"
                value={logKeyword}
                onChange={(e) => setLogKeyword(e.target.value)}
                sx={{ minWidth: 240 }}
              />
              <FormControlLabel
                control={<Checkbox checked={logCaseSensitive} onChange={(e) => setLogCaseSensitive(e.target.checked)} />}
                label="大小写敏感"
              />
              <FormControlLabel
                control={<Checkbox checked={logMatchOnly} onChange={(e) => setLogMatchOnly(e.target.checked)} />}
                label="仅显示匹配行"
              />
              <FormControlLabel
                control={<Checkbox checked={logFollow} onChange={(e) => setLogFollow(e.target.checked)} />}
                label="跟随刷新"
              />
            </Stack>
            <Stack direction={{ xs: "column", md: "row" }} spacing={1.5} alignItems={{ md: "center" }}>
              <Typography variant="body2" color="text.secondary">
                当前显示 {visibleLogLines.length} 行
              </Typography>
              <Typography variant="body2" color="text.secondary">
                匹配 {matchedLogCount} 行
              </Typography>
              <Button size="small" onClick={clearLogFilters}>清空筛选</Button>
              <Button size="small" onClick={() => void refreshLogs()} disabled={logsLoading}>
                {logsLoading ? "刷新中..." : "立即刷新"}
              </Button>
              <Button size="small" onClick={() => void copyLogs()} disabled={!rawLogsText}>复制日志</Button>
            </Stack>
            <TextField
              multiline
              minRows={16}
              fullWidth
              value={rawLogsText}
              InputProps={{ readOnly: true }}
              placeholder={logsLoading ? "日志加载中..." : "暂无日志"}
            />
          </Stack>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => void openTerminal()}>打开终端</Button>
          <Button onClick={downloadLogs} disabled={!rawLogsText}>导出日志</Button>
          <Button onClick={() => setLogsOpen(false)}>关闭</Button>
        </DialogActions>
      </Dialog>

      <TerminalDialog
        open={terminalOpen}
        title={selectedName ? `Pod 终端 - ${selectedName}` : "Pod 终端"}
        wsPath={terminalWsPath}
        onClose={() => setTerminalOpen(false)}
      />
    </>
  );
}

function formatTimestamp(date: Date) {
  const yyyy = String(date.getFullYear());
  const mm = String(date.getMonth() + 1).padStart(2, "0");
  const dd = String(date.getDate()).padStart(2, "0");
  const hh = String(date.getHours()).padStart(2, "0");
  const mi = String(date.getMinutes()).padStart(2, "0");
  const ss = String(date.getSeconds()).padStart(2, "0");
  return `${yyyy}${mm}${dd}-${hh}${mi}${ss}`;
}
