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
import { useNamespaceStore } from "../stores/useNamespaceStore";

export default function NamespacePage() {
  const items = useNamespaceStore((s) => s.items);
  const loading = useNamespaceStore((s) => s.loading);
  const error = useNamespaceStore((s) => s.error);
  const load = useNamespaceStore((s) => s.load);
  const create = useNamespaceStore((s) => s.create);
  const remove = useNamespaceStore((s) => s.remove);
  const fetchYaml = useNamespaceStore((s) => s.fetchYaml);

  const [open, setOpen] = useState(false);
  const [newName, setNewName] = useState("");
  const [yamlText, setYamlText] = useState("");
  const [yamlTitle, setYamlTitle] = useState("");

  useEffect(() => {
    void load();
  }, [load]);

  if (loading) {
    return <CircularProgress />;
  }

  return (
    <Stack spacing={2}>
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Typography variant="h5">名称空间管理（MVP）</Typography>
        <Button variant="contained" onClick={() => setOpen(true)}>
          新建名称空间
        </Button>
      </Stack>

      {error && <Alert severity="error">{error}</Alert>}

      {items.map((ns) => (
        <Card key={ns.name} variant="outlined">
          <CardContent>
            <Typography variant="h6">{ns.name}</Typography>
            <Typography variant="body2" color="text.secondary">
              状态：{ns.status}，存在时间：{ns.age}
            </Typography>
            <Stack direction="row" spacing={1} sx={{ mt: 2 }}>
              <Button
                size="small"
                onClick={async () => {
                  const text = await fetchYaml(ns.name);
                  setYamlTitle(`YAML - ${ns.name}`);
                  setYamlText(text);
                }}
              >
                查看 YAML
              </Button>
              <Button
                size="small"
                component="a"
                href={`/api/v1/namespaces/${ns.name}/yaml/download`}
              >
                下载 YAML
              </Button>
              <Button
                size="small"
                color="error"
                onClick={async () => {
                  if (window.confirm(`确认删除名称空间 ${ns.name} ?`)) {
                    await remove(ns.name);
                  }
                }}
              >
                删除
              </Button>
            </Stack>
          </CardContent>
        </Card>
      ))}

      <Dialog open={open} onClose={() => setOpen(false)} fullWidth maxWidth="xs">
        <DialogTitle>新建名称空间</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="名称空间名称"
            fullWidth
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>取消</Button>
          <Button
            variant="contained"
            onClick={async () => {
              await create(newName);
              setNewName("");
              setOpen(false);
            }}
          >
            创建
          </Button>
        </DialogActions>
      </Dialog>

      <Dialog
        open={yamlText !== ""}
        onClose={() => setYamlText("")}
        fullWidth
        maxWidth="md"
      >
        <DialogTitle>{yamlTitle}</DialogTitle>
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
            {yamlText}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setYamlText("")}>关闭</Button>
        </DialogActions>
      </Dialog>
    </Stack>
  );
}
