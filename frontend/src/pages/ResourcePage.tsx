import { Alert, Button, Stack, Typography } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import { useResourceStore } from "../stores/useResourceStore";

type Mode = "services" | "configmaps" | "secrets";

type ServiceItem = {
  name: string;
  namespace: string;
  type: string;
  clusterIP: string;
  ports: string;
  pods: number;
  age: string;
};

type ConfigMapItem = {
  name: string;
  namespace: string;
  dataCount: number;
  age: string;
};

type SecretItem = {
  name: string;
  namespace: string;
  type: string;
  data: Record<string, string>;
  age: string;
};

export default function ResourcePage() {
  const services = useResourceStore((s) => s.services);
  const configMaps = useResourceStore((s) => s.configMaps);
  const secrets = useResourceStore((s) => s.secrets);
  const loading = useResourceStore((s) => s.loading);
  const error = useResourceStore((s) => s.error);
  const load = useResourceStore((s) => s.load);
  const [mode, setMode] = useState<Mode>("services");
  const [selectedService, setSelectedService] = useState<ServiceItem | null>(null);
  const [selectedConfig, setSelectedConfig] = useState<ConfigMapItem | null>(null);
  const [selectedSecret, setSelectedSecret] = useState<SecretItem | null>(null);

  useEffect(() => {
    void load();
  }, [load]);

  const serviceColumns = useMemo(
    () => [
      { key: "name", header: "名称", render: (r: ServiceItem) => r.name },
      { key: "ns", header: "命名空间", render: (r: ServiceItem) => r.namespace },
      { key: "type", header: "类型", render: (r: ServiceItem) => r.type },
      { key: "ports", header: "端口", render: (r: ServiceItem) => r.ports },
      { key: "age", header: "Age", render: (r: ServiceItem) => r.age }
    ],
    []
  );
  const configColumns = useMemo(
    () => [
      { key: "name", header: "名称", render: (r: ConfigMapItem) => r.name },
      { key: "ns", header: "命名空间", render: (r: ConfigMapItem) => r.namespace },
      { key: "count", header: "数据项", render: (r: ConfigMapItem) => r.dataCount },
      { key: "age", header: "Age", render: (r: ConfigMapItem) => r.age }
    ],
    []
  );
  const secretColumns = useMemo(
    () => [
      { key: "name", header: "名称", render: (r: SecretItem) => r.name },
      { key: "ns", header: "命名空间", render: (r: SecretItem) => r.namespace },
      { key: "type", header: "类型", render: (r: SecretItem) => r.type },
      { key: "age", header: "Age", render: (r: SecretItem) => r.age }
    ],
    []
  );

  return (
    <>
      <PageScaffold
        title="服务与配置"
        description="统一管理 Service/ConfigMap/Secret，Secret 默认脱敏"
        actions={
          <Stack direction="row" spacing={1}>
            <Button
              variant={mode === "services" ? "contained" : "outlined"}
              onClick={() => setMode("services")}
            >
              Service
            </Button>
            <Button
              variant={mode === "configmaps" ? "contained" : "outlined"}
              onClick={() => setMode("configmaps")}
            >
              ConfigMap
            </Button>
            <Button
              variant={mode === "secrets" ? "contained" : "outlined"}
              onClick={() => setMode("secrets")}
            >
              Secret
            </Button>
          </Stack>
        }
      >
        {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}
        {mode === "services" && (
          <ResourceTable
            loading={loading}
            rows={services}
            rowKey={(r) => r.name}
            columns={serviceColumns}
            onRowClick={(r) => setSelectedService(r)}
          />
        )}
        {mode === "configmaps" && (
          <ResourceTable
            loading={loading}
            rows={configMaps}
            rowKey={(r) => r.name}
            columns={configColumns}
            onRowClick={(r) => setSelectedConfig(r)}
          />
        )}
        {mode === "secrets" && (
          <ResourceTable
            loading={loading}
            rows={secrets}
            rowKey={(r) => r.name}
            columns={secretColumns}
            onRowClick={(r) => setSelectedSecret(r)}
          />
        )}
      </PageScaffold>

      <DetailDrawer
        open={Boolean(selectedService)}
        title={selectedService ? `Service 详情 - ${selectedService.name}` : "Service 详情"}
        onClose={() => setSelectedService(null)}
      >
        {selectedService && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedService.name}</Typography>
            <Typography variant="body2">命名空间：{selectedService.namespace}</Typography>
            <Typography variant="body2">类型：{selectedService.type}</Typography>
            <Typography variant="body2">ClusterIP：{selectedService.clusterIP}</Typography>
            <Typography variant="body2">端口：{selectedService.ports}</Typography>
            <Typography variant="body2">关联 Pod：{selectedService.pods}</Typography>
          </Stack>
        )}
      </DetailDrawer>

      <DetailDrawer
        open={Boolean(selectedConfig)}
        title={selectedConfig ? `ConfigMap 详情 - ${selectedConfig.name}` : "ConfigMap 详情"}
        onClose={() => setSelectedConfig(null)}
      >
        {selectedConfig && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedConfig.name}</Typography>
            <Typography variant="body2">命名空间：{selectedConfig.namespace}</Typography>
            <Typography variant="body2">数据项：{selectedConfig.dataCount}</Typography>
            <Typography variant="body2">Age：{selectedConfig.age}</Typography>
          </Stack>
        )}
      </DetailDrawer>

      <DetailDrawer
        open={Boolean(selectedSecret)}
        title={selectedSecret ? `Secret 详情 - ${selectedSecret.name}` : "Secret 详情"}
        onClose={() => setSelectedSecret(null)}
      >
        {selectedSecret && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedSecret.name}</Typography>
            <Typography variant="body2">命名空间：{selectedSecret.namespace}</Typography>
            <Typography variant="body2">类型：{selectedSecret.type}</Typography>
            <Typography variant="body2">Age：{selectedSecret.age}</Typography>
            <Typography variant="body2" color="text.secondary">
              脱敏数据：
            </Typography>
            {Object.entries(selectedSecret.data).map(([k, v]) => (
              <Typography key={k} variant="body2" color="text.secondary">
                {k}: {v}
              </Typography>
            ))}
          </Stack>
        )}
      </DetailDrawer>
    </>
  );
}
