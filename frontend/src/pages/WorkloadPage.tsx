import { Alert, Button, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
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
  const [logsOpen, setLogsOpen] = useState(false);
  const [logsText, setLogsText] = useState("");

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    setMode(initialMode);
    setSelectedName("");
  }, [initialMode]);

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
  }

  async function saveYaml(yaml: string) {
    if (!selectedName) return;
    let ok = false;
    if (mode === "deployments") {
      ok = await saveDeploymentYAML(selectedName, yaml);
    } else if (mode === "pods") {
      ok = await savePodYAML(selectedName, yaml);
    } else if (mode === "statefulsets") {
      ok = await saveStatefulSetYAML(selectedName, yaml);
    } else if (mode === "daemonsets") {
      ok = await saveDaemonSetYAML(selectedName, yaml);
    } else if (mode === "jobs") {
      ok = await saveJobYAML(selectedName, yaml);
    } else {
      ok = await saveCronJobYAML(selectedName, yaml);
    }
    if (ok) {
      setYamlOpen(false);
    }
  }

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
              <Button size="small" onClick={openYaml}>查看/编辑 YAML</Button>
              {mode === "pods" && (
                <Button
                  size="small"
                  onClick={async () => {
                    if (!selectedName) return;
                    const logs = await getPodLogs(selectedName);
                    setLogsText(logs);
                    setLogsOpen(true);
                  }}
                >
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
      />

      <YamlDialog
        open={logsOpen}
        title={selectedName ? `Pod 日志 - ${selectedName}` : "Pod 日志"}
        yaml={logsText}
        onClose={() => setLogsOpen(false)}
      />
    </>
  );
}
