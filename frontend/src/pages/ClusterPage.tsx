import {
  Button,
  Card,
  CardContent,
  Chip,
  CircularProgress,
  Stack,
  Typography
} from "@mui/material";
import { useEffect } from "react";
import { useClusterStore } from "../stores/useClusterStore";
import { useAuthStore } from "../stores/useAuthStore";

export default function ClusterPage() {
  const clusters = useClusterStore((s) => s.clusters);
  const current = useClusterStore((s) => s.current);
  const loading = useClusterStore((s) => s.loading);
  const switching = useClusterStore((s) => s.switching);
  const load = useClusterStore((s) => s.load);
  const switchCluster = useClusterStore((s) => s.switchCluster);
  const canWorkloadWrite = useAuthStore((s) => s.canWorkloadWrite);

  useEffect(() => {
    void load();
  }, [load]);

  if (loading) {
    return <CircularProgress />;
  }

  return (
    <Stack spacing={2}>
      <Typography variant="h5">集群管理（MVP）</Typography>
      {clusters.map((cluster) => (
        <Card key={cluster.name} variant="outlined">
          <CardContent>
            <Stack direction="row" spacing={2} alignItems="center" flexWrap="wrap">
              <Typography variant="h6">{cluster.name}</Typography>
              <Chip
                size="small"
                color={cluster.status === "ready" ? "success" : "default"}
                label={cluster.status}
              />
              {current === cluster.name && (
                <Chip size="small" color="primary" label="当前集群" />
              )}
            </Stack>
            <Typography variant="body2" color="text.secondary">
              版本：{cluster.version}，节点数：{cluster.nodes}
            </Typography>
            <Stack direction="row" justifyContent="flex-end" sx={{ mt: 2 }}>
              <Button
                size="small"
                variant="contained"
                disabled={
                  !canWorkloadWrite() ||
                  current === cluster.name ||
                  switching === cluster.name
                }
                onClick={() => {
                  void switchCluster(cluster.name);
                }}
              >
                {switching === cluster.name ? "切换中..." : "切换到该集群"}
              </Button>
            </Stack>
          </CardContent>
        </Card>
      ))}
    </Stack>
  );
}
