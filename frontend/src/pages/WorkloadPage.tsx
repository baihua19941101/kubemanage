import { Alert, Box, Button, Stack, TextField, Typography } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import YamlDialog from "../components/framework/YamlDialog";
import { useAuthStore } from "../stores/useAuthStore";
import { useWorkloadStore } from "../stores/useWorkloadStore";

type Mode = "deployments" | "pods";

type DeployRow = {
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

export default function WorkloadPage() {
  const deployments = useWorkloadStore((s) => s.deployments);
  const pods = useWorkloadStore((s) => s.pods);
  const loading = useWorkloadStore((s) => s.loading);
  const error = useWorkloadStore((s) => s.error);
  const load = useWorkloadStore((s) => s.load);
  const getDeploymentYAML = useWorkloadStore((s) => s.getDeploymentYAML);
  const saveDeploymentYAML = useWorkloadStore((s) => s.saveDeploymentYAML);
  const getPodYAML = useWorkloadStore((s) => s.getPodYAML);
  const savePodYAML = useWorkloadStore((s) => s.savePodYAML);
  const getPodLogs = useWorkloadStore((s) => s.getPodLogs);
  const canWorkloadWrite = useAuthStore((s) => s.canWorkloadWrite);

  const [mode, setMode] = useState<Mode>("deployments");
  const [keyword, setKeyword] = useState("");
  const [selectedDeploy, setSelectedDeploy] = useState<DeployRow | null>(null);
  const [selectedPod, setSelectedPod] = useState<PodRow | null>(null);
  const [yamlOpen, setYamlOpen] = useState(false);
  const [yamlText, setYamlText] = useState("");
  const [logsOpen, setLogsOpen] = useState(false);
  const [logsText, setLogsText] = useState("");

  useEffect(() => {
    void load();
  }, [load]);

  const filteredDeployments = useMemo(
    () =>
      deployments.filter((d) =>
        d.name.toLowerCase().includes(keyword.toLowerCase().trim())
      ),
    [deployments, keyword]
  );
  const filteredPods = useMemo(
    () =>
      pods.filter((p) => p.name.toLowerCase().includes(keyword.toLowerCase().trim())),
    [pods, keyword]
  );

  const deploymentColumns = [
    { key: "name", header: "名称", render: (r: DeployRow) => r.name },
    { key: "ns", header: "命名空间", render: (r: DeployRow) => r.namespace },
    { key: "image", header: "镜像", render: (r: DeployRow) => r.image },
    { key: "replicas", header: "副本", render: (r: DeployRow) => `${r.ready}/${r.replicas}` },
    { key: "age", header: "Age", render: (r: DeployRow) => r.age }
  ];
  const podColumns = [
    { key: "name", header: "名称", render: (r: PodRow) => r.name },
    { key: "ns", header: "命名空间", render: (r: PodRow) => r.namespace },
    { key: "status", header: "状态", render: (r: PodRow) => r.status },
    { key: "node", header: "节点", render: (r: PodRow) => r.node },
    { key: "age", header: "Age", render: (r: PodRow) => r.age }
  ];

  return (
    <>
      <PageScaffold
        title="工作负载管理"
        description="统一管理 Deployment/Pod，支持 YAML 编辑与日志查看"
        actions={
          <Stack direction="row" spacing={1}>
            <Button
              variant={mode === "deployments" ? "contained" : "outlined"}
              onClick={() => setMode("deployments")}
            >
              Deployment
            </Button>
            <Button
              variant={mode === "pods" ? "contained" : "outlined"}
              onClick={() => setMode("pods")}
            >
              Pod
            </Button>
          </Stack>
        }
        toolbar={
          <TextField
            size="small"
            label="按名称筛选"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            sx={{ width: 260 }}
          />
        }
      >
        {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}
        {mode === "deployments" ? (
          <ResourceTable
            loading={loading}
            rows={filteredDeployments}
            rowKey={(r) => r.name}
            columns={deploymentColumns}
            onRowClick={(r) => setSelectedDeploy(r)}
          />
        ) : (
          <ResourceTable
            loading={loading}
            rows={filteredPods}
            rowKey={(r) => r.name}
            columns={podColumns}
            onRowClick={(r) => setSelectedPod(r)}
          />
        )}
      </PageScaffold>

      <DetailDrawer
        open={Boolean(selectedDeploy)}
        title={selectedDeploy ? `Deployment 详情 - ${selectedDeploy.name}` : "Deployment 详情"}
        onClose={() => setSelectedDeploy(null)}
        actions={
          selectedDeploy ? (
            <Button
              size="small"
              onClick={async () => {
                if (!selectedDeploy) return;
                const y = await getDeploymentYAML(selectedDeploy.name);
                setYamlText(y);
                setYamlOpen(true);
              }}
            >
              查看/编辑 YAML
            </Button>
          ) : null
        }
      >
        {selectedDeploy && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedDeploy.name}</Typography>
            <Typography variant="body2">命名空间：{selectedDeploy.namespace}</Typography>
            <Typography variant="body2">镜像：{selectedDeploy.image}</Typography>
            <Typography variant="body2">
              副本：{selectedDeploy.ready}/{selectedDeploy.replicas}
            </Typography>
          </Stack>
        )}
      </DetailDrawer>

      <DetailDrawer
        open={Boolean(selectedPod)}
        title={selectedPod ? `Pod 详情 - ${selectedPod.name}` : "Pod 详情"}
        onClose={() => setSelectedPod(null)}
        actions={
          selectedPod ? (
            <Stack direction="row" spacing={1}>
              <Button
                size="small"
                onClick={async () => {
                  if (!selectedPod) return;
                  const y = await getPodYAML(selectedPod.name);
                  setYamlText(y);
                  setYamlOpen(true);
                }}
              >
                查看/编辑 YAML
              </Button>
              <Button
                size="small"
                onClick={async () => {
                  if (!selectedPod) return;
                  const logs = await getPodLogs(selectedPod.name);
                  setLogsText(logs);
                  setLogsOpen(true);
                }}
              >
                查看日志
              </Button>
            </Stack>
          ) : null
        }
      >
        {selectedPod && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedPod.name}</Typography>
            <Typography variant="body2">命名空间：{selectedPod.namespace}</Typography>
            <Typography variant="body2">状态：{selectedPod.status}</Typography>
            <Typography variant="body2">节点：{selectedPod.node}</Typography>
            <Typography variant="body2">IP：{selectedPod.ip}</Typography>
          </Stack>
        )}
      </DetailDrawer>

      <YamlDialog
        open={yamlOpen}
        title="YAML 编辑"
        yaml={yamlText}
        onClose={() => setYamlOpen(false)}
        onSave={
          !canWorkloadWrite()
            ? undefined
            : async (yaml) => {
                if (selectedDeploy) {
                  const ok = await saveDeploymentYAML(selectedDeploy.name, yaml);
                  if (ok) {
                    setYamlOpen(false);
                  }
                  return;
                }
                if (selectedPod) {
                  const ok = await savePodYAML(selectedPod.name, yaml);
                  if (ok) {
                    setYamlOpen(false);
                  }
                }
              }
        }
      />

      <YamlDialog
        open={logsOpen}
        title={selectedPod ? `Pod 日志 - ${selectedPod.name}` : "Pod 日志"}
        yaml={logsText}
        onClose={() => setLogsOpen(false)}
      />
    </>
  );
}
