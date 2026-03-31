import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  TextField,
  Typography
} from "@mui/material";
import { useEffect, useState } from "react";
import { useWorkloadStore } from "../stores/useWorkloadStore";

type EditorMode = "deployment" | "pod" | "";

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

  const [editorOpen, setEditorOpen] = useState(false);
  const [editorMode, setEditorMode] = useState<EditorMode>("");
  const [editorName, setEditorName] = useState("");
  const [editorYAML, setEditorYAML] = useState("");
  const [logsOpen, setLogsOpen] = useState(false);
  const [logsTitle, setLogsTitle] = useState("");
  const [logsText, setLogsText] = useState("");

  useEffect(() => {
    void load();
  }, [load]);

  if (loading) {
    return <CircularProgress />;
  }

  return (
    <Stack spacing={3}>
      <Typography variant="h5">工作负载管理（MVP）</Typography>
      {error && <Alert severity="error">{error}</Alert>}

      <Typography variant="h6">Deployment</Typography>
      {deployments.map((d) => (
        <Card key={d.name} variant="outlined">
          <CardContent>
            <Typography variant="h6">{d.name}</Typography>
            <Typography variant="body2" color="text.secondary">
              命名空间：{d.namespace}，镜像：{d.image}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              副本：{d.ready}/{d.replicas}，Age：{d.age}
            </Typography>
            <Stack direction="row" spacing={1} sx={{ mt: 2 }}>
              <Button
                size="small"
                onClick={async () => {
                  const yaml = await getDeploymentYAML(d.name);
                  setEditorMode("deployment");
                  setEditorName(d.name);
                  setEditorYAML(yaml);
                  setEditorOpen(true);
                }}
              >
                查看/编辑 YAML
              </Button>
            </Stack>
          </CardContent>
        </Card>
      ))}

      <Typography variant="h6">Pod</Typography>
      {pods.map((p) => (
        <Card key={p.name} variant="outlined">
          <CardContent>
            <Typography variant="h6">{p.name}</Typography>
            <Typography variant="body2" color="text.secondary">
              命名空间：{p.namespace}，状态：{p.status}，重启：{p.restarts}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              节点：{p.node}，IP：{p.ip}，Age：{p.age}
            </Typography>
            <Stack direction="row" spacing={1} sx={{ mt: 2 }}>
              <Button
                size="small"
                onClick={async () => {
                  const yaml = await getPodYAML(p.name);
                  setEditorMode("pod");
                  setEditorName(p.name);
                  setEditorYAML(yaml);
                  setEditorOpen(true);
                }}
              >
                查看/编辑 YAML
              </Button>
              <Button
                size="small"
                onClick={async () => {
                  const logs = await getPodLogs(p.name);
                  setLogsTitle(`Pod 日志 - ${p.name}`);
                  setLogsText(logs);
                  setLogsOpen(true);
                }}
              >
                查看日志
              </Button>
            </Stack>
          </CardContent>
        </Card>
      ))}

      <Dialog
        open={editorOpen}
        onClose={() => setEditorOpen(false)}
        fullWidth
        maxWidth="md"
      >
        <DialogTitle>YAML 编辑 - {editorName}</DialogTitle>
        <DialogContent>
          <TextField
            multiline
            minRows={16}
            fullWidth
            value={editorYAML}
            onChange={(e) => setEditorYAML(e.target.value)}
            sx={{ mt: 1 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditorOpen(false)}>关闭</Button>
          <Button
            variant="contained"
            onClick={async () => {
              const ok =
                editorMode === "deployment"
                  ? await saveDeploymentYAML(editorName, editorYAML)
                  : await savePodYAML(editorName, editorYAML);
              if (ok) {
                setEditorOpen(false);
              }
            }}
          >
            保存 YAML
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog open={logsOpen} onClose={() => setLogsOpen(false)} fullWidth maxWidth="md">
        <DialogTitle>{logsTitle}</DialogTitle>
        <DialogContent>
          <Box
            component="pre"
            sx={{
              m: 0,
              p: 2,
              bgcolor: "#0d1117",
              color: "#c9d1d9",
              borderRadius: 1,
              overflowX: "auto"
            }}
          >
            {logsText}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setLogsOpen(false)}>关闭</Button>
        </DialogActions>
      </Dialog>
    </Stack>
  );
}
