import { Alert, Button, Stack, Typography } from "@mui/material";
import { useEffect, useMemo, useState } from "react";
import DetailDrawer from "../components/framework/DetailDrawer";
import PageScaffold from "../components/framework/PageScaffold";
import ResourceTable from "../components/framework/ResourceTable";
import { useResourceStore } from "../stores/useResourceStore";

type Mode = "services" | "ingresses" | "hpas" | "configmaps" | "secrets";

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

type IngressItem = {
  name: string;
  namespace: string;
  className: string;
  hosts: string[];
  address: string;
  tls: boolean;
  age: string;
};

type HPAItem = {
  name: string;
  namespace: string;
  targetKind: string;
  targetName: string;
  minReplicas: number;
  maxReplicas: number;
  currentReplicas: number;
  targetCPUPercent: number;
  currentCPUPercent: number;
  age: string;
};

type HPATarget = {
  kind: string;
  name: string;
  namespace: string;
  currentReplicas: number;
  desiredReplicas: number;
};

export default function ResourcePage() {
  const services = useResourceStore((s) => s.services);
  const ingresses = useResourceStore((s) => s.ingresses);
  const hpas = useResourceStore((s) => s.hpas);
  const configMaps = useResourceStore((s) => s.configMaps);
  const secrets = useResourceStore((s) => s.secrets);
  const loading = useResourceStore((s) => s.loading);
  const error = useResourceStore((s) => s.error);
  const load = useResourceStore((s) => s.load);
  const getIngressServices = useResourceStore((s) => s.getIngressServices);
  const getHPATarget = useResourceStore((s) => s.getHPATarget);
  const [mode, setMode] = useState<Mode>("services");
  const [selectedService, setSelectedService] = useState<ServiceItem | null>(null);
  const [selectedIngress, setSelectedIngress] = useState<IngressItem | null>(null);
  const [selectedHPA, setSelectedHPA] = useState<HPAItem | null>(null);
  const [selectedConfig, setSelectedConfig] = useState<ConfigMapItem | null>(null);
  const [selectedSecret, setSelectedSecret] = useState<SecretItem | null>(null);
  const [ingressServices, setIngressServices] = useState<ServiceItem[]>([]);
  const [hpaTarget, setHPATarget] = useState<HPATarget | null>(null);

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
  const ingressColumns = useMemo(
    () => [
      { key: "name", header: "名称", render: (r: IngressItem) => r.name },
      { key: "ns", header: "命名空间", render: (r: IngressItem) => r.namespace },
      { key: "class", header: "Class", render: (r: IngressItem) => r.className },
      { key: "hosts", header: "Hosts", render: (r: IngressItem) => r.hosts.join(", ") },
      { key: "tls", header: "TLS", render: (r: IngressItem) => (r.tls ? "是" : "否") }
    ],
    []
  );
  const hpaColumns = useMemo(
    () => [
      { key: "name", header: "名称", render: (r: HPAItem) => r.name },
      { key: "ns", header: "命名空间", render: (r: HPAItem) => r.namespace },
      { key: "target", header: "目标", render: (r: HPAItem) => `${r.targetKind}/${r.targetName}` },
      { key: "replicas", header: "副本", render: (r: HPAItem) => `${r.currentReplicas} (${r.minReplicas}-${r.maxReplicas})` },
      { key: "cpu", header: "CPU", render: (r: HPAItem) => `${r.currentCPUPercent}% / ${r.targetCPUPercent}%` }
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
              variant={mode === "ingresses" ? "contained" : "outlined"}
              onClick={() => setMode("ingresses")}
            >
              Ingress
            </Button>
            <Button
              variant={mode === "hpas" ? "contained" : "outlined"}
              onClick={() => setMode("hpas")}
            >
              HPA
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
        {mode === "ingresses" && (
          <ResourceTable
            loading={loading}
            rows={ingresses}
            rowKey={(r) => r.name}
            columns={ingressColumns}
            onRowClick={async (r) => {
              setSelectedIngress(r);
              setIngressServices(await getIngressServices(r.name));
            }}
          />
        )}
        {mode === "hpas" && (
          <ResourceTable
            loading={loading}
            rows={hpas}
            rowKey={(r) => r.name}
            columns={hpaColumns}
            onRowClick={async (r) => {
              setSelectedHPA(r);
              setHPATarget(await getHPATarget(r.name));
            }}
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
        open={Boolean(selectedIngress)}
        title={selectedIngress ? `Ingress 详情 - ${selectedIngress.name}` : "Ingress 详情"}
        onClose={() => {
          setSelectedIngress(null);
          setIngressServices([]);
        }}
      >
        {selectedIngress && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedIngress.name}</Typography>
            <Typography variant="body2">命名空间：{selectedIngress.namespace}</Typography>
            <Typography variant="body2">Class：{selectedIngress.className}</Typography>
            <Typography variant="body2">Hosts：{selectedIngress.hosts.join(", ")}</Typography>
            <Typography variant="body2">Address：{selectedIngress.address}</Typography>
            <Typography variant="body2">TLS：{selectedIngress.tls ? "是" : "否"}</Typography>
            <Typography variant="body2" color="text.secondary">关联 Service：</Typography>
            {ingressServices.map((item) => (
              <Typography key={item.name} variant="body2" color="text.secondary">
                {item.namespace}/{item.name} ({item.ports})
              </Typography>
            ))}
          </Stack>
        )}
      </DetailDrawer>

      <DetailDrawer
        open={Boolean(selectedHPA)}
        title={selectedHPA ? `HPA 详情 - ${selectedHPA.name}` : "HPA 详情"}
        onClose={() => {
          setSelectedHPA(null);
          setHPATarget(null);
        }}
      >
        {selectedHPA && (
          <Stack spacing={1}>
            <Typography variant="body2">名称：{selectedHPA.name}</Typography>
            <Typography variant="body2">命名空间：{selectedHPA.namespace}</Typography>
            <Typography variant="body2">目标：{selectedHPA.targetKind}/{selectedHPA.targetName}</Typography>
            <Typography variant="body2">副本范围：{selectedHPA.minReplicas} - {selectedHPA.maxReplicas}</Typography>
            <Typography variant="body2">当前副本：{selectedHPA.currentReplicas}</Typography>
            <Typography variant="body2">CPU：{selectedHPA.currentCPUPercent}% / {selectedHPA.targetCPUPercent}%</Typography>
            {hpaTarget && (
              <>
                <Typography variant="body2" color="text.secondary">关联目标：</Typography>
                <Typography variant="body2" color="text.secondary">
                  {hpaTarget.namespace}/{hpaTarget.kind}/{hpaTarget.name}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  当前/期望副本：{hpaTarget.currentReplicas}/{hpaTarget.desiredReplicas}
                </Typography>
              </>
            )}
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
