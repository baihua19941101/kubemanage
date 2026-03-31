import {
  Alert,
  Card,
  CardContent,
  CircularProgress,
  Divider,
  Stack,
  Typography
} from "@mui/material";
import { useEffect } from "react";
import { useResourceStore } from "../stores/useResourceStore";

export default function ResourcePage() {
  const services = useResourceStore((s) => s.services);
  const configMaps = useResourceStore((s) => s.configMaps);
  const secrets = useResourceStore((s) => s.secrets);
  const loading = useResourceStore((s) => s.loading);
  const error = useResourceStore((s) => s.error);
  const load = useResourceStore((s) => s.load);

  useEffect(() => {
    void load();
  }, [load]);

  if (loading) {
    return <CircularProgress />;
  }

  return (
    <Stack spacing={3}>
      <Typography variant="h5">服务与配置（MVP）</Typography>
      {error && <Alert severity="error">{error}</Alert>}

      <Typography variant="h6">Service</Typography>
      {services.map((item) => (
        <Card key={item.name} variant="outlined">
          <CardContent>
            <Typography variant="h6">{item.name}</Typography>
            <Typography variant="body2" color="text.secondary">
              命名空间：{item.namespace}，类型：{item.type}，ClusterIP：{item.clusterIP}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              端口：{item.ports}，关联 Pod：{item.pods}，Age：{item.age}
            </Typography>
          </CardContent>
        </Card>
      ))}

      <Divider />

      <Typography variant="h6">ConfigMap</Typography>
      {configMaps.map((item) => (
        <Card key={item.name} variant="outlined">
          <CardContent>
            <Typography variant="h6">{item.name}</Typography>
            <Typography variant="body2" color="text.secondary">
              命名空间：{item.namespace}，数据项：{item.dataCount}，Age：{item.age}
            </Typography>
          </CardContent>
        </Card>
      ))}

      <Divider />

      <Typography variant="h6">Secret（脱敏）</Typography>
      {secrets.map((item) => (
        <Card key={item.name} variant="outlined">
          <CardContent>
            <Typography variant="h6">{item.name}</Typography>
            <Typography variant="body2" color="text.secondary">
              命名空间：{item.namespace}，类型：{item.type}，Age：{item.age}
            </Typography>
            <Stack spacing={0.5} sx={{ mt: 1 }}>
              {Object.entries(item.data).map(([k, v]) => (
                <Typography key={k} variant="body2" color="text.secondary">
                  {k}: {v}
                </Typography>
              ))}
            </Stack>
          </CardContent>
        </Card>
      ))}
    </Stack>
  );
}
